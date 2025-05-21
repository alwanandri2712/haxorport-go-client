package transport

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/alwanandri2712/haxorport-go-client/internal/domain/model"
	"github.com/alwanandri2712/haxorport-go-client/internal/domain/port"
)


type TunnelRepository struct {
	client      *Client
	logger      port.Logger
	tunnels     map[string]*model.Tunnel
	connections map[string]net.Conn
	mutex       sync.RWMutex
}


func NewTunnelRepository(client *Client, logger port.Logger) *TunnelRepository {
	repo := &TunnelRepository{
		client:      client,
		logger:      logger,
		tunnels:     make(map[string]*model.Tunnel),
		connections: make(map[string]net.Conn),
		mutex:       sync.RWMutex{},
	}


	client.RegisterHandler(model.MessageTypeData, repo.handleDataMessage)

	return repo
}


func (r *TunnelRepository) Register(config model.TunnelConfig) (*model.Tunnel, error) {
	// Pastikan klien terhubung
	if !r.client.IsConnected() {
		if err := r.client.Connect(); err != nil {
			return nil, fmt.Errorf("failed to connect to server: %v", err)
		}
	}


	response, err := r.client.SendRegisterTunnel(config)
	if err != nil {
		return nil, fmt.Errorf("failed to register tunnel: %v", err)
	}


	if !response.Success {
		return nil, fmt.Errorf("tunnel registration failed: %s", response.Error)
	}


	tunnel := model.NewTunnel(response.TunnelID, config)


	if config.Type == model.TunnelTypeHTTP {
		tunnel.SetHTTPInfo(response.URL)
	} else if config.Type == model.TunnelTypeTCP {
		tunnel.SetTCPInfo(response.RemotePort)
	}


	r.mutex.Lock()
	r.tunnels[response.TunnelID] = tunnel
	r.mutex.Unlock()




	if config.Type == model.TunnelTypeTCP {
		go r.startTunnelListener(tunnel)
	}

	return tunnel, nil
}


func (r *TunnelRepository) Unregister(tunnelID string) error {
	// Pastikan klien terhubung
	if !r.client.IsConnected() {
		if err := r.client.Connect(); err != nil {
			return fmt.Errorf("failed to connect to server: %v", err)
		}
	}


	if err := r.client.SendUnregisterTunnel(tunnelID); err != nil {
		return fmt.Errorf("failed to remove tunnel: %v", err)
	}


	r.mutex.Lock()
	delete(r.tunnels, tunnelID)
	r.mutex.Unlock()

	return nil
}


func (r *TunnelRepository) GetAll() []*model.Tunnel {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	tunnels := make([]*model.Tunnel, 0, len(r.tunnels))
	for _, tunnel := range r.tunnels {
		tunnels = append(tunnels, tunnel)
	}

	return tunnels
}


func (r *TunnelRepository) GetByID(tunnelID string) (*model.Tunnel, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	tunnel, exists := r.tunnels[tunnelID]
	if !exists {
		return nil, fmt.Errorf("tunnel with ID %s not found", tunnelID)
	}

	return tunnel, nil
}


func (r *TunnelRepository) SendData(tunnelID string, connectionID string, data []byte) error {
	return r.client.SendData(tunnelID, connectionID, data)
}


func (r *TunnelRepository) HandleData(tunnelID string, connectionID string, data []byte) error {
	r.mutex.RLock()
	conn, exists := r.connections[connectionID]
	r.mutex.RUnlock()

	if !exists {
		return fmt.Errorf("connection with ID %s not found", connectionID)
	}


	_, err := conn.Write(data)
	if err != nil {
		r.logger.Error("Failed to send data to local connection: %v", err)
		return err
	}

	return nil
}


func (r *TunnelRepository) handleDataMessage(msg *model.Message) error {

	var payload model.DataPayload
	if err := msg.ParsePayload(&payload); err != nil {
		return fmt.Errorf("failed to parse data payload: %v", err)
	}


	return r.HandleData(payload.TunnelID, payload.ConnectionID, payload.Data)
}

// startTunnelListener starts a listener for a tunnel.
func (r *TunnelRepository) startTunnelListener(tunnel *model.Tunnel) {
	localAddr := fmt.Sprintf("%s:%d", tunnel.Config.LocalAddr, tunnel.Config.LocalPort)
	r.logger.Info("Starting listener for tunnel %s on %s", tunnel.ID, localAddr)

	listener, err := net.Listen("tcp", localAddr)
	if err != nil {
		r.logger.Error("Failed to create listener for tunnel %s: %v", tunnel.ID, err)
		return
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			r.logger.Error("Failed to accept connection for tunnel %s: %v", tunnel.ID, err)
			break
		}

		connectionID := fmt.Sprintf("%s-%d", tunnel.ID, time.Now().UnixNano())

		r.mutex.Lock()
		r.connections[connectionID] = conn
		r.mutex.Unlock()

		go r.handleConnection(tunnel.ID, connectionID, conn)
	}
}

// handleConnection handles a connection to a tunnel.
func (r *TunnelRepository) handleConnection(tunnelID string, connectionID string, conn net.Conn) {
	defer func() {
		conn.Close()
		r.mutex.Lock()
		delete(r.connections, connectionID)
		r.mutex.Unlock()
		r.logger.Info("Closing connection %s for tunnel %s", connectionID, tunnelID)
	}()

	buffer := make([]byte, 4096)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err != io.EOF {
				r.logger.Error("Failed to read data from local connection: %v", err)
			}
			break
		}

		if err := r.SendData(tunnelID, connectionID, buffer[:n]); err != nil {
			r.logger.Debug("Sending data to server: %d bytes", n)
			r.logger.Error("Failed to send data to server: %v", err)
			break
		}
	}
}

var _ port.TunnelRepository = (*TunnelRepository)(nil)
