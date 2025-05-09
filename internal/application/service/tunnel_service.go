package service

import (
	"fmt"

	"github.com/alwanandri2712/haxorport-go-client/internal/domain/model"
	"github.com/alwanandri2712/haxorport-go-client/internal/domain/port"
)

// TunnelService adalah service untuk mengelola tunnel
type TunnelService struct {
	tunnelRepo port.TunnelRepository
	logger     port.Logger
}

// NewTunnelService membuat instance TunnelService baru
func NewTunnelService(tunnelRepo port.TunnelRepository, logger port.Logger) *TunnelService {
	return &TunnelService{
		tunnelRepo: tunnelRepo,
		logger:     logger,
	}
}

// CreateHTTPTunnel membuat tunnel HTTP baru
func (s *TunnelService) CreateHTTPTunnel(localPort int, subdomain string, auth *model.TunnelAuth) (*model.Tunnel, error) {
	s.logger.Info("Membuat tunnel HTTP untuk port lokal %d dengan subdomain %s", localPort, subdomain)

	// Buat konfigurasi tunnel
	tunnelConfig := model.TunnelConfig{
		Type:      model.TunnelTypeHTTP,
		LocalPort: localPort,
		Subdomain: subdomain,
		Auth:      auth,
	}

	// Daftarkan tunnel
	tunnel, err := s.tunnelRepo.Register(tunnelConfig)
	if err != nil {
		return nil, fmt.Errorf("gagal mendaftarkan tunnel HTTP: %v", err)
	}

	s.logger.Info("Tunnel HTTP berhasil dibuat dengan URL: %s", tunnel.URL)
	// Tidak menampilkan statistik tunnel

	return tunnel, nil
}

// CreateTCPTunnel membuat tunnel TCP baru
func (s *TunnelService) CreateTCPTunnel(localPort int, remotePort int) (*model.Tunnel, error) {
	s.logger.Info("Membuat tunnel TCP untuk port lokal %d dengan port remote %d", localPort, remotePort)

	// Buat konfigurasi tunnel
	tunnelConfig := model.TunnelConfig{
		Type:       model.TunnelTypeTCP,
		LocalPort:  localPort,
		RemotePort: remotePort,
	}

	// Daftarkan tunnel
	tunnel, err := s.tunnelRepo.Register(tunnelConfig)
	if err != nil {
		return nil, fmt.Errorf("gagal mendaftarkan tunnel TCP: %v", err)
	}

	s.logger.Info("Tunnel TCP berhasil dibuat dengan port remote: %d", tunnel.RemotePort)
	// Tidak menampilkan statistik tunnel

	return tunnel, nil
}

// CloseTunnel menutup tunnel
func (s *TunnelService) CloseTunnel(tunnelID string) error {
	s.logger.Info("Menutup tunnel dengan ID: %s", tunnelID)

	// Dapatkan tunnel
	tunnel, err := s.tunnelRepo.GetByID(tunnelID)
	if err != nil {
		return fmt.Errorf("tunnel tidak ditemukan: %v", err)
	}

	s.logger.Info("Menutup tunnel %s dengan tipe %s", tunnelID, tunnel.Config.Type)

	// Hapus tunnel
	if err := s.tunnelRepo.Unregister(tunnelID); err != nil {
		return fmt.Errorf("gagal menghapus tunnel: %v", err)
	}

	s.logger.Info("Tunnel berhasil ditutup: %s", tunnelID)

	return nil
}

// GetAllTunnels mengembalikan semua tunnel yang aktif
func (s *TunnelService) GetAllTunnels() []*model.Tunnel {
	return s.tunnelRepo.GetAll()
}

// GetTunnelByID mengembalikan tunnel berdasarkan ID
func (s *TunnelService) GetTunnelByID(tunnelID string) (*model.Tunnel, error) {
	return s.tunnelRepo.GetByID(tunnelID)
}
