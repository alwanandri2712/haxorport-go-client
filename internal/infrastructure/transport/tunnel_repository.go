package transport

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/haxorport/haxor-client/internal/domain/model"
	"github.com/haxorport/haxor-client/internal/domain/port"
)

// TunnelRepository adalah implementasi port.TunnelRepository
type TunnelRepository struct {
	client      *Client
	logger      port.Logger
	tunnels     map[string]*model.Tunnel
	connections map[string]net.Conn
	mutex       sync.RWMutex
}

// NewTunnelRepository membuat instance TunnelRepository baru
func NewTunnelRepository(client *Client, logger port.Logger) *TunnelRepository {
	repo := &TunnelRepository{
		client:      client,
		logger:      logger,
		tunnels:     make(map[string]*model.Tunnel),
		connections: make(map[string]net.Conn),
		mutex:       sync.RWMutex{},
	}

	// Daftarkan handler untuk pesan data
	client.RegisterHandler(model.MessageTypeData, repo.handleDataMessage)

	return repo
}

// Register mendaftarkan tunnel baru ke server
func (r *TunnelRepository) Register(config model.TunnelConfig) (*model.Tunnel, error) {
	// Pastikan klien terhubung
	if !r.client.IsConnected() {
		if err := r.client.Connect(); err != nil {
			return nil, fmt.Errorf("gagal terhubung ke server: %v", err)
		}
	}

	// Kirim permintaan pendaftaran tunnel
	response, err := r.client.SendRegisterTunnel(config)
	if err != nil {
		return nil, fmt.Errorf("gagal mendaftarkan tunnel: %v", err)
	}

	// Periksa respons
	if !response.Success {
		return nil, fmt.Errorf("pendaftaran tunnel gagal: %s", response.Error)
	}

	// Buat objek tunnel
	tunnel := model.NewTunnel(response.TunnelID, config)

	// Set informasi tunnel berdasarkan tipe
	if config.Type == model.TunnelTypeHTTP {
		tunnel.SetHTTPInfo(response.URL)
	} else if config.Type == model.TunnelTypeTCP {
		tunnel.SetTCPInfo(response.RemotePort)
	}

	// Simpan tunnel
	r.mutex.Lock()
	r.tunnels[response.TunnelID] = tunnel
	r.mutex.Unlock()

	// Mulai listener untuk tunnel tanpa menampilkan statistik
	go r.startTunnelListener(tunnel)

	return tunnel, nil
}

// Unregister menghapus tunnel dari server
func (r *TunnelRepository) Unregister(tunnelID string) error {
	// Pastikan klien terhubung
	if !r.client.IsConnected() {
		if err := r.client.Connect(); err != nil {
			return fmt.Errorf("gagal terhubung ke server: %v", err)
		}
	}

	// Kirim permintaan penghapusan tunnel
	if err := r.client.SendUnregisterTunnel(tunnelID); err != nil {
		return fmt.Errorf("gagal menghapus tunnel: %v", err)
	}

	// Hapus tunnel dari map
	r.mutex.Lock()
	delete(r.tunnels, tunnelID)
	r.mutex.Unlock()

	return nil
}

// GetAll mengembalikan semua tunnel yang aktif
func (r *TunnelRepository) GetAll() []*model.Tunnel {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	tunnels := make([]*model.Tunnel, 0, len(r.tunnels))
	for _, tunnel := range r.tunnels {
		tunnels = append(tunnels, tunnel)
	}

	return tunnels
}

// GetByID mengembalikan tunnel berdasarkan ID
func (r *TunnelRepository) GetByID(tunnelID string) (*model.Tunnel, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	tunnel, exists := r.tunnels[tunnelID]
	if !exists {
		return nil, fmt.Errorf("tunnel dengan ID %s tidak ditemukan", tunnelID)
	}

	return tunnel, nil
}

// SendData mengirim data melalui tunnel
func (r *TunnelRepository) SendData(tunnelID string, connectionID string, data []byte) error {
	return r.client.SendData(tunnelID, connectionID, data)
}

// HandleData menangani data yang diterima dari server
func (r *TunnelRepository) HandleData(tunnelID string, connectionID string, data []byte) error {
	r.mutex.RLock()
	conn, exists := r.connections[connectionID]
	r.mutex.RUnlock()

	if !exists {
		return fmt.Errorf("koneksi dengan ID %s tidak ditemukan", connectionID)
	}

	// Kirim data ke koneksi lokal
	_, err := conn.Write(data)
	if err != nil {
		r.logger.Error("Gagal mengirim data ke koneksi lokal: %v", err)
		return err
	}

	return nil
}

// handleDataMessage menangani pesan data dari server
func (r *TunnelRepository) handleDataMessage(msg *model.Message) error {
	// Parse payload
	var payload model.DataPayload
	if err := msg.ParsePayload(&payload); err != nil {
		return fmt.Errorf("gagal mengurai payload data: %v", err)
	}

	// Tangani data
	return r.HandleData(payload.TunnelID, payload.ConnectionID, payload.Data)
}

// startTunnelListener memulai listener untuk tunnel
func (r *TunnelRepository) startTunnelListener(tunnel *model.Tunnel) {
	localAddr := fmt.Sprintf("localhost:%d", tunnel.Config.LocalPort)
	r.logger.Info("Memulai listener untuk tunnel %s pada %s", tunnel.ID, localAddr)

	// Buat listener
	listener, err := net.Listen("tcp", localAddr)
	if err != nil {
		r.logger.Error("Gagal membuat listener untuk tunnel %s: %v", tunnel.ID, err)
		return
	}
	defer listener.Close()

	for {
		// Terima koneksi
		conn, err := listener.Accept()
		if err != nil {
			r.logger.Error("Gagal menerima koneksi untuk tunnel %s: %v", tunnel.ID, err)
			break
		}

		// Buat ID koneksi
		connectionID := fmt.Sprintf("%s-%d", tunnel.ID, time.Now().UnixNano())

		// Simpan koneksi
		r.mutex.Lock()
		r.connections[connectionID] = conn
		r.mutex.Unlock()

		// Tangani koneksi
		go r.handleConnection(tunnel.ID, connectionID, conn)
	}
}

// handleConnection menangani koneksi untuk tunnel
func (r *TunnelRepository) handleConnection(tunnelID string, connectionID string, conn net.Conn) {
	defer func() {
		conn.Close()
		r.mutex.Lock()
		delete(r.connections, connectionID)
		r.mutex.Unlock()
	}()

	buffer := make([]byte, 4096)
	for {
		// Baca data dari koneksi lokal
		n, err := conn.Read(buffer)
		if err != nil {
			if err != io.EOF {
				r.logger.Error("Gagal membaca data dari koneksi lokal: %v", err)
			}
			break
		}

		// Kirim data ke server
		if err := r.SendData(tunnelID, connectionID, buffer[:n]); err != nil {
			r.logger.Error("Gagal mengirim data ke server: %v", err)
			break
		}
	}
}

// Ensure TunnelRepository implements port.TunnelRepository
var _ port.TunnelRepository = (*TunnelRepository)(nil)
