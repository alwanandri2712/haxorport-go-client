package config

import (
        "fmt"
        "os"
        "path/filepath"
        "strings"

        "github.com/spf13/viper"
)

// Config adalah struktur konfigurasi klien haxorport
type Config struct {
        // ServerAddress adalah alamat server haxorport
        ServerAddress string `mapstructure:"server_address" yaml:"server_address"`
        // ControlPort adalah port untuk control plane
        ControlPort int `mapstructure:"control_port" yaml:"control_port"`
        // AuthToken adalah token untuk autentikasi ke server
        AuthToken string `mapstructure:"auth_token" yaml:"auth_token"`
        // LogLevel adalah level logging (debug, info, warn, error)
        LogLevel string `mapstructure:"log_level" yaml:"log_level"`
        // LogFile adalah path ke file log (kosong untuk stdout)
        LogFile string `mapstructure:"log_file" yaml:"log_file"`
        // Tunnels adalah daftar tunnel yang akan dibuat saat startup
        Tunnels []TunnelConfig `mapstructure:"tunnels" yaml:"tunnels"`
}

// TunnelConfig adalah konfigurasi untuk tunnel
type TunnelConfig struct {
        // Name adalah nama tunnel (opsional)
        Name string `mapstructure:"name" yaml:"name"`
        // Type adalah tipe tunnel (http, tcp)
        Type string `mapstructure:"type" yaml:"type"`
        // LocalPort adalah port lokal yang akan di-tunnel
        LocalPort int `mapstructure:"local_port" yaml:"local_port"`
        // Subdomain adalah subdomain yang diminta (untuk HTTP, opsional)
        Subdomain string `mapstructure:"subdomain" yaml:"subdomain,omitempty"`
        // RemotePort adalah port remote yang diminta (untuk TCP, opsional)
        RemotePort int `mapstructure:"remote_port" yaml:"remote_port,omitempty"`
        // Auth adalah konfigurasi autentikasi untuk tunnel (opsional)
        Auth *AuthConfig `mapstructure:"auth" yaml:"auth,omitempty"`
}

// AuthConfig adalah konfigurasi autentikasi untuk tunnel
type AuthConfig struct {
        // Type adalah tipe autentikasi (basic, header)
        Type string `mapstructure:"type" yaml:"type"`
        // Username adalah username untuk autentikasi basic
        Username string `mapstructure:"username" yaml:"username,omitempty"`
        // Password adalah password untuk autentikasi basic
        Password string `mapstructure:"password" yaml:"password,omitempty"`
        // HeaderName adalah nama header untuk autentikasi header
        HeaderName string `mapstructure:"header_name" yaml:"header_name,omitempty"`
        // HeaderValue adalah nilai header untuk autentikasi header
        HeaderValue string `mapstructure:"header_value" yaml:"header_value,omitempty"`
}

// DefaultConfig mengembalikan konfigurasi default
func DefaultConfig() *Config {
        return &Config{
                ServerAddress: "localhost",
                ControlPort:   8080,
                LogLevel:      "info",
                LogFile:       "",
                Tunnels:       []TunnelConfig{},
        }
}

// LoadConfig memuat konfigurasi dari file
func LoadConfig(configPath string) (*Config, error) {
        config := DefaultConfig()

        // Jika configPath kosong, cari di lokasi default
        if configPath == "" {
                // Cari di direktori home pengguna
                homeDir, err := os.UserHomeDir()
                if err == nil {
                        // Cek ~/.haxorport/config.yaml
                        homePath := filepath.Join(homeDir, ".haxorport", "config.yaml")
                        if _, err := os.Stat(homePath); err == nil {
                                configPath = homePath
                        }
                }

                // Jika masih kosong, cek di direktori saat ini
                if configPath == "" {
                        if _, err := os.Stat("config.yaml"); err == nil {
                                configPath = "config.yaml"
                        }
                }
        }

        // Jika configPath masih kosong, gunakan konfigurasi default
        if configPath == "" {
                return config, nil
        }

        // Muat konfigurasi dari file
        viper.SetConfigFile(configPath)
        if err := viper.ReadInConfig(); err != nil {
                return nil, fmt.Errorf("error membaca file konfigurasi: %v", err)
        }

        // Unmarshal ke struct Config
        if err := viper.Unmarshal(config); err != nil {
                return nil, fmt.Errorf("error parsing konfigurasi: %v", err)
        }

        return config, nil
}

// SaveConfig menyimpan konfigurasi ke file
func SaveConfig(config *Config, configPath string) error {
        // Jika configPath kosong, gunakan lokasi default
        if configPath == "" {
                homeDir, err := os.UserHomeDir()
                if err != nil {
                        return fmt.Errorf("error mendapatkan direktori home: %v", err)
                }

                // Buat direktori ~/.haxorport jika belum ada
                configDir := filepath.Join(homeDir, ".haxorport")
                if err := os.MkdirAll(configDir, 0755); err != nil {
                        return fmt.Errorf("error membuat direktori konfigurasi: %v", err)
                }

                configPath = filepath.Join(configDir, "config.yaml")
        }

        // Set konfigurasi di viper
        viper.SetConfigFile(configPath)
        
        // Set nilai konfigurasi di viper
        viper.Set("server_address", config.ServerAddress)
        viper.Set("control_port", config.ControlPort)
        viper.Set("auth_token", config.AuthToken)
        viper.Set("log_level", config.LogLevel)
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
