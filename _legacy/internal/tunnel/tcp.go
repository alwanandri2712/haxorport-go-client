package tunnel

import (
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/haxorport/client/internal/config"
)

// TCPTunnel adalah tunnel untuk protokol TCP
type TCPTunnel struct {
	*Tunnel
	remotePort int
}

// NewTCPTunnel membuat instance TCPTunnel baru
func NewTCPTunnel(id string, config config.TunnelConfig, client *Client) *TCPTunnel {
	tunnel := NewTunnel(id, config, client)
	return &TCPTunnel{
		Tunnel:     tunnel,
		remotePort: config.RemotePort,
	}
}

// Start memulai tunnel TCP
func (t *TCPTunnel) Start() error {
	// Buat listener untuk menerima koneksi lokal
	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", t.Config.LocalPort))
	if err != nil {
		return fmt.Errorf("error membuat listener: %v", err)
	}

	t.listener = listener
	t.client.logger.Info("Tunnel TCP %s mendengarkan di localhost:%d -> remote:%d", t.ID, t.Config.LocalPort, t.remotePort)

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

// handleConnection menangani koneksi TCP
func (t *TCPTunnel) handleConnection(conn net.Conn) {
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

	t.client.logger.Debug("Koneksi TCP baru untuk tunnel %s: %s", t.ID, connID)

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

// ProxyTCP membuat proxy TCP antara koneksi lokal dan remote
func ProxyTCP(local, remote net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	
	// Copy data dari local ke remote
	_, err := io.Copy(remote, local)
	if err != nil && err != io.EOF {
		// Ignore error
	}
}
