package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/haxorport/haxor-client/internal/domain/model"
	"github.com/spf13/cobra"
)

var (
	// Config command flags
	configServerAddress string
	configControlPort   int
	configAuthToken     string
	configLogLevel      string
	configLogFile       string
)

// configCmd adalah command untuk mengelola konfigurasi
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Mengelola konfigurasi",
	Long:  `Mengelola konfigurasi Haxorport Client.`,
}

// configShowCmd adalah command untuk menampilkan konfigurasi
var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Menampilkan konfigurasi",
	Long:  `Menampilkan konfigurasi Haxorport Client.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Tampilkan konfigurasi
		fmt.Println("Konfigurasi Haxorport Client:")
		fmt.Printf("Server Address: %s\n", Container.Config.ServerAddress)
		fmt.Printf("Control Port: %d\n", Container.Config.ControlPort)
		fmt.Printf("Auth Token: %s\n", maskString(Container.Config.AuthToken))
		fmt.Printf("Log Level: %s\n", Container.Config.LogLevel)
		fmt.Printf("Log File: %s\n", Container.Config.LogFile)

		// Tampilkan tunnel
		if len(Container.Config.Tunnels) > 0 {
			fmt.Println("\nTunnel:")
			for i, tunnel := range Container.Config.Tunnels {
				fmt.Printf("  %d. %s (%s)\n", i+1, tunnel.Name, tunnel.Type)
				fmt.Printf("     Local Port: %d\n", tunnel.LocalPort)
				if tunnel.Type == model.TunnelTypeHTTP {
					fmt.Printf("     Subdomain: %s\n", tunnel.Subdomain)
				} else if tunnel.Type == model.TunnelTypeTCP {
					fmt.Printf("     Remote Port: %d\n", tunnel.RemotePort)
				}
				if tunnel.Auth != nil {
					fmt.Printf("     Auth: %s\n", tunnel.Auth.Type)
				}
			}
		}
	},
}

// configSetCmd adalah command untuk mengatur konfigurasi
var configSetCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "Mengatur konfigurasi",
	Long: `Mengatur konfigurasi Haxorport Client.
Contoh:
  haxor config set server_address example.com
  haxor config set control_port 8080
  haxor config set auth_token my-token
  haxor config set log_level debug
  haxor config set log_file /path/to/log.txt`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		value := args[1]

		// Update konfigurasi
		switch key {
		case "server_address":
			Container.ConfigService.SetServerAddress(Container.Config, value)
		case "control_port":
			port, err := strconv.Atoi(value)
			if err != nil {
				fmt.Printf("Error: Port harus berupa angka: %v\n", err)
				os.Exit(1)
			}
			Container.ConfigService.SetControlPort(Container.Config, port)
		case "auth_token":
			Container.ConfigService.SetAuthToken(Container.Config, value)
		case "log_level":
			Container.ConfigService.SetLogLevel(Container.Config, value)
		case "log_file":
			Container.ConfigService.SetLogFile(Container.Config, value)
		default:
			fmt.Printf("Error: Kunci konfigurasi tidak valid: %s\n", key)
			os.Exit(1)
		}

		// Simpan konfigurasi
		if err := Container.ConfigService.SaveConfig(Container.Config, ConfigPath); err != nil {
			fmt.Printf("Error: Gagal menyimpan konfigurasi: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Konfigurasi %s berhasil diubah menjadi %s\n", key, value)
	},
}

// configAddTunnelCmd adalah command untuk menambahkan tunnel ke konfigurasi
var configAddTunnelCmd = &cobra.Command{
	Use:   "add-tunnel",
	Short: "Menambahkan tunnel ke konfigurasi",
	Long: `Menambahkan tunnel ke konfigurasi Haxorport Client.
