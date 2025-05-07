package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// demoCmd adalah command untuk menampilkan demo output
var demoCmd = &cobra.Command{
	Use:   "demo",
	Short: "Menampilkan demo output",
	Long:  `Menampilkan demo output untuk melihat tampilan yang diinginkan.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Tampilkan pesan sederhana menggunakan os.Stderr
		fmt.Fprintf(os.Stderr, "Ini adalah pesan demo sederhana\n")
		fmt.Fprintf(os.Stderr, "Jika Anda dapat melihat pesan ini, output berfungsi dengan baik\n")

		// Flush stderr untuk memastikan output ditampilkan
		os.Stderr.Sync()

		// Tunggu beberapa detik
		time.Sleep(5 * time.Second)

		// Tampilkan pesan penutup
		fmt.Fprintf(os.Stderr, "Demo selesai\n")
	},
}

func init() {
	RootCmd.AddCommand(demoCmd)
}
