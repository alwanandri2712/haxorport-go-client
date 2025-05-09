package service

import (
	"fmt"

	"github.com/alwanandri2712/haxorport-go-client/internal/domain/model"
	"github.com/alwanandri2712/haxorport-go-client/internal/domain/port"
)

// ConfigService adalah service untuk mengelola konfigurasi
type ConfigService struct {
	configRepo port.ConfigRepository
	logger     port.Logger
}

// NewConfigService membuat instance ConfigService baru
func NewConfigService(configRepo port.ConfigRepository, logger port.Logger) *ConfigService {
	return &ConfigService{
		configRepo: configRepo,
		logger:     logger,
	}
}

// LoadConfig memuat konfigurasi dari file
func (s *ConfigService) LoadConfig(configPath string) (*model.Config, error) {
	// Jika configPath kosong, gunakan path default
	if configPath == "" {
		var err error
		configPath, err = s.configRepo.GetDefaultPath()
		if err != nil {
			return nil, fmt.Errorf("gagal mendapatkan path default: %v", err)
		}
	}
	
	// Muat konfigurasi
	config, err := s.configRepo.Load(configPath)
	if err != nil {
		s.logger.Warn("Gagal memuat konfigurasi dari %s: %v", configPath, err)
		// Kembalikan konfigurasi default jika gagal memuat
		return model.NewConfig(), nil
	}
	
	s.logger.Info("Konfigurasi dimuat dari %s", configPath)
	
	return config, nil
}

// SaveConfig menyimpan konfigurasi ke file
func (s *ConfigService) SaveConfig(config *model.Config, configPath string) error {
	// Jika configPath kosong, gunakan path default
	if configPath == "" {
		var err error
		configPath, err = s.configRepo.GetDefaultPath()
		if err != nil {
			return fmt.Errorf("gagal mendapatkan path default: %v", err)
		}
	}
	
	// Simpan konfigurasi
	if err := s.configRepo.Save(config, configPath); err != nil {
		return fmt.Errorf("gagal menyimpan konfigurasi: %v", err)
	}
	
	s.logger.Info("Konfigurasi disimpan ke %s", configPath)
	
	return nil
}

// SetServerAddress mengatur alamat server
func (s *ConfigService) SetServerAddress(config *model.Config, serverAddress string) {
	config.ServerAddress = serverAddress
}

// SetControlPort mengatur port control plane
func (s *ConfigService) SetControlPort(config *model.Config, controlPort int) {
	config.ControlPort = controlPort
}

// SetAuthToken mengatur token autentikasi
func (s *ConfigService) SetAuthToken(config *model.Config, authToken string) {
	config.AuthToken = authToken
}

// SetLogLevel mengatur level logging
func (s *ConfigService) SetLogLevel(config *model.Config, logLevel string) {
	config.LogLevel = model.LogLevel(logLevel)
}

// SetLogFile mengatur file log
func (s *ConfigService) SetLogFile(config *model.Config, logFile string) {
	config.LogFile = logFile
}

// AddTunnel menambahkan tunnel ke konfigurasi
func (s *ConfigService) AddTunnel(config *model.Config, tunnel model.TunnelConfig) {
	config.AddTunnel(tunnel)
}

// RemoveTunnel menghapus tunnel dari konfigurasi
func (s *ConfigService) RemoveTunnel(config *model.Config, name string) bool {
	return config.RemoveTunnel(name)
}

// GetTunnel mengembalikan tunnel dari konfigurasi
func (s *ConfigService) GetTunnel(config *model.Config, name string) *model.TunnelConfig {
	return config.GetTunnel(name)
}
