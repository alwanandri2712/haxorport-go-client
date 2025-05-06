package main

import (
        "fmt"
        "os"
        "os/signal"
        "path/filepath"
        "syscall"

        "github.com/haxorport/client/internal/config"
        "github.com/haxorport/client/internal/logger"
        "github.com/haxorport/client/internal/tunnel"
        "github.com/spf13/cobra"
)

var (
        // Versi klien
        version = "0.1.0"
        
        // Flag global
        configPath string
        serverAddr string
        serverPort int
        authToken  string
        logLevel   string
        logFile    string
        
        // Config
        cfg *config.Config
        log *logger.Logger
        
        // Client
        client *tunnel.Client
)

// rootCmd adalah perintah root untuk CLI
var rootCmd = &cobra.Command{
        Use:   "hxp",
        Short: "Haxorport - Tunneling HTTP dan TCP",
        Long: `Haxorport adalah platform tunneling yang memungkinkan pengembang
mengekspos layanan lokal ke internet melalui tunnel yang aman.

Haxorport mendukung tunneling HTTP dan TCP, dengan fitur-fitur seperti:
- Subdomain kustom untuk tunnel HTTP
- Port kustom untuk tunnel TCP
- Autentikasi basic dan header untuk tunnel HTTP
- Konfigurasi yang dapat disimpan untuk penggunaan berulang
- Manajemen tunnel yang mudah`,
        Example: `  # Membuat tunnel HTTP untuk aplikasi web di port 8080
  hxp http 8080 --subdomain myapp

  # Membuat tunnel HTTP dengan autentikasi basic
  hxp http 8080 --subdomain myapp --auth-user admin --auth-pass secret

  # Membuat tunnel TCP untuk server SSH di port 22
  hxp tcp 22 --remote-port 2222

  # Menggunakan file konfigurasi kustom
  hxp --config /path/to/config.yaml http 8080

  # Melihat status tunnel yang aktif
  hxp status

  # Menampilkan konfigurasi saat ini
  hxp config show`,
        Version: version,
        PersistentPreRun: func(cmd *cobra.Command, args []string) {
                // Muat konfigurasi
                var err error
                cfg, err = config.LoadConfig(configPath)
                if err != nil {
                        fmt.Printf("Error memuat konfigurasi: %v\n", err)
                        os.Exit(1)
                }
                
                // Override konfigurasi dengan flag
                if serverAddr != "" {
                        cfg.ServerAddress = serverAddr
                }
                if serverPort != 0 {
                        cfg.ControlPort = serverPort
                }
                if authToken != "" {
                        cfg.AuthToken = authToken
                }
                if logLevel != "" {
                        cfg.LogLevel = logLevel
                }
                if logFile != "" {
                        cfg.LogFile = logFile
                }
                
                // Inisialisasi logger
                if cfg.LogFile != "" {
                        // Pastikan direktori log ada
                        logDir := filepath.Dir(cfg.LogFile)
                        if err := os.MkdirAll(logDir, 0755); err != nil {
                                fmt.Printf("Error membuat direktori log: %v\n", err)
                                os.Exit(1)
                        }
                        
                        // Buka file log
                        logFile, err := os.OpenFile(cfg.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
                        if err != nil {
                                fmt.Printf("Error membuka file log: %v\n", err)
                                os.Exit(1)
                        }
                        
                        log = logger.NewLogger(logFile, cfg.LogLevel)
                } else {
                        log = logger.NewLogger(os.Stdout, cfg.LogLevel)
                }
                
                // Inisialisasi klien
                client = tunnel.NewClient(cfg, log)
        },
}

// httpCmd adalah perintah untuk membuat tunnel HTTP
var httpCmd = &cobra.Command{
        Use:   "http [local_port]",
        Short: "Membuat tunnel HTTP untuk mengekspos layanan web lokal ke internet",
        Long: `Perintah 'http' membuat tunnel HTTP yang mengekspos layanan web lokal ke internet.
Anda dapat menentukan subdomain kustom dan menambahkan autentikasi untuk mengamankan akses.

Tunnel HTTP memungkinkan Anda:
- Mengekspos aplikasi web lokal ke internet dengan URL publik
- Menguji webhook dan integrasi dengan layanan eksternal
- Berbagi aplikasi web lokal dengan orang lain tanpa deployment
- Mengamankan akses dengan autentikasi basic atau header kustom`,
        Example: `  # Membuat tunnel HTTP sederhana untuk aplikasi web di port 8080
  hxp http 8080

  # Menggunakan subdomain kustom
  hxp http 8080 --subdomain myapp

  # Menambahkan autentikasi basic
  hxp http 8080 --subdomain myapp --auth-user admin --auth-pass secret

  # Menggunakan autentikasi header kustom
  hxp http 8080 --auth-header-name "X-API-Key" --auth-header-value "secret-key"`,
        Args:  cobra.ExactArgs(1),
        Run: func(cmd *cobra.Command, args []string) {
                // Parse port lokal
                var localPort int
                if _, err := fmt.Sscanf(args[0], "%d", &localPort); err != nil {
                        log.Error("Port lokal tidak valid: %v", err)
                        os.Exit(1)
                }
                
                // Dapatkan flag
                subdomain, _ := cmd.Flags().GetString("subdomain")
                username, _ := cmd.Flags().GetString("auth-user")
                password, _ := cmd.Flags().GetString("auth-pass")
                headerName, _ := cmd.Flags().GetString("auth-header-name")
                headerValue, _ := cmd.Flags().GetString("auth-header-value")
                
                // Buat konfigurasi tunnel
                tunnelConfig := config.TunnelConfig{
                        Type:      "http",
                        LocalPort: localPort,
                        Subdomain: subdomain,
                }
                
                // Tambahkan autentikasi jika diperlukan
                if username != "" && password != "" {
                        tunnelConfig.Auth = &config.AuthConfig{
                                Type:     "basic",
                                Username: username,
                                Password: password,
                        }
                } else if headerName != "" && headerValue != "" {
                        tunnelConfig.Auth = &config.AuthConfig{
                                Type:        "header",
                                HeaderName:  headerName,
                                HeaderValue: headerValue,
                        }
                }
                
                // Hubungkan ke server
                if err := client.Connect(); err != nil {
                        log.Error("Gagal menghubungkan ke server: %v", err)
                        os.Exit(1)
                }
                
                // Daftarkan tunnel
                tunnel, err := client.RegisterTunnel(tunnelConfig)
                if err != nil {
                        log.Error("Gagal mendaftarkan tunnel: %v", err)
                        os.Exit(1)
                }
                
                // Tangkap sinyal untuk shutdown yang bersih
                sigChan := make(chan os.Signal, 1)
                signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
                
                log.Info("Tekan Ctrl+C untuk keluar")
                
                // Tunggu sinyal
                <-sigChan
                
                // Hapus tunnel
                if err := client.UnregisterTunnel(tunnel.ID); err != nil {
                        log.Error("Gagal menghapus tunnel: %v", err)
                }
                
                // Tutup koneksi
                client.Close()
        },
}

// tcpCmd adalah perintah untuk membuat tunnel TCP
var tcpCmd = &cobra.Command{
        Use:   "tcp [local_port]",
        Short: "Membuat tunnel TCP untuk mengekspos layanan TCP lokal ke internet",
        Long: `Perintah 'tcp' membuat tunnel TCP yang mengekspos layanan TCP lokal ke internet.
Anda dapat menentukan port remote yang akan digunakan untuk mengakses layanan lokal Anda.

Tunnel TCP memungkinkan Anda:
- Mengekspos layanan TCP seperti SSH, database, atau game server ke internet
- Mengakses layanan lokal dari jarak jauh tanpa VPN
- Berbagi layanan TCP dengan orang lain tanpa konfigurasi firewall atau router
- Menguji aplikasi client-server dengan koneksi TCP yang nyata`,
        Example: `  # Membuat tunnel TCP untuk server SSH lokal di port 22
  hxp tcp 22 --remote-port 2222

  # Membuat tunnel TCP untuk database MySQL lokal di port 3306
  hxp tcp 3306 --remote-port 33060

  # Membuat tunnel TCP untuk server game di port 25565
  hxp tcp 25565 --remote-port 25565

  # Menggunakan server dan token autentikasi kustom
  hxp --server haxorport.example.com --token your-auth-token tcp 22 --remote-port 2222`,
        Args:  cobra.ExactArgs(1),
        Run: func(cmd *cobra.Command, args []string) {
                // Parse port lokal
                var localPort int
                if _, err := fmt.Sscanf(args[0], "%d", &localPort); err != nil {
                        log.Error("Port lokal tidak valid: %v", err)
                        os.Exit(1)
                }
                
                // Dapatkan flag
                remotePort, _ := cmd.Flags().GetInt("remote-port")
                
                // Buat konfigurasi tunnel
                tunnelConfig := config.TunnelConfig{
                        Type:       "tcp",
                        LocalPort:  localPort,
                        RemotePort: remotePort,
                }
                
                // Hubungkan ke server
                if err := client.Connect(); err != nil {
                        log.Error("Gagal menghubungkan ke server: %v", err)
                        os.Exit(1)
                }
                
                // Daftarkan tunnel
                tunnel, err := client.RegisterTunnel(tunnelConfig)
                if err != nil {
                        log.Error("Gagal mendaftarkan tunnel: %v", err)
                        os.Exit(1)
                }
                
                // Tangkap sinyal untuk shutdown yang bersih
                sigChan := make(chan os.Signal, 1)
                signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
                
                log.Info("Tekan Ctrl+C untuk keluar")
                
                // Tunggu sinyal
                <-sigChan
                
                // Hapus tunnel
                if err := client.UnregisterTunnel(tunnel.ID); err != nil {
                        log.Error("Gagal menghapus tunnel: %v", err)
                }
                
                // Tutup koneksi
                client.Close()
        },
}

// statusCmd adalah perintah untuk melihat status tunnel
var statusCmd = &cobra.Command{
        Use:   "status",
        Short: "Melihat status semua tunnel yang aktif",
        Long: `Perintah 'status' menampilkan informasi tentang semua tunnel yang aktif saat ini.
Informasi yang ditampilkan meliputi:
- Jenis tunnel (HTTP atau TCP)
- Subdomain atau port remote yang digunakan
- Port lokal yang diteruskan
- ID unik tunnel

Perintah ini berguna untuk:
- Memantau tunnel yang sedang berjalan
- Mendapatkan URL atau port untuk mengakses layanan Anda
- Memverifikasi bahwa tunnel berfungsi dengan benar`,
        Example: `  # Melihat status semua tunnel
  hxp status

  # Melihat status dengan server kustom
  hxp --server haxorport.example.com status`,
        Run: func(cmd *cobra.Command, args []string) {
                // Hubungkan ke server
                if err := client.Connect(); err != nil {
                        log.Error("Gagal menghubungkan ke server: %v", err)
                        os.Exit(1)
                }
                
                // Dapatkan daftar tunnel
                tunnels := client.GetTunnels()
                
                // Tampilkan daftar tunnel
                fmt.Println("Daftar tunnel:")
                if len(tunnels) == 0 {
                        fmt.Println("  Tidak ada tunnel aktif")
                } else {
                        for _, t := range tunnels {
                                if t.Config.Type == "http" {
                                        fmt.Printf("  [HTTP] %s -> %s (ID: %s)\n", t.Config.Subdomain, fmt.Sprintf("localhost:%d", t.Config.LocalPort), t.ID)
                                } else {
                                        fmt.Printf("  [TCP] %d -> localhost:%d (ID: %s)\n", t.Config.RemotePort, t.Config.LocalPort, t.ID)
                                }
                        }
                }
                
                // Tutup koneksi
                client.Close()
        },
}

// configCmd adalah perintah untuk mengelola konfigurasi
var configCmd = &cobra.Command{
        Use:   "config",
        Short: "Mengelola konfigurasi klien haxorport",
        Long: `Perintah 'config' memungkinkan Anda mengelola konfigurasi klien haxorport.
Anda dapat melihat, mengubah, dan menyimpan pengaturan konfigurasi seperti:
- Alamat server haxorport
- Port control plane
- Token autentikasi
- Level logging
- File log

Konfigurasi disimpan dalam file YAML dan dapat digunakan kembali untuk sesi berikutnya.`,
        Example: `  # Melihat konfigurasi saat ini
  hxp config show

  # Mengatur alamat server
  hxp config set server_address haxorport.example.com

  # Mengatur token autentikasi
  hxp config set auth_token your-auth-token

  # Mendapatkan nilai konfigurasi tertentu
  hxp config get server_address`,
}

// configSetCmd adalah perintah untuk mengatur konfigurasi
var configSetCmd = &cobra.Command{
        Use:   "set [key] [value]",
        Short: "Mengatur nilai konfigurasi",
        Long: `Perintah 'config set' memungkinkan Anda mengatur nilai untuk kunci konfigurasi tertentu.
Kunci konfigurasi yang didukung:
- server_address: Alamat server haxorport (contoh: haxorport.example.com)
- control_port: Port control plane server (contoh: 8080)
- auth_token: Token autentikasi untuk mengakses server
- log_level: Level logging (debug, info, warn, error)
- log_file: Path ke file log (kosong untuk stdout)

Perubahan akan disimpan ke file konfigurasi dan akan digunakan untuk sesi berikutnya.`,
        Example: `  # Mengatur alamat server
  hxp config set server_address haxorport.example.com

  # Mengatur port control plane
  hxp config set control_port 8080

  # Mengatur token autentikasi
  hxp config set auth_token your-auth-token

  # Mengatur level logging
  hxp config set log_level debug

  # Mengatur file log
  hxp config set log_file /path/to/haxorport.log`,
        Args:  cobra.ExactArgs(2),
        Run: func(cmd *cobra.Command, args []string) {
                key := args[0]
                value := args[1]
                
                // Atur nilai konfigurasi
                switch key {
                case "server_address":
                        cfg.ServerAddress = value
                case "control_port":
                        var port int
                        if _, err := fmt.Sscanf(value, "%d", &port); err != nil {
                                log.Error("Nilai port tidak valid: %v", err)
                                os.Exit(1)
                        }
                        cfg.ControlPort = port
                case "auth_token":
                        cfg.AuthToken = value
                case "log_level":
                        cfg.LogLevel = value
                case "log_file":
                        cfg.LogFile = value
                default:
                        log.Error("Kunci konfigurasi tidak dikenal: %s", key)
                        os.Exit(1)
                }
                
                // Simpan konfigurasi
                if err := config.SaveConfig(cfg, configPath); err != nil {
                        log.Error("Gagal menyimpan konfigurasi: %v", err)
                        os.Exit(1)
                }
                
                log.Info("Konfigurasi disimpan")
        },
}

// configGetCmd adalah perintah untuk mendapatkan nilai konfigurasi
var configGetCmd = &cobra.Command{
        Use:   "get [key]",
        Short: "Mendapatkan nilai konfigurasi",
        Long: `Perintah 'config get' memungkinkan Anda mendapatkan nilai dari kunci konfigurasi tertentu.
Kunci konfigurasi yang didukung:
- server_address: Alamat server haxorport
- control_port: Port control plane server
- auth_token: Token autentikasi untuk mengakses server
- log_level: Level logging (debug, info, warn, error)
- log_file: Path ke file log (kosong untuk stdout)

Nilai konfigurasi akan ditampilkan di stdout.`,
        Example: `  # Mendapatkan alamat server
  hxp config get server_address

  # Mendapatkan port control plane
  hxp config get control_port

  # Mendapatkan token autentikasi
  hxp config get auth_token

  # Mendapatkan level logging
  hxp config get log_level

  # Mendapatkan file log
  hxp config get log_file`,
        Args:  cobra.ExactArgs(1),
        Run: func(cmd *cobra.Command, args []string) {
                key := args[0]
                
                // Dapatkan nilai konfigurasi
                switch key {
                case "server_address":
                        fmt.Println(cfg.ServerAddress)
                case "control_port":
                        fmt.Println(cfg.ControlPort)
                case "auth_token":
                        fmt.Println(cfg.AuthToken)
                case "log_level":
                        fmt.Println(cfg.LogLevel)
                case "log_file":
                        fmt.Println(cfg.LogFile)
                default:
                        log.Error("Kunci konfigurasi tidak dikenal: %s", key)
                        os.Exit(1)
                }
        },
}

// configShowCmd adalah perintah untuk menampilkan konfigurasi
var configShowCmd = &cobra.Command{
        Use:   "show",
        Short: "Menampilkan seluruh konfigurasi",
        Long: `Perintah 'config show' menampilkan seluruh konfigurasi klien haxorport saat ini.
Informasi yang ditampilkan meliputi:
- Alamat server haxorport
- Port control plane
- Token autentikasi
- Level logging
- File log
- Daftar tunnel yang terkonfigurasi (jika ada)

Perintah ini berguna untuk:
- Memverifikasi konfigurasi saat ini
- Memeriksa tunnel yang telah dikonfigurasi
- Memecahkan masalah koneksi`,
        Example: `  # Menampilkan seluruh konfigurasi
  hxp config show

  # Menampilkan konfigurasi dari file konfigurasi kustom
  hxp --config /path/to/config.yaml config show`,
        Run: func(cmd *cobra.Command, args []string) {
                fmt.Println("Konfigurasi:")
                fmt.Printf("  server_address: %s\n", cfg.ServerAddress)
                fmt.Printf("  control_port: %d\n", cfg.ControlPort)
                fmt.Printf("  auth_token: %s\n", cfg.AuthToken)
                fmt.Printf("  log_level: %s\n", cfg.LogLevel)
                fmt.Printf("  log_file: %s\n", cfg.LogFile)
                
                fmt.Println("Tunnel:")
                if len(cfg.Tunnels) == 0 {
                        fmt.Println("  Tidak ada tunnel terkonfigurasi")
                } else {
                        for i, t := range cfg.Tunnels {
                                fmt.Printf("  %d. %s (%s)\n", i+1, t.Name, t.Type)
                                fmt.Printf("     local_port: %d\n", t.LocalPort)
                                if t.Type == "http" && t.Subdomain != "" {
                                        fmt.Printf("     subdomain: %s\n", t.Subdomain)
                                } else if t.Type == "tcp" && t.RemotePort != 0 {
                                        fmt.Printf("     remote_port: %d\n", t.RemotePort)
                                }
                                if t.Auth != nil {
                                        fmt.Printf("     auth: %s\n", t.Auth.Type)
                                }
                        }
                }
        },
}

func init() {
        // Flag global
        rootCmd.PersistentFlags().StringVar(&configPath, "config", "", "Path ke file konfigurasi (default: ~/.haxorport/config.yaml)")
        rootCmd.PersistentFlags().StringVar(&serverAddr, "server", "", "Alamat server haxorport (contoh: haxorport.example.com)")
        rootCmd.PersistentFlags().IntVar(&serverPort, "port", 0, "Port control plane server (default: 8080)")
        rootCmd.PersistentFlags().StringVar(&authToken, "token", "", "Token autentikasi untuk mengakses server")
        rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "", "Level logging: debug, info, warn, error (default: info)")
        rootCmd.PersistentFlags().StringVar(&logFile, "log-file", "", "Path ke file log, kosong untuk menggunakan stdout")
        
        // Flag untuk httpCmd
        httpCmd.Flags().String("subdomain", "", "Subdomain kustom untuk tunnel HTTP (contoh: myapp)")
        httpCmd.Flags().String("auth-user", "", "Username untuk autentikasi basic HTTP")
        httpCmd.Flags().String("auth-pass", "", "Password untuk autentikasi basic HTTP")
        httpCmd.Flags().String("auth-header-name", "", "Nama header untuk autentikasi HTTP berbasis header")
        httpCmd.Flags().String("auth-header-value", "", "Nilai header untuk autentikasi HTTP berbasis header")
        
        // Flag untuk tcpCmd
        tcpCmd.Flags().Int("remote-port", 0, "Port remote untuk tunnel TCP (wajib diisi)")
        
        // Tambahkan deskripsi penggunaan
        httpCmd.SetUsageTemplate(`Penggunaan:
  {{.UseLine}}

Deskripsi:
  {{.Long}}

Contoh:
{{.Example}}

Flag:
{{.LocalFlags.FlagUsages}}

Flag Global:
{{.InheritedFlags.FlagUsages}}
`)

        tcpCmd.SetUsageTemplate(`Penggunaan:
  {{.UseLine}}

Deskripsi:
  {{.Long}}

Contoh:
{{.Example}}

Flag:
{{.LocalFlags.FlagUsages}}

Flag Global:
{{.InheritedFlags.FlagUsages}}
`)

        statusCmd.SetUsageTemplate(`Penggunaan:
  {{.UseLine}}

Deskripsi:
  {{.Long}}

Contoh:
{{.Example}}

Flag Global:
{{.InheritedFlags.FlagUsages}}
`)

        configCmd.SetUsageTemplate(`Penggunaan:
  {{.UseLine}}

Deskripsi:
  {{.Long}}

Contoh:
{{.Example}}

Subcommand:
  set         Mengatur nilai konfigurasi
  get         Mendapatkan nilai konfigurasi
  show        Menampilkan seluruh konfigurasi

Flag Global:
{{.InheritedFlags.FlagUsages}}
`)
        
        // Tambahkan subcommand ke rootCmd
        rootCmd.AddCommand(httpCmd)
        rootCmd.AddCommand(tcpCmd)
        rootCmd.AddCommand(statusCmd)
        
        // Tambahkan subcommand ke configCmd
        configCmd.AddCommand(configSetCmd)
        configCmd.AddCommand(configGetCmd)
        configCmd.AddCommand(configShowCmd)
        
        // Tambahkan configCmd ke rootCmd
        rootCmd.AddCommand(configCmd)
}

func main() {
        if err := rootCmd.Execute(); err != nil {
                fmt.Println(err)
                os.Exit(1)
        }
}
