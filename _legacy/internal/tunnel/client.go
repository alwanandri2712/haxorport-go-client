package tunnel

import (
	"encoding/json"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/haxorport/client/internal/config"
	"github.com/haxorport/client/internal/logger"
	"github.com/haxorport/client/internal/proto"
)

// Client adalah klien untuk berkomunikasi dengan server haxorport
type Client struct {
	config     *config.Config
	logger     *logger.Logger
	conn       *websocket.Conn
	tunnels    map[string]*Tunnel
	tunnelsMu  sync.RWMutex
	done       chan struct{}
	reconnect  chan struct{}
	pingTicker *time.Ticker
}

// NewClient membuat instance Client baru
func NewClient(config *config.Config, logger *logger.Logger) *Client {
	return &Client{
		config:    config,
		logger:    logger,
		tunnels:   make(map[string]*Tunnel),
		done:      make(chan struct{}),
		reconnect: make(chan struct{}, 1),
	}
}

// Connect menghubungkan ke server haxorport
func (c *Client) Connect() error {
	// Buat URL WebSocket
	u := url.URL{
		Scheme: "ws",
		Host:   fmt.Sprintf("%s:%d", c.config.ServerAddress, c.config.ControlPort),
		Path:   "/control",
	}

	c.logger.Info("Menghubungkan ke server haxorport di %s", u.String())

	// Hubungkan ke server
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return fmt.Errorf("error menghubungkan ke server: %v", err)
	}

	c.conn = conn

	// Autentikasi jika token disediakan
	if c.config.AuthToken != "" {
		if err := c.authenticate(); err != nil {
			c.conn.Close()
			return fmt.Errorf("error autentikasi: %v", err)
		}
	}

	// Mulai ping ticker untuk menjaga koneksi tetap hidup
	c.pingTicker = time.NewTicker(30 * time.Second)

	// Mulai goroutine untuk membaca pesan
	go c.readPump()

	// Mulai goroutine untuk mengirim ping
	go c.pingPump()

	c.logger.Info("Terhubung ke server haxorport")

	return nil
}

// authenticate melakukan autentikasi ke server
func (c *Client) authenticate() error {
	// Buat payload autentikasi
	authPayload := proto.AuthPayload{
		Token: c.config.AuthToken,
	}

	// Buat pesan autentikasi
	authMsg, err := proto.NewMessage(proto.MessageTypeAuth, authPayload)
	if err != nil {
		return fmt.Errorf("error membuat pesan autentikasi: %v", err)
	}

	// Kirim pesan autentikasi
	if err := c.sendMessage(authMsg); err != nil {
		return fmt.Errorf("error mengirim pesan autentikasi: %v", err)
	}

	c.logger.Debug("Pesan autentikasi terkirim")

	return nil
}

// RegisterTunnel mendaftarkan tunnel baru ke server
func (c *Client) RegisterTunnel(tunnelConfig config.TunnelConfig) (*Tunnel, error) {
	// Buat payload pendaftaran
	registerPayload := proto.RegisterPayload{
		TunnelType: tunnelConfig.Type,
		LocalPort:  tunnelConfig.LocalPort,
	}

	// Tambahkan subdomain jika ada
	if tunnelConfig.Subdomain != "" {
		registerPayload.Subdomain = tunnelConfig.Subdomain
	}

	// Tambahkan remote port jika ada
	if tunnelConfig.RemotePort != 0 {
		registerPayload.RemotePort = tunnelConfig.RemotePort
	}

	// Tambahkan autentikasi jika ada
	if tunnelConfig.Auth != nil {
		registerPayload.Auth = &proto.TunnelAuth{
			Type:        tunnelConfig.Auth.Type,
			Username:    tunnelConfig.Auth.Username,
			Password:    tunnelConfig.Auth.Password,
			HeaderName:  tunnelConfig.Auth.HeaderName,
			HeaderValue: tunnelConfig.Auth.HeaderValue,
		}
	}

	// Buat pesan pendaftaran
	registerMsg, err := proto.NewMessage(proto.MessageTypeRegister, registerPayload)
	if err != nil {
		return nil, fmt.Errorf("error membuat pesan pendaftaran: %v", err)
	}

	// Kirim pesan pendaftaran
	if err := c.sendMessage(registerMsg); err != nil {
		return nil, fmt.Errorf("error mengirim pesan pendaftaran: %v", err)
	}

	c.logger.Debug("Pesan pendaftaran terkirim untuk tunnel %s:%d", tunnelConfig.Type, tunnelConfig.LocalPort)

	// Tunggu respons dari server
	// Respons akan ditangani di readPump dan tunnel akan ditambahkan ke map tunnels
	// Kita perlu menunggu beberapa saat untuk memastikan respons diterima
	time.Sleep(500 * time.Millisecond)

	// Periksa apakah tunnel berhasil didaftarkan
	c.tunnelsMu.RLock()
	defer c.tunnelsMu.RUnlock()

	// Cari tunnel berdasarkan konfigurasi
	for _, tunnel := range c.tunnels {
		if tunnel.Config.Type == tunnelConfig.Type && tunnel.Config.LocalPort == tunnelConfig.LocalPort {
			return tunnel, nil
		}
	}

	return nil, fmt.Errorf("tunnel tidak berhasil didaftarkan")
}

