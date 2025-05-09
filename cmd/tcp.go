package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

var (
	// TCP command flags
	tcpLocalPort  int
	tcpRemotePort int
)

// tcpCmd adalah command untuk membuat TCP tunnel
var tcpCmd = &cobra.Command{
	Use:   "tcp",
	Short: "Membuat TCP tunnel",
	Long: `Membuat TCP tunnel untuk mengekspos layanan TCP lokal ke internet.
Contoh:
  haxor tcp --port 22 --remote-port 2222
  haxor tcp --port 5432`,
	Run: func(cmd *cobra.Command, args []string) {
		// Validasi parameter
		if tcpLocalPort <= 0 {
			fmt.Println("Error: Port lokal harus lebih besar dari 0")
			os.Exit(1)
		}

		// Pastikan client terhubung
		if !Container.Client.IsConnected() {
			if err := Container.Client.Connect(); err != nil {
				fmt.Printf("Error: Gagal terhubung ke server: %v\n", err)
				os.Exit(1)
			}
		}
		
		// Periksa validasi token jika auth diaktifkan
		if Container.Config.AuthEnabled {
			// Periksa apakah data pengguna tersedia (berarti token sudah divalidasi)
			userData := Container.Client.GetUserData()
			if userData == nil {
				fmt.Println("Error: Token autentikasi tidak valid atau belum divalidasi")
				os.Exit(1)
			}
			
			// Periksa batas tunnel
			reached, used, limit := Container.Client.CheckTunnelLimit()
			if reached {
				fmt.Printf("Error: Batas tunnel tercapai (%d/%d). Upgrade langganan Anda.\n", used, limit)
				os.Exit(1)
			}
			
			// Informasi pengguna sudah ditampilkan di log
		}

		// Jalankan client dengan reconnect otomatis
		Container.Client.RunWithReconnect()

		// Buat tunnel
		tunnel, err := Container.TunnelService.CreateTCPTunnel(tcpLocalPort, tcpRemotePort)
		if err != nil {
			fmt.Printf("Error: Gagal membuat tunnel: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Tunnel TCP berhasil dibuat!\n")
		fmt.Printf("Port Remote: %d\n", tunnel.RemotePort)
		fmt.Printf("Port Lokal: %d\n", tunnel.Config.LocalPort)

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
	RootCmd.AddCommand(tcpCmd)

	// Tambahkan flag
	tcpCmd.Flags().IntVarP(&tcpLocalPort, "port", "p", 0, "Port lokal yang akan di-tunnel")
	tcpCmd.Flags().IntVarP(&tcpRemotePort, "remote-port", "r", 0, "Port remote yang diminta (opsional)")

	// Tandai flag yang diperlukan
	tcpCmd.MarkFlagRequired("port")
}
