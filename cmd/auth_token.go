package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// authTokenCmd adalah command untuk mengatur token autentikasi
var authTokenCmd = &cobra.Command{
	Use:   "auth-token [token]",
	Short: "Mengatur token autentikasi",
	Long: `Mengatur token autentikasi untuk koneksi ke server.
Contoh:
  haxor auth-token mFZzPMtTyzZfmF28TWqm_atuaTPwF2WWeExA9CfNS`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Ambil token dari argumen
		token := args[0]

		// Pastikan token tidak kosong
		if token == "" {
			fmt.Println("Error: Token tidak boleh kosong")
			os.Exit(1)
		}

		// Aktifkan autentikasi
		Container.Config.AuthEnabled = true
		Container.Config.AuthToken = token

		// Simpan konfigurasi
		err := Container.ConfigRepository.Save(Container.Config, ConfigPath)
		if err != nil {
			fmt.Printf("Error: Gagal menyimpan konfigurasi: %v\n", err)
			os.Exit(1)
		}

		// Validasi token
		if err := Container.Client.Connect(); err != nil {
			fmt.Printf("Error: Gagal memvalidasi token: %v\n", err)
			os.Exit(1)
		}

		// Ambil data pengguna
		userData := Container.Client.GetUserData()
		if userData == nil {
			fmt.Println("Error: Gagal mendapatkan data pengguna")
			os.Exit(1)
		}

		// Tampilkan informasi pengguna
		fmt.Println("\n=================================================")
		fmt.Println("✅ TOKEN BERHASIL DIATUR DAN DIVALIDASI!")
		fmt.Println("=================================================")
		fmt.Printf("👤 Pengguna: %s (%s)\n", userData.Fullname, userData.Email)
		fmt.Printf("🔑 Langganan: %s\n", userData.Subscription.Name)
		fmt.Printf("📊 Batas Tunnel: %d/%d\n", userData.Subscription.Limits.Tunnels.Used, userData.Subscription.Limits.Tunnels.Limit)
		
		// Tampilkan fitur langganan
		fmt.Println("\n📋 Fitur Langganan:")
		if userData.Subscription.Features.CustomDomains {
			fmt.Println("  ✓ Domain Kustom")
		}
		if userData.Subscription.Features.Analytics {
			fmt.Println("  ✓ Analitik")
		}
		if userData.Subscription.Features.PrioritySupport {
			fmt.Println("  ✓ Dukungan Prioritas")
		}
		
		fmt.Println("\n🔒 Token telah disimpan dalam konfigurasi")
		fmt.Println("=================================================")
	},
}

func init() {
	RootCmd.AddCommand(authTokenCmd)
}