// UnregisterTunnel menghapus tunnel dari server
func (c *Client) UnregisterTunnel(tunnelID string) error {
	// Periksa apakah tunnel ada
	c.tunnelsMu.RLock()
	tunnel, exists := c.tunnels[tunnelID]
	c.tunnelsMu.RUnlock()

	if !exists {
		return fmt.Errorf("tunnel %s tidak ditemukan", tunnelID)
	}

	// Buat payload penghapusan
	unregisterPayload := proto.UnregisterPayload{
		TunnelID: tunnelID,
	}

	// Buat pesan penghapusan
	unregisterMsg, err := proto.NewMessage(proto.MessageTypeUnregister, unregisterPayload)
	if err != nil {
		return fmt.Errorf("error membuat pesan penghapusan: %v", err)
	}

	// Kirim pesan penghapusan
	if err := c.sendMessage(unregisterMsg); err != nil {
		return fmt.Errorf("error mengirim pesan penghapusan: %v", err)
	}

	c.logger.Debug("Pesan penghapusan terkirim untuk tunnel %s", tunnelID)

	// Hentikan tunnel
	tunnel.Stop()

	// Hapus tunnel dari map
	c.tunnelsMu.Lock()
	delete(c.tunnels, tunnelID)
	c.tunnelsMu.Unlock()

	c.logger.Info("Tunnel %s dihapus", tunnelID)

	return nil
}

// Close menutup koneksi ke server
func (c *Client) Close() {
	// Hentikan semua tunnel
	c.tunnelsMu.Lock()
	for _, tunnel := range c.tunnels {
		tunnel.Stop()
	}
	c.tunnels = make(map[string]*Tunnel)
	c.tunnelsMu.Unlock()

	// Hentikan ping ticker
	if c.pingTicker != nil {
		c.pingTicker.Stop()
	}

	// Tutup channel done
	close(c.done)

	// Tutup koneksi WebSocket
	if c.conn != nil {
		c.conn.Close()
	}

	c.logger.Info("Koneksi ke server haxorport ditutup")
}

// sendMessage mengirim pesan ke server
func (c *Client) sendMessage(msg *proto.Message) error {
	// Marshal pesan
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("error marshaling pesan: %v", err)
	}

	// Kirim pesan
	c.conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	err = c.conn.WriteMessage(websocket.TextMessage, msgBytes)
	c.conn.SetWriteDeadline(time.Time{})
	if err != nil {
		return fmt.Errorf("error mengirim pesan: %v", err)
	}

	return nil
}

// readPump membaca pesan dari server
func (c *Client) readPump() {
	defer func() {
		c.logger.Debug("readPump berhenti")
		// Trigger reconnect
		select {
		case c.reconnect <- struct{}{}:
		default:
		}
	}()

	for {
		select {
		case <-c.done:
			return
		default:
			// Baca pesan
			_, msgBytes, err := c.conn.ReadMessage()
			if err != nil {
				c.logger.Error("Error membaca pesan: %v", err)
				return
			}

			// Parse pesan
			var msg proto.Message
			if err := json.Unmarshal(msgBytes, &msg); err != nil {
				c.logger.Error("Error parsing pesan: %v", err)
				continue
			}

			// Tangani pesan berdasarkan tipe
			switch msg.Type {
			case proto.MessageTypeRegister:
				c.handleRegisterResponse(&msg)
			case proto.MessageTypeData:
				c.handleDataMessage(&msg)
			case proto.MessageTypePong:
				// Tidak perlu melakukan apa-apa untuk pong
			case proto.MessageTypeError:
				c.handleErrorMessage(&msg)
			default:
				c.logger.Warn("Tipe pesan tidak dikenal: %s", msg.Type)
			}
		}
	}
}

// pingPump mengirim ping ke server secara berkala
func (c *Client) pingPump() {
	defer func() {
		c.logger.Debug("pingPump berhenti")
	}()

	for {
		select {
		case <-c.done:
			return
		case <-c.pingTicker.C:
			// Buat pesan ping
			pingMsg, err := proto.NewMessage(proto.MessageTypePing, nil)
			if err != nil {
				c.logger.Error("Error membuat pesan ping: %v", err)
				continue
			}

			// Kirim pesan ping
			if err := c.sendMessage(pingMsg); err != nil {
				c.logger.Error("Error mengirim pesan ping: %v", err)
				// Trigger reconnect
				select {
				case c.reconnect <- struct{}{}:
				default:
				}
				return
			}

			c.logger.Debug("Ping terkirim")
		}
	}
}

