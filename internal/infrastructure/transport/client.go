package transport

import (
	"encoding/json"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/haxorport/haxor-client/internal/domain/model"
	"github.com/haxorport/haxor-client/internal/domain/port"
)

// Client adalah implementasi port.Client
type Client struct {
	serverAddr   string
	serverPort   int
	authToken    string
	conn         *websocket.Conn
	isConnected  bool
	reconnecting bool
	mutex        sync.Mutex
	logger       port.Logger
	handlers     map[model.MessageType]func(*model.Message) error
}

// NewClient membuat instance Client baru
func NewClient(serverAddress string, controlPort int, authToken string, logger port.Logger) *Client {
	return &Client{
		serverAddr:   serverAddress,
		serverPort:   controlPort,
		authToken:    authToken,
		isConnected:  false,
		reconnecting: false,
		logger:       logger,
		handlers:     make(map[model.MessageType]func(*model.Message) error),
	}
}

// Connect menghubungkan ke server haxorport
func (c *Client) Connect() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.isConnected {
		return nil
	}

	serverURL := fmt.Sprintf("ws://%s:%d/control", c.serverAddr, c.serverPort)
	c.logger.Info("Menghubungkan ke server: %s", serverURL)

	// Parse URL
	u, err := url.Parse(serverURL)
	if err != nil {
		return fmt.Errorf("URL tidak valid: %v", err)
	}

	// Buat koneksi WebSocket
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return fmt.Errorf("gagal terhubung ke server: %v", err)
	}

	c.conn = conn
	c.isConnected = true

	// Siapkan pesan autentikasi
	authPayload := model.AuthPayload{
		Token: c.authToken,
	}
	authMessage, err := model.NewMessage(model.MessageTypeAuth, authPayload)
	if err != nil {
		c.Close()
		return fmt.Errorf("gagal membuat pesan autentikasi: %v", err)
	}

	// Marshal pesan ke JSON
	data, err := json.Marshal(authMessage)
	if err != nil {
		c.Close()
		return fmt.Errorf("gagal mengkonversi pesan autentikasi ke JSON: %v", err)
	}

	// Kirim pesan langsung tanpa memanggil sendMessage (untuk menghindari deadlock)
	if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
		c.logger.Error("Gagal mengirim pesan autentikasi: %v", err)
		c.Close()
		return fmt.Errorf("gagal mengirim autentikasi: %v", err)
	}

	// Mulai goroutine untuk membaca pesan
	go c.readPump()

	c.logger.Info("Terhubung ke server: %s", serverURL)

	return nil
}

// Close menutup koneksi ke server
func (c *Client) Close() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if !c.isConnected {
		return
	}

	c.logger.Info("Menutup koneksi ke server")

	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}

	c.isConnected = false
}

// IsConnected mengembalikan status koneksi
func (c *Client) IsConnected() bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.isConnected
}

// RunWithReconnect menjalankan klien dengan reconnect otomatis
func (c *Client) RunWithReconnect() {
	c.mutex.Lock()
	if c.reconnecting {
		c.mutex.Unlock()
		return
	}
	c.reconnecting = true
	c.mutex.Unlock()

	go func() {
		for {
			if !c.IsConnected() {
				c.logger.Info("Mencoba menghubungkan kembali ke server...")
				if err := c.Connect(); err != nil {
					c.logger.Error("Gagal menghubungkan kembali: %v", err)
					time.Sleep(5 * time.Second)
					continue
				}
			}
			time.Sleep(1 * time.Second)
		}
	}()

	// Kirim ping secara berkala
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			if c.IsConnected() {
				pingMessage, err := model.NewMessage(model.MessageTypePing, nil)
				if err != nil {
					c.logger.Error("Gagal membuat pesan ping: %v", err)
					continue
				}

				if err := c.sendMessage(pingMessage); err != nil {
					c.logger.Error("Gagal mengirim ping: %v", err)
					c.Close()
				}
			}
		}
	}()
}

// RegisterHandler mendaftarkan handler untuk tipe pesan tertentu
func (c *Client) RegisterHandler(msgType model.MessageType, handler func(*model.Message) error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.handlers[msgType] = handler
}

