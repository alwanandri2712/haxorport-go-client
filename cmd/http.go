package cmd

import (
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/haxorport/haxor-client/internal/domain/model"
	"github.com/spf13/cobra"
)

var (
	// HTTP command flags
	httpLocalPort int
	httpSubdomain string
	httpAuthType  string
	httpUsername  string
	httpPassword  string
	httpHeader    string
	httpValue     string
)

// httpCmd adalah command untuk membuat HTTP tunnel
var httpCmd = &cobra.Command{
	Use:   "http [target_url]",
	Short: "Membuat HTTP tunnel",
	Long: `Membuat HTTP tunnel untuk mengekspos layanan HTTP lokal ke internet.
Contoh:
  haxor http http://localhost:8080
  haxor http --port 8080 --subdomain myapp
  haxor http --port 3000 --auth basic --username user --password pass`,
	Run: func(cmd *cobra.Command, args []string) {
		// Periksa apakah ada argumen URL
		if len(args) > 0 {
			// Parse URL dari argumen
			targetURL := args[0]

			// Ekstrak port dan host dari URL
			u, err := url.Parse(targetURL)
			if err != nil {
				fmt.Printf("Error: URL tidak valid: %v\n", err)
				os.Exit(1)
			}

			// Ekstrak port dari URL
			port := u.Port()
			if port == "" {
				// Default port berdasarkan skema
				if u.Scheme == "https" {
					port = "443"
				} else {
					port = "80"
				}
			}

			// Konversi port ke integer
			portInt, err := strconv.Atoi(port)
			if err != nil {
				fmt.Printf("Error: Port tidak valid: %v\n", err)
				os.Exit(1)
			}

			// Set port lokal
			httpLocalPort = portInt

			// Generate subdomain otomatis jika tidak ditentukan
			if httpSubdomain == "" {
				// Gunakan timestamp untuk membuat subdomain unik
				timestamp := time.Now().UnixNano() / int64(time.Millisecond)
				httpSubdomain = fmt.Sprintf("haxor-%x", timestamp%0xFFFFFF)
			}
		}

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

		// Pastikan client terhubung
		if !Container.Client.IsConnected() {
			if err := Container.Client.Connect(); err != nil {
				fmt.Printf("Error: Gagal terhubung ke server: %v\n", err)
				os.Exit(1)
			}
		}

		// Jalankan client dengan reconnect otomatis
		Container.Client.RunWithReconnect()

		// Buat tunnel
		tunnel, err := Container.TunnelService.CreateHTTPTunnel(httpLocalPort, httpSubdomain, auth)
		if err != nil {
			fmt.Printf("Error: Gagal membuat tunnel: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Tunnel HTTP berhasil dibuat!\n")
		fmt.Printf("URL: %s\n", tunnel.URL)
		fmt.Printf("Port Lokal: %d\n", tunnel.Config.LocalPort)
		if auth != nil {
			fmt.Printf("Autentikasi: %s\n", auth.Type)
		}

		// Tampilkan informasi tunnel tanpa statistik

		// Tunggu sinyal untuk keluar
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh

		// Tutup tunnel
		if err := Container.TunnelService.CloseTunnel(tunnel.ID); err != nil {
			fmt.Printf("Error: Gagal menutup tunnel: %v\n", err)
		} else {
			fmt.Println("Tunnel ditutup")
		}
	},
}

func init() {
	RootCmd.AddCommand(httpCmd)

	// Tambahkan flag
	httpCmd.Flags().IntVarP(&httpLocalPort, "port", "p", 0, "Port lokal yang akan di-tunnel")
	httpCmd.Flags().StringVarP(&httpSubdomain, "subdomain", "s", "", "Subdomain yang diminta (opsional)")
	httpCmd.Flags().StringVarP(&httpAuthType, "auth", "a", "", "Tipe autentikasi (basic, header)")
	httpCmd.Flags().StringVarP(&httpUsername, "username", "u", "", "Username untuk autentikasi basic")
	httpCmd.Flags().StringVarP(&httpPassword, "password", "w", "", "Password untuk autentikasi basic")
	httpCmd.Flags().StringVar(&httpHeader, "header", "", "Nama header untuk autentikasi header")
	httpCmd.Flags().StringVar(&httpValue, "value", "", "Nilai header untuk autentikasi header")

	// Tandai flag yang diperlukan
	httpCmd.MarkFlagRequired("port")
}