// handleRegisterResponse menangani respons pendaftaran tunnel
func (c *Client) handleRegisterResponse(msg *proto.Message) {
	// Parse payload
	var payload proto.RegisterResponsePayload
	if err := msg.ParsePayload(&payload); err != nil {
		c.logger.Error("Error parsing payload respons pendaftaran: %v", err)
		return
	}

	// Periksa apakah pendaftaran berhasil
	if !payload.Success {
		c.logger.Error("Pendaftaran tunnel gagal: %s", payload.Error)
		return
	}

	// Buat objek tunnel
	tunnelType := "http"
	if payload.RemotePort > 0 {
		tunnelType = "tcp"
	}

	tunnelConfig := config.TunnelConfig{
		Type:      tunnelType,
		LocalPort: 0, // Akan diisi nanti
	}

	if tunnelType == "http" {
		tunnelConfig.Subdomain = payload.URL
	} else {
		tunnelConfig.RemotePort = payload.RemotePort
	}

	tunnel := NewTunnel(payload.TunnelID, tunnelConfig, c)

	// Tambahkan tunnel ke map
	c.tunnelsMu.Lock()
	c.tunnels[payload.TunnelID] = tunnel
	c.tunnelsMu.Unlock()

	c.logger.Info("Tunnel terdaftar: %s, ID: %s", tunnelType, payload.TunnelID)
	if tunnelType == "http" {
		c.logger.Info("URL tunnel: %s", payload.URL)
	} else {
		c.logger.Info("Port remote tunnel: %d", payload.RemotePort)
	}

	// Mulai tunnel
	go tunnel.Start()
}

// handleDataMessage menangani pesan data
func (c *Client) handleDataMessage(msg *proto.Message) {
	// Parse payload
	var payload proto.DataPayload
	if err := msg.ParsePayload(&payload); err != nil {
		c.logger.Error("Error parsing payload data: %v", err)
		return
	}

	// Dapatkan tunnel
	c.tunnelsMu.RLock()
	tunnel, exists := c.tunnels[payload.TunnelID]
	c.tunnelsMu.RUnlock()

	if !exists {
		c.logger.Warn("Tunnel %s tidak ditemukan untuk pesan data", payload.TunnelID)
		return
	}

	// Teruskan data ke tunnel
	if err := tunnel.HandleData(payload.ConnectionID, payload.Data); err != nil {
		c.logger.Error("Error menangani data: %v", err)
	}
}

// handleErrorMessage menangani pesan error
func (c *Client) handleErrorMessage(msg *proto.Message) {
	// Parse payload
	var payload proto.ErrorPayload
	if err := msg.ParsePayload(&payload); err != nil {
		c.logger.Error("Error parsing payload error: %v", err)
		return
	}

	c.logger.Error("Error dari server: %s - %s", payload.Code, payload.Message)
}

// RunWithReconnect menjalankan klien dengan reconnect otomatis
func (c *Client) RunWithReconnect() {
	backoff := 1 * time.Second
	maxBackoff := 60 * time.Second

	for {
		// Coba connect
		err := c.Connect()
		if err != nil {
			c.logger.Error("Gagal menghubungkan ke server: %v", err)
			c.logger.Info("Mencoba kembali dalam %v...", backoff)
			time.Sleep(backoff)
			
			// Tingkatkan backoff
			backoff *= 2
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
			
			continue
		}

		// Reset backoff
		backoff = 1 * time.Second

		// Tunggu sinyal reconnect atau done
		select {
		case <-c.done:
			return
		case <-c.reconnect:
			c.logger.Info("Koneksi terputus, mencoba menghubungkan kembali...")
			// Tutup koneksi lama
			if c.conn != nil {
				c.conn.Close()
			}
			// Hentikan ping ticker
			if c.pingTicker != nil {
				c.pingTicker.Stop()
			}
			// Tunggu sebentar sebelum reconnect
			time.Sleep(backoff)
		}
	}
}

// GetTunnels mengembalikan daftar tunnel yang aktif
func (c *Client) GetTunnels() []*Tunnel {
	c.tunnelsMu.RLock()
	defer c.tunnelsMu.RUnlock()

	tunnels := make([]*Tunnel, 0, len(c.tunnels))
	for _, tunnel := range c.tunnels {
		tunnels = append(tunnels, tunnel)
	}

	return tunnels
}

// GetTunnel mengembalikan tunnel berdasarkan ID
func (c *Client) GetTunnel(tunnelID string) (*Tunnel, error) {
	c.tunnelsMu.RLock()
	defer c.tunnelsMu.RUnlock()

	tunnel, exists := c.tunnels[tunnelID]
	if !exists {
		return nil, fmt.Errorf("tunnel %s tidak ditemukan", tunnelID)
	}

	return tunnel, nil
}
