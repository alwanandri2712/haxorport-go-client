package cmd

import (
	"fmt"
	"os"

	"github.com/alwanandri2712/haxorport-go-client/internal/di"
	"github.com/spf13/cobra"
)

var (
	// Container is the dependency injection container
	Container *di.Container

	// ConfigPath is the path to the configuration file
	ConfigPath string

	// RootCmd is the root command for CLI
	RootCmd = &cobra.Command{
		Use:   "haxor",
		Short: "Haxorport Client - HTTP and TCP Tunneling",
		Long: `Haxorport Client is a tool for creating HTTP and TCP tunnels.
With Haxorport, you can expose local services to the internet.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Initialize container
			Container = di.NewContainer()
			if err := Container.Initialize(ConfigPath); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			// Close container
			if Container != nil {
				Container.Close()
			}
		},
	}
)

// Execute runs the root command
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Add global flags
	RootCmd.PersistentFlags().StringVarP(&ConfigPath, "config", "c", "", "Path to configuration file (default: ~/.haxorport/config.yaml)")
}
