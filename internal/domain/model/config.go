package model

import (
	"os"
	"path/filepath"
)

// LogLevel mendefinisikan level logging
type LogLevel string

const (
	// LogLevelDebug adalah level untuk pesan debug
	LogLevelDebug LogLevel = "debug"
	// LogLevelInfo adalah level untuk pesan informasi
	LogLevelInfo LogLevel = "info"
	// LogLevelWarn adalah level untuk pesan peringatan
	LogLevelWarn LogLevel = "warn"
	// LogLevelError adalah level untuk pesan error
	LogLevelError LogLevel = "error"
)

// Config adalah struktur konfigurasi klien haxorport
type Config struct {
	// ServerAddress adalah alamat server haxorport
	ServerAddress string
	// ControlPort adalah port untuk control plane
	ControlPort int
	// DataPort adalah port untuk data plane
	DataPort int
	// AuthEnabled adalah flag untuk mengaktifkan autentikasi
	AuthEnabled bool
	// AuthToken adalah token untuk autentikasi ke server
	AuthToken string
	// AuthValidationURL adalah URL untuk validasi token (kosong untuk menggunakan default)
	AuthValidationURL string
	// TLSEnabled adalah flag untuk mengaktifkan TLS
	TLSEnabled bool
	// TLSCert adalah path ke file sertifikat TLS
	TLSCert string
	// TLSKey adalah path ke file kunci TLS
	TLSKey string
	// LogLevel adalah level logging (debug, info, warn, error)
	LogLevel LogLevel
	// LogFile adalah path ke file log (kosong untuk stdout)
	LogFile string
	// BaseDomain adalah domain dasar untuk subdomain tunnel
	BaseDomain string
	// Tunnels adalah daftar tunnel yang akan dibuat saat startup
	Tunnels []TunnelConfig
}

// NewConfig membuat instance Config baru dengan nilai default
func NewConfig() *Config {
	return &Config{
		ServerAddress:     "control.haxorport.online",
		ControlPort:       443,
		DataPort:          8081,
		AuthEnabled:       false,
		AuthToken:         "",
		AuthValidationURL: "https://haxorport.online/AuthToken/validate",
		TLSEnabled:        false,
		TLSCert:           "",
		TLSKey:            "",
		LogLevel:          LogLevelWarn,
		LogFile:           "",
		BaseDomain:        "haxorport.online",
		Tunnels:           []TunnelConfig{},
	}
}

// AddTunnel menambahkan tunnel ke konfigurasi
func (c *Config) AddTunnel(tunnel TunnelConfig) {
	c.Tunnels = append(c.Tunnels, tunnel)
}

// RemoveTunnel menghapus tunnel dari konfigurasi berdasarkan nama
func (c *Config) RemoveTunnel(name string) bool {
	for i, tunnel := range c.Tunnels {
		if tunnel.Name == name {
			// Hapus tunnel dari slice
			c.Tunnels = append(c.Tunnels[:i], c.Tunnels[i+1:]...)
			return true
		}
	}
	return false
}

// GetTunnel mengembalikan tunnel berdasarkan nama
func (c *Config) GetTunnel(name string) *TunnelConfig {
	for _, tunnel := range c.Tunnels {
		if tunnel.Name == name {
			return &tunnel
		}
	}
	return nil
}

// GetConfigFilePath mengembalikan path ke file konfigurasi
func (c *Config) GetConfigFilePath() string {
	// Tentukan direktori konfigurasi berdasarkan user
	configDir := "/etc/haxorport"
	
	// Jika bukan root, gunakan direktori home
	if os.Getuid() != 0 {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			configDir = filepath.Join(homeDir, ".haxorport")
		}
	}
	
	// Path file konfigurasi
	return filepath.Join(configDir, "config.yaml")
}
