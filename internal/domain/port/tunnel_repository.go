package port

import "github.com/haxorport/haxor-client/internal/domain/model"

// TunnelRepository adalah interface untuk mengakses dan memanipulasi tunnel
type TunnelRepository interface {
	// Register mendaftarkan tunnel baru ke server
	Register(config model.TunnelConfig) (*model.Tunnel, error)
	
	// Unregister menghapus tunnel dari server
	Unregister(tunnelID string) error
	
	// GetAll mengembalikan semua tunnel yang aktif
	GetAll() []*model.Tunnel
	
	// GetByID mengembalikan tunnel berdasarkan ID
	GetByID(tunnelID string) (*model.Tunnel, error)
	
	// SendData mengirim data melalui tunnel
	SendData(tunnelID string, connectionID string, data []byte) error
	
	// HandleData menangani data yang diterima dari server
	HandleData(tunnelID string, connectionID string, data []byte) error
}
