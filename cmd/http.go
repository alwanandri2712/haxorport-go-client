package cmd

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/alwanandri2712/haxorport-go-client/internal/domain/model"
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

// httpCmd is the command to create an HTTP tunnel
var httpCmd = &cobra.Command{
	Use:   "http [target_url]",
	Short: "Create an HTTP tunnel",
	Long: `Create an HTTP tunnel to expose local HTTP services to the internet.
Examples:
  haxor http http://localhost:8080
  haxor http --port 8080 --subdomain myapp
  haxor http --port 3000 --auth basic --username user --password pass`,
	Run: func(cmd *cobra.Command, args []string) {
		// Check if URL argument is provided
		if len(args) > 0 {
			// Parse URL from argument
			targetURL := args[0]

			// Extract port and host from URL
			u, err := url.Parse(targetURL)
			if err != nil {
				fmt.Printf("Error: URL tidak valid: %v\n", err)
				os.Exit(1)
			}

			// Extract port from URL
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
				// Gunakan timestamp untuk membuat subdomain unik tanpa awalan "haxor-"
				timestamp := time.Now().UnixNano() / int64(time.Millisecond)
				httpSubdomain = fmt.Sprintf("%x", timestamp%0xFFFFFF)
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

		// Periksa konfigurasi token terlebih dahulu
		if Container.Config.AuthEnabled {
			if Container.Config.AuthToken == "" {
				fmt.Println("\n===================================================")
				fmt.Println("⚠️ ERROR: Auth token tidak ditemukan dalam konfigurasi")
				fmt.Println("===================================================")
				fmt.Println("Anda perlu menambahkan token autentikasi ke file konfigurasi:")
				fmt.Printf("  %s\n\n", Container.Config.GetConfigFilePath())
				fmt.Println("Cara menambahkan token:")
				fmt.Println("1. Edit file konfigurasi dengan editor teks")
				fmt.Println("2. Temukan baris 'auth_token: \"\"' dan ganti dengan token Anda")
				fmt.Println("3. Simpan file dan jalankan kembali perintah ini")
				fmt.Println("\nAnda bisa mendapatkan token di dashboard Haxorport:")
				fmt.Println("  https://haxorport.online/dashboard")
				fmt.Println("===================================================")
				os.Exit(1)
			}
		}

		// Pastikan client terhubung
		if !Container.Client.IsConnected() {
			if err := Container.Client.Connect(); err != nil {
				fmt.Println("\n===================================================")
				fmt.Printf("⚠️ ERROR: Gagal terhubung ke server\n")
				fmt.Println("===================================================")
				fmt.Printf("Detail error: %v\n", err)
				
				// Berikan saran berdasarkan jenis error
				if Container.Config.AuthEnabled {
					fmt.Println("\nKemungkinan penyebab:")
					fmt.Println("1. Token autentikasi tidak valid")
					fmt.Println("2. Server tidak dapat dijangkau")
					fmt.Println("3. Koneksi internet bermasalah")
					fmt.Println("\nSaran:")
					fmt.Println("- Periksa token autentikasi Anda di file konfigurasi")
					fmt.Printf("  %s\n", Container.Config.GetConfigFilePath())
					fmt.Println("- Pastikan Anda terhubung ke internet")
					fmt.Println("- Coba jalankan './setup.sh' untuk memperbarui konfigurasi")
				}
				fmt.Println("===================================================")
				os.Exit(1)
			}
		}
		
		// Periksa validasi token jika auth diaktifkan
		if Container.Config.AuthEnabled {
			// Check if user data is available (means token has been validated)
			userData := Container.Client.GetUserData()
			if userData == nil {
				fmt.Println("\n===================================================")
				fmt.Println("⚠️ ERROR: Token autentikasi tidak valid")
				fmt.Println("===================================================")
				fmt.Println("Token yang Anda berikan tidak valid atau tidak dapat divalidasi.")
				fmt.Println("\nSaran:")
				fmt.Println("1. Periksa apakah token sudah benar di file konfigurasi:")
				fmt.Printf("   %s\n", Container.Config.GetConfigFilePath())
				fmt.Println("2. Pastikan Anda menggunakan token yang valid dari dashboard Haxorport")
				fmt.Println("   https://haxorport.online/dashboard")
				fmt.Println("===================================================")
				os.Exit(1)
			}
			
			// Periksa batas tunnel
			reached, used, limit := Container.Client.CheckTunnelLimit()
			if reached {
				fmt.Println("\n===================================================")
				fmt.Printf("⚠️ ERROR: Batas tunnel tercapai (%d/%d)\n", used, limit)
				fmt.Println("===================================================")
				fmt.Println("Anda telah mencapai batas tunnel untuk langganan Anda.")
				fmt.Println("\nSaran:")
				fmt.Println("- Tutup beberapa tunnel yang tidak digunakan")
				fmt.Println("- Upgrade langganan Anda untuk mendapatkan lebih banyak tunnel")
				fmt.Println("  https://haxorport.online/pricing")
				fmt.Println("===================================================")
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

		// Tulis ke file log untuk debugging
		logFile, err := os.OpenFile("output.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err == nil {
			defer logFile.Close()
			fmt.Fprintf(logFile, "Tunnel berhasil dibuat: %s\n", tunnel.URL)
		}

		// Gunakan fmt.Fprintf dengan os.Stderr untuk memastikan output ditampilkan
		fmt.Fprintf(os.Stderr, "\n=================================================\n")
		fmt.Fprintf(os.Stderr, "✅ TUNNEL BERHASIL DIBUAT!\n")
		fmt.Fprintf(os.Stderr, "=================================================\n")
		fmt.Fprintf(os.Stderr, "🌐 URL Tunnel: %s\n", tunnel.URL)
		fmt.Fprintf(os.Stderr, "🔌 Port Lokal: %d\n", tunnel.Config.LocalPort)
		fmt.Fprintf(os.Stderr, "🆔 Tunnel ID: %s\n", tunnel.ID)

		// Tampilkan informasi tambahan
		if auth != nil {
			fmt.Fprintf(os.Stderr, "🔒 Autentikasi: %s\n", auth.Type)
		}
		// Informasi server tidak ditampilkan
		fmt.Fprintf(os.Stderr, "📝 Log File: %s\n", Container.Config.LogFile)

		// Tambahkan instruksi untuk mengakses URL
		fmt.Fprintf(os.Stderr, "\n📌 Untuk mengakses layanan Anda, buka URL di atas di browser\n")
		fmt.Fprintf(os.Stderr, "   atau gunakan curl:\n")
		fmt.Fprintf(os.Stderr, "   curl %s\n", tunnel.URL)

		fmt.Fprintf(os.Stderr, "=================================================\n")
		fmt.Fprintf(os.Stderr, "📋 Tekan Ctrl+C untuk menghentikan tunnel\n")
		fmt.Fprintf(os.Stderr, "=================================================\n")

		// Flush stderr untuk memastikan output ditampilkan
		os.Stderr.Sync()

		// Gunakan log.Printf untuk menampilkan output
		log.Printf("Tunnel berhasil dibuat: %s", tunnel.URL)

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

	// Port hanya wajib jika URL tidak diberikan
	// httpCmd.MarkFlagRequired("port")
}
