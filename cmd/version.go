package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version adalah versi aplikasi
const Version = "1.0.0"

// versionCmd adalah command untuk menampilkan versi
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Menampilkan versi",
	Long:  `Menampilkan versi Haxorport Client.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Haxorport Client v%s\n", Version)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
