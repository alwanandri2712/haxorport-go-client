package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/haxorport/haxor-client/internal/domain/model"
	"github.com/haxorport/haxor-client/internal/domain/port"
	"github.com/spf13/viper"
)

// ConfigRepository adalah implementasi port.ConfigRepository
type ConfigRepository struct{}

// NewConfigRepository membuat instance ConfigRepository baru
func NewConfigRepository() *ConfigRepository {
	return &ConfigRepository{}
}

// Load memuat konfigurasi dari file
func (r *ConfigRepository) Load(configPath string) (*model.Config, error) {
	config := model.NewConfig()

	// Jika configPath kosong, cari di lokasi default
	if configPath == "" {
		var err error
		configPath, err = r.GetDefaultPath()
		if err != nil {
			return nil, err
		}
	}

	// Periksa apakah file ada
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return config, nil
	}

	// Muat konfigurasi dari file
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error membaca file konfigurasi: %v", err)
	}

	// Mapping dari viper ke struct Config
	config.ServerAddress = viper.GetString("server_address")
	config.ControlPort = viper.GetInt("control_port")
	config.DataPort = viper.GetInt("data_port")
	config.AuthEnabled = viper.GetBool("auth_enabled")
	config.AuthToken = viper.GetString("auth_token")
	config.TLSEnabled = viper.GetBool("tls_enabled")
	config.TLSCert = viper.GetString("tls_cert")
	config.TLSKey = viper.GetString("tls_key")
	config.BaseDomain = viper.GetString("base_domain")
	config.LogLevel = model.LogLevel(viper.GetString("log_level"))
	config.LogFile = viper.GetString("log_file")

	// Muat tunnel
	var tunnelConfigs []model.TunnelConfig
	if err := viper.UnmarshalKey("tunnels", &tunnelConfigs); err != nil {
		return nil, fmt.Errorf("error parsing konfigurasi tunnel: %v", err)
	}
	config.Tunnels = tunnelConfigs

	return config, nil
}

// Save menyimpan konfigurasi ke file
func (r *ConfigRepository) Save(config *model.Config, configPath string) error {
	// Jika configPath kosong, gunakan lokasi default
	if configPath == "" {
		var err error
		configPath, err = r.GetDefaultPath()
		if err != nil {
			return err
		}
	}

	// Buat direktori jika belum ada
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("error membuat direktori konfigurasi: %v", err)
	}

	// Set konfigurasi di viper
	viper.SetConfigFile(configPath)
	
	// Set nilai konfigurasi di viper
	viper.Set("server_address", config.ServerAddress)
	viper.Set("control_port", config.ControlPort)
	viper.Set("data_port", config.DataPort)
	viper.Set("auth_enabled", config.AuthEnabled)
	viper.Set("auth_token", config.AuthToken)
	viper.Set("tls_enabled", config.TLSEnabled)
	viper.Set("tls_cert", config.TLSCert)
	viper.Set("tls_key", config.TLSKey)
	viper.Set("base_domain", config.BaseDomain)
	viper.Set("log_level", string(config.LogLevel))
	viper.Set("log_file", config.LogFile)
	viper.Set("tunnels", config.Tunnels)

	// Simpan ke file
	if err := viper.WriteConfig(); err != nil {
		// Jika file tidak ada, buat baru
		if strings.Contains(err.Error(), "no such file") {
			return viper.SafeWriteConfig()
		}
		return fmt.Errorf("error menyimpan konfigurasi: %v", err)
	}

	return nil
}

// GetDefaultPath mengembalikan path default untuk file konfigurasi
func (r *ConfigRepository) GetDefaultPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error mendapatkan direktori home: %v", err)
	}

	return filepath.Join(homeDir, ".haxorport", "config.yaml"), nil
}

// Ensure ConfigRepository implements port.ConfigRepository
var _ port.ConfigRepository = (*ConfigRepository)(nil)