// sendMessage mengirim pesan ke server
func (c *Client) sendMessage(msg *model.Message) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if !c.isConnected || c.conn == nil {
		return fmt.Errorf("tidak terhubung ke server")
	}

	// Marshal pesan ke JSON
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("gagal mengkonversi pesan ke JSON: %v", err)
	}

	c.logger.Debug("Mengirim pesan: %s", string(data))

	// Kirim pesan
	if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
		c.logger.Error("Gagal mengirim pesan: %v", err)
		c.isConnected = false
		return fmt.Errorf("gagal mengirim pesan: %v", err)
	}

	return nil
}

// readPump membaca pesan dari server
func (c *Client) readPump() {
	defer c.Close()

	for {
		// Baca pesan dari WebSocket
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			c.logger.Error("Gagal membaca pesan: %v", err)
			break
		}

		c.logger.Debug("Menerima pesan: %s", string(data))

		// Parse pesan
		var msg model.Message
		if err := json.Unmarshal(data, &msg); err != nil {
			c.logger.Error("Gagal mengurai pesan: %v", err)
			continue
		}

		// Tangani pesan pong secara khusus
		if msg.Type == model.MessageTypePong {
			c.logger.Debug("Menerima pong dari server")
			continue
		}

		// Tangani pesan dengan handler yang terdaftar
		c.mutex.Lock()
		handler, exists := c.handlers[msg.Type]
		c.mutex.Unlock()

		if exists {
			if err := handler(&msg); err != nil {
				c.logger.Error("Error menangani pesan %s: %v", msg.Type, err)
			}
		} else {
			c.logger.Warn("Tidak ada handler untuk pesan tipe: %s", msg.Type)
		}
	}
}

// SendRegisterTunnel mengirim permintaan pendaftaran tunnel
func (c *Client) SendRegisterTunnel(config model.TunnelConfig) (*model.RegisterResponsePayload, error) {
	// Buat channel untuk menerima respons
	responseCh := make(chan *model.RegisterResponsePayload, 1)
	errCh := make(chan error, 1)

	// Daftarkan handler untuk pesan register
	c.RegisterHandler(model.MessageTypeRegister, func(msg *model.Message) error {
		var response model.RegisterResponsePayload
		if err := msg.ParsePayload(&response); err != nil {
			errCh <- fmt.Errorf("gagal mengurai respons: %v", err)
			return err
		}

		responseCh <- &response
		return nil
	})

	// Daftarkan handler untuk pesan error
	c.RegisterHandler(model.MessageTypeError, func(msg *model.Message) error {
		var errorPayload model.ErrorPayload
		if err := msg.ParsePayload(&errorPayload); err != nil {
			errCh <- fmt.Errorf("gagal mengurai pesan error: %v", err)
			return err
		}

		errCh <- fmt.Errorf("error dari server: %s - %s", errorPayload.Code, errorPayload.Message)
		return nil
	})

	// Kirim permintaan pendaftaran
	payload := model.RegisterPayload{
		TunnelType: string(config.Type),
		Subdomain:  config.Subdomain,
		LocalPort:  config.LocalPort,
		RemotePort: config.RemotePort,
		Auth:       config.Auth,
	}

	msg, err := model.NewMessage(model.MessageTypeRegister, payload)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat pesan: %v", err)
	}

	if err := c.sendMessage(msg); err != nil {
		return nil, err
	}

	// Tunggu respons atau error
	select {
	case response := <-responseCh:
		if !response.Success {
			return nil, fmt.Errorf("pendaftaran tunnel gagal: %s", response.Error)
		}
		return response, nil
	case err := <-errCh:
		return nil, err
	case <-time.After(10 * time.Second):
		return nil, fmt.Errorf("timeout menunggu respons dari server")
	}
}

// SendUnregisterTunnel mengirim permintaan penghapusan tunnel
func (c *Client) SendUnregisterTunnel(tunnelID string) error {
	payload := model.UnregisterPayload{
		TunnelID: tunnelID,
	}

	msg, err := model.NewMessage(model.MessageTypeUnregister, payload)
	if err != nil {
		return fmt.Errorf("gagal membuat pesan: %v", err)
	}

	return c.sendMessage(msg)
}

// SendData mengirim data melalui tunnel
func (c *Client) SendData(tunnelID string, connectionID string, data []byte) error {
	payload := model.DataPayload{
		TunnelID:     tunnelID,
		ConnectionID: connectionID,
		Data:         data,
	}

	msg, err := model.NewMessage(model.MessageTypeData, payload)
	if err != nil {
		return fmt.Errorf("gagal membuat pesan: %v", err)
	}

	return c.sendMessage(msg)
}

// Ensure Client implements port.Client
var _ port.Client = (*Client)(nil)
