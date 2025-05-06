package model

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
	// AuthToken adalah token untuk autentikasi ke server
	AuthToken string
	// LogLevel adalah level logging (debug, info, warn, error)
	LogLevel LogLevel
	// LogFile adalah path ke file log (kosong untuk stdout)
	LogFile string
	// Tunnels adalah daftar tunnel yang akan dibuat saat startup
	Tunnels []TunnelConfig
}

// NewConfig membuat instance Config baru dengan nilai default
func NewConfig() *Config {
	return &Config{
		ServerAddress: "control.haxorport.online",
		ControlPort:   443,
		LogLevel:      LogLevelInfo,
		LogFile:       "",
		Tunnels:       []TunnelConfig{},
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
