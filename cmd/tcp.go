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
	tcpLocalPort  int
	tcpRemotePort int
	tcpLocalAddr  string
)

var tcpCmd = &cobra.Command{
	Use:   "tcp",
	Short: "Create a TCP tunnel",
	Long: `Create a TCP tunnel to expose local TCP services to the internet.
Examples:
  haxor tcp --port 22 --remote-port 2222
  haxor tcp --port 5432`,
	Run: func(cmd *cobra.Command, args []string) {
		if tcpLocalPort <= 0 {
			fmt.Println("Error: Local port must be greater than 0")
			os.Exit(1)
		}

		if !Container.Client.IsConnected() {
			if err := Container.Client.Connect(); err != nil {
				fmt.Printf("Error: Failed to connect to server: %v\n", err)
				os.Exit(1)
			}
		}

		if Container.Config.AuthEnabled {
			userData := Container.Client.GetUserData()
			if userData == nil {
				fmt.Println("Error: Invalid or unvalidated authentication token")
				os.Exit(1)
			}

			reached, used, limit := Container.Client.CheckTunnelLimit()
			if reached {
				fmt.Printf("Error: Tunnel limit reached (%d/%d). Please upgrade your subscription.\n", used, limit)
				os.Exit(1)
			}
		}

		Container.Client.RunWithReconnect()

		localHost := "127.0.0.1"
		localPort := tcpLocalPort

		host, _, err := net.SplitHostPort(tcpLocalAddr)
		if err == nil {
			if host != "" {
				localHost = host
			}
		} else if !strings.Contains(tcpLocalAddr, ":") {
			localHost = tcpLocalAddr
		} else {
			fmt.Printf("Error: Invalid local address format: %s\n", tcpLocalAddr)
			os.Exit(1)
		}

		tunnelConfig := model.TunnelConfig{
			Type:       model.TunnelTypeTCP,
			LocalAddr:  localHost,
			LocalPort:  localPort,
			RemotePort: tcpRemotePort,
		}

		tunnel, err := Container.TunnelService.CreateTCPTunnel(tunnelConfig)
		if err != nil {
			fmt.Printf("Error: Failed to create tunnel: %v\n", err)
			os.Exit(1)
		}

		log.Printf("Creating TCP tunnel for %s:%d with remote port %d", tunnelConfig.LocalAddr, tunnelConfig.LocalPort, tunnelConfig.RemotePort)

		fmt.Printf("TCP tunnel created successfully!\n")
		fmt.Printf("Remote Port: %d\n", tunnel.RemotePort)
		fmt.Printf("Local Port: %d\n", tunnel.Config.LocalPort)

		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh

		if err := Container.TunnelService.CloseTunnel(tunnel.ID); err != nil {
			fmt.Printf("Error: Failed to close tunnel: %v\n", err)
		} else {
			fmt.Println("Tunnel closed")
		}
	},
}

func init() {
	RootCmd.AddCommand(tcpCmd)

	tcpCmd.Flags().IntVarP(&tcpLocalPort, "port", "p", 0, "Local port to tunnel")
	tcpCmd.Flags().IntVarP(&tcpRemotePort, "remote-port", "r", 0, "Requested remote port (optional)")
	tcpCmd.Flags().StringVarP(&tcpLocalAddr, "local-addr", "l", "127.0.0.1", "Local address to forward to (default: 127.0.0.1)")

	tcpCmd.MarkFlagRequired("port")
}
