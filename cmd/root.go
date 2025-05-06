package cmd

import (
	"fmt"
	"os"

	"github.com/haxorport/haxor-client/internal/di"
	"github.com/spf13/cobra"
)

var (
	// Container adalah container untuk dependency injection
	Container *di.Container

	// ConfigPath adalah path ke file konfigurasi
	ConfigPath string

	// RootCmd adalah command root untuk CLI
	RootCmd = &cobra.Command{
		Use:   "haxor",
		Short: "Haxorport Client - Tunneling HTTP dan TCP",
		Long: `Haxorport Client adalah alat untuk membuat tunnel HTTP dan TCP.
Dengan Haxorport, Anda dapat mengekspos layanan lokal ke internet.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Inisialisasi container
			Container = di.NewContainer()
			if err := Container.Initialize(ConfigPath); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			// Tutup container
			if Container != nil {
				Container.Close()
			}
		},
	}
)

// Execute menjalankan command root
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Tambahkan flag global
	RootCmd.PersistentFlags().StringVarP(&ConfigPath, "config", "c", "", "Path ke file konfigurasi (default: ~/.haxorport/config.yaml)")
}