Contoh:
  haxor config add-tunnel --name web --type http --port 8080 --subdomain myapp
  haxor config add-tunnel --name ssh --type tcp --port 22 --remote-port 2222`,
	Run: func(cmd *cobra.Command, args []string) {
		// Validasi parameter
		if httpLocalPort <= 0 {
			fmt.Println("Error: Port lokal harus lebih besar dari 0")
			os.Exit(1)
		}

		// Buat auth jika diperlukan
		var auth *model.TunnelAuth
		if httpAuthType != "" {
			auth = &model.TunnelAuth{}
			switch httpAuthType {
			case "basic":
				auth.Type = model.AuthTypeBasic
				auth.Username = httpUsername
				auth.Password = httpPassword
				if auth.Username == "" || auth.Password == "" {
					fmt.Println("Error: Username dan password diperlukan untuk auth basic")
					os.Exit(1)
				}
			case "header":
				auth.Type = model.AuthTypeHeader
				auth.HeaderName = httpHeader
				auth.HeaderValue = httpValue
				if auth.HeaderName == "" || auth.HeaderValue == "" {
					fmt.Println("Error: Nama dan nilai header diperlukan untuk auth header")
					os.Exit(1)
				}
			default:
				fmt.Printf("Error: Tipe auth tidak valid: %s\n", httpAuthType)
				os.Exit(1)
			}
		}

		// Buat konfigurasi tunnel
		tunnelConfig := model.TunnelConfig{
			Name:      httpSubdomain,
			LocalPort: httpLocalPort,
		}

		// Set tipe tunnel
		tunnelType := cmd.Flag("type").Value.String()
		switch tunnelType {
		case "http":
			tunnelConfig.Type = model.TunnelTypeHTTP
			tunnelConfig.Subdomain = httpSubdomain
		case "tcp":
			tunnelConfig.Type = model.TunnelTypeTCP
			tunnelConfig.RemotePort = tcpRemotePort
		default:
			fmt.Printf("Error: Tipe tunnel tidak valid: %s\n", tunnelType)
			os.Exit(1)
		}

		// Set auth
		tunnelConfig.Auth = auth

		// Tambahkan tunnel ke konfigurasi
		Container.ConfigService.AddTunnel(Container.Config, tunnelConfig)

		// Simpan konfigurasi
		if err := Container.ConfigService.SaveConfig(Container.Config, ConfigPath); err != nil {
			fmt.Printf("Error: Gagal menyimpan konfigurasi: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Tunnel berhasil ditambahkan ke konfigurasi")
	},
}

// configRemoveTunnelCmd adalah command untuk menghapus tunnel dari konfigurasi
var configRemoveTunnelCmd = &cobra.Command{
	Use:   "remove-tunnel [name]",
	Short: "Menghapus tunnel dari konfigurasi",
	Long: `Menghapus tunnel dari konfigurasi Haxorport Client.
Contoh:
  haxor config remove-tunnel web`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		// Hapus tunnel dari konfigurasi
		if !Container.ConfigService.RemoveTunnel(Container.Config, name) {
			fmt.Printf("Error: Tunnel dengan nama %s tidak ditemukan\n", name)
			os.Exit(1)
		}

		// Simpan konfigurasi
		if err := Container.ConfigService.SaveConfig(Container.Config, ConfigPath); err != nil {
			fmt.Printf("Error: Gagal menyimpan konfigurasi: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Tunnel %s berhasil dihapus dari konfigurasi\n", name)
	},
}

// maskString menyembunyikan sebagian string
func maskString(s string) string {
	if len(s) <= 4 {
		return "****"
	}
	return s[:2] + "****" + s[len(s)-2:]
}

func init() {
	RootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configAddTunnelCmd)
	configCmd.AddCommand(configRemoveTunnelCmd)

	// Tambahkan flag untuk add-tunnel
	configAddTunnelCmd.Flags().StringP("name", "n", "", "Nama tunnel")
	configAddTunnelCmd.Flags().StringP("type", "t", "", "Tipe tunnel (http, tcp)")
	configAddTunnelCmd.Flags().IntVarP(&httpLocalPort, "port", "p", 0, "Port lokal yang akan di-tunnel")
	configAddTunnelCmd.Flags().StringVarP(&httpSubdomain, "subdomain", "s", "", "Subdomain yang diminta (untuk HTTP)")
	configAddTunnelCmd.Flags().IntVarP(&tcpRemotePort, "remote-port", "r", 0, "Port remote yang diminta (untuk TCP)")
	configAddTunnelCmd.Flags().StringVarP(&httpAuthType, "auth", "a", "", "Tipe autentikasi (basic, header)")
	configAddTunnelCmd.Flags().StringVarP(&httpUsername, "username", "u", "", "Username untuk autentikasi basic")
	configAddTunnelCmd.Flags().StringVarP(&httpPassword, "password", "w", "", "Password untuk autentikasi basic")
	configAddTunnelCmd.Flags().StringVar(&httpHeader, "header", "", "Nama header untuk autentikasi header")
	configAddTunnelCmd.Flags().StringVar(&httpValue, "value", "", "Nilai header untuk autentikasi header")

	// Tandai flag yang diperlukan
	configAddTunnelCmd.MarkFlagRequired("type")
	configAddTunnelCmd.MarkFlagRequired("port")
}
