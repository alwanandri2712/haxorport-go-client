package cmd

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/alwanandri2712/haxorport-go-client/internal/domain/model"
	"github.com/spf13/cobra"
)

var (
	// TCP command flags
	tcpLocalPort  int
	tcpRemotePort int
	tcpLocalAddr  string
)

// tcpCmd is the command to create a TCP tunnel
var tcpCmd = &cobra.Command{
	Use:   "tcp",
	Short: "Create a TCP tunnel",
	Long: `Create a TCP tunnel to expose local TCP services to the internet.
Examples:
  haxor tcp --port 22 --remote-port 2222
  haxor tcp --port 5432`,
	Run: func(cmd *cobra.Command, args []string) {
		// Validate parameters
		if tcpLocalPort <= 0 {
			fmt.Println("Error: Local port must be greater than 0")
			os.Exit(1)
		}

		// Ensure client is connected
		if !Container.Client.IsConnected() {
			if err := Container.Client.Connect(); err != nil {
				fmt.Printf("Error: Failed to connect to server: %v\n", err)
				os.Exit(1)
			}
		}
		
		// Check token validation if auth is enabled
		if Container.Config.AuthEnabled {
			// Check if user data is available (means token has been validated)
			userData := Container.Client.GetUserData()
			if userData == nil {
				fmt.Println("Error: Invalid or unvalidated authentication token")
				os.Exit(1)
			}
			
			// Check tunnel limit
			reached, used, limit := Container.Client.CheckTunnelLimit()
			if reached {
				fmt.Printf("Error: Tunnel limit reached (%d/%d). Please upgrade your subscription.\n", used, limit)
				os.Exit(1)
			}
			
			// User information is already displayed in the log
		}

		// Run client with auto-reconnect
		Container.Client.RunWithReconnect()

		// Parse local address and port
		localHost := "127.0.0.1"
		localPort := tcpLocalPort
		
		// If tcpLocalAddr contains port, parse it
		host, _, err := net.SplitHostPort(tcpLocalAddr)
		if err == nil {
			// If parsing succeeds, use the host part
			if host != "" {
				localHost = host
			}
		} else if !strings.Contains(tcpLocalAddr, ":") {
			// If it's just a hostname/ip without port
			localHost = tcpLocalAddr
		} else {
			fmt.Printf("Error: Invalid local address format: %s\n", tcpLocalAddr)
			os.Exit(1)
		}

		// Create tunnel config
		tunnelConfig := model.TunnelConfig{
			Type:       model.TunnelTypeTCP,
			LocalAddr:  localHost,
			LocalPort:  localPort,
			RemotePort: tcpRemotePort,
		}

		// Create tunnel
		tunnel, err := Container.TunnelService.CreateTCPTunnel(tunnelConfig)
		if err != nil {
			fmt.Printf("Error: Failed to create tunnel: %v\n", err)
			os.Exit(1)
		}

		log.Printf("Membuat tunnel TCP untuk %s:%d dengan port remote %d", tunnelConfig.LocalAddr, tunnelConfig.LocalPort, tunnelConfig.RemotePort)

		fmt.Printf("TCP tunnel created successfully!\n")
		fmt.Printf("Remote Port: %d\n", tunnel.RemotePort)
		fmt.Printf("Local Port: %d\n", tunnel.Config.LocalPort)

		// Display tunnel information without statistics

		// Wait for exit signal
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh

		// Close tunnel
		if err := Container.TunnelService.CloseTunnel(tunnel.ID); err != nil {
			fmt.Printf("Error: Failed to close tunnel: %v\n", err)
		} else {
			fmt.Println("Tunnel closed")
		}
	},
}

func init() {
	RootCmd.AddCommand(tcpCmd)

	// Add flags
	tcpCmd.Flags().IntVarP(&tcpLocalPort, "port", "p", 0, "Local port to tunnel")
	tcpCmd.Flags().IntVarP(&tcpRemotePort, "remote-port", "r", 0, "Requested remote port (optional)")
	tcpCmd.Flags().StringVarP(&tcpLocalAddr, "local-addr", "l", "127.0.0.1", "Local address to forward to (default: 127.0.0.1)")

	// Mark required flags
	tcpCmd.MarkFlagRequired("port")
}
