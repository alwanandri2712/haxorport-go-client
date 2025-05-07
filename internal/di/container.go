package di

import (
	"os"

	"github.com/haxorport/haxor-client/internal/application/service"
	"github.com/haxorport/haxor-client/internal/domain/model"
	"github.com/haxorport/haxor-client/internal/infrastructure/config"
	"github.com/haxorport/haxor-client/internal/infrastructure/logger"
	"github.com/haxorport/haxor-client/internal/infrastructure/transport"
)

// Container adalah container untuk dependency injection
type Container struct {
	// Logger
	Logger *logger.Logger

	// Repositories
	ConfigRepository *config.ConfigRepository

	// Services
	ConfigService *service.ConfigService
	TunnelService *service.TunnelService

	// Client
	Client *transport.Client

	// TunnelRepository
	TunnelRepository *transport.TunnelRepository

	// Config
	Config *model.Config
}

// NewContainer membuat instance Container baru
func NewContainer() *Container {
	return &Container{}
}

// Initialize menginisialisasi container
func (c *Container) Initialize(configPath string) error {
	// Inisialisasi logger
	c.Logger = logger.NewLogger(os.Stdout, "info")

	// Inisialisasi config repository
	c.ConfigRepository = config.NewConfigRepository()

	// Inisialisasi config service
	c.ConfigService = service.NewConfigService(c.ConfigRepository, c.Logger)

	// Muat konfigurasi
	var err error
	c.Config, err = c.ConfigService.LoadConfig(configPath)
	if err != nil {
		return err
	}

	// Setel level logger berdasarkan konfigurasi
	c.Logger.SetLevel(string(c.Config.LogLevel))

	// Jika file log ditentukan, gunakan file logger tetapi tetap tampilkan output ke terminal
	if c.Config.LogFile != "" {
		_, err := logger.NewFileLogger(c.Config.LogFile, string(c.Config.LogLevel))
		if err != nil {
			c.Logger.Error("Gagal membuat file logger: %v", err)
		} else {
			// Tetap gunakan logger yang ada (ke stdout)
			c.Logger.Info("Log juga akan ditulis ke file: %s", c.Config.LogFile)
		}
	}

	// Inisialisasi client
	c.Client = transport.NewClient(c.Config, c.Logger)

	// Inisialisasi tunnel repository
	c.TunnelRepository = transport.NewTunnelRepository(c.Client, c.Logger)

	// Inisialisasi tunnel service
	c.TunnelService = service.NewTunnelService(c.TunnelRepository, c.Logger)

	// Daftarkan handler untuk pesan HTTP request
	c.Client.RegisterHandler(model.MessageTypeHTTPRequest, c.Client.HandleHTTPRequestMessage)

	return nil
}

// Close menutup semua resource
func (c *Container) Close() {
	// Tutup client
	if c.Client != nil {
		c.Client.Close()
	}

	// Tutup logger
	if c.Logger != nil {
		c.Logger.Close()
	}
}
