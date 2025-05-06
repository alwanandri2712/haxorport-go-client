package tunnel

import (
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/haxorport/client/internal/config"
	"github.com/haxorport/client/internal/proto"
)

// Tunnel merepresentasikan tunnel yang dibuat oleh klien
type Tunnel struct {
	ID            string
	Config        config.TunnelConfig
	client        *Client
	connections   map[string]net.Conn
	connectionsMu sync.RWMutex
	listener      net.Listener
	done          chan struct{}
}

// NewTunnel membuat instance Tunnel baru
func NewTunnel(id string, config config.TunnelConfig, client *Client) *Tunnel {
	return &Tunnel{
		ID:          id,
		Config:      config,
		client:      client,
		connections: make(map[string]net.Conn),
		done:        make(chan struct{}),
	}
}

// Start memulai tunnel
func (t *Tunnel) Start() error {
	// Buat listener untuk menerima koneksi lokal
	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", t.Config.LocalPort))
	if err != nil {
		return fmt.Errorf("error membuat listener: %v", err)
	}

	t.listener = listener
	t.client.logger.Info("Tunnel %s mendengarkan di localhost:%d", t.ID, t.Config.LocalPort)

	// Loop untuk menerima koneksi
	go func() {
		for {
			select {
			case <-t.done:
				return
			default:
				// Terima koneksi
				conn, err := listener.Accept()
				if err != nil {
					if isClosedError(err) {
						return
					}
					t.client.logger.Error("Error menerima koneksi: %v", err)
					continue
				}

				// Tangani koneksi dalam goroutine
				go t.handleConnection(conn)
			}
		}
	}()

	return nil
}

// Stop menghentikan tunnel
func (t *Tunnel) Stop() {
	// Tutup channel done
	close(t.done)

	// Tutup listener
	if t.listener != nil {
		t.listener.Close()
	}

	// Tutup semua koneksi
	t.connectionsMu.Lock()
	for _, conn := range t.connections {
		conn.Close()
	}
	t.connections = make(map[string]net.Conn)
	t.connectionsMu.Unlock()

	t.client.logger.Info("Tunnel %s dihentikan", t.ID)
}

// handleConnection menangani koneksi lokal
func (t *Tunnel) handleConnection(conn net.Conn) {
	// Buat ID koneksi
	connID := generateRandomString(16)

	// Simpan koneksi
	t.connectionsMu.Lock()
	t.connections[connID] = conn
	t.connectionsMu.Unlock()

	// Pastikan koneksi ditutup dan dihapus ketika fungsi selesai
	defer func() {
		conn.Close()
		t.connectionsMu.Lock()
		delete(t.connections, connID)
		t.connectionsMu.Unlock()
	}()

	t.client.logger.Debug("Koneksi baru untuk tunnel %s: %s", t.ID, connID)

	// Buat buffer untuk membaca data
	buffer := make([]byte, 4096)

	// Loop untuk membaca data dari koneksi lokal
	for {
		select {
		case <-t.done:
			return
		default:
			// Baca data
			n, err := conn.Read(buffer)
			if err != nil {
				if err != io.EOF {
					t.client.logger.Error("Error membaca data: %v", err)
				}
				return
			}

			// Kirim data ke server
			if err := t.sendData(connID, buffer[:n]); err != nil {
				t.client.logger.Error("Error mengirim data: %v", err)
				return
			}
		}
	}
}

// sendData mengirim data ke server
func (t *Tunnel) sendData(connID string, data []byte) error {
	// Buat payload data
	dataPayload := proto.DataPayload{
		TunnelID:     t.ID,
		ConnectionID: connID,
		Data:         data,
	}

	// Buat pesan data
	dataMsg, err := proto.NewMessage(proto.MessageTypeData, dataPayload)
	if err != nil {
		return fmt.Errorf("error membuat pesan data: %v", err)
	}

	// Kirim pesan data
	if err := t.client.sendMessage(dataMsg); err != nil {
		return fmt.Errorf("error mengirim pesan data: %v", err)
	}

	return nil
}

// HandleData menangani data yang diterima dari server
func (t *Tunnel) HandleData(connID string, data []byte) error {
	// Dapatkan koneksi
	t.connectionsMu.RLock()
	conn, exists := t.connections[connID]
	t.connectionsMu.RUnlock()

	if !exists {
		return fmt.Errorf("koneksi %s tidak ditemukan", connID)
	}

	// Kirim data ke koneksi lokal
	_, err := conn.Write(data)
	if err != nil {
		return fmt.Errorf("error menulis data: %v", err)
	}

	return nil
}

// isClosedError memeriksa apakah error disebabkan oleh koneksi yang ditutup
func isClosedError(err error) bool {
	if err == nil {
		return false
	}
	return err.Error() == "use of closed network connection"
}

// generateRandomString menghasilkan string acak dengan panjang tertentu
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[i%len(charset)]
	}
	return string(result)
}
