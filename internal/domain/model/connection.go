package model

// Connection merepresentasikan koneksi tunnel
type Connection struct {
	// ID adalah ID unik koneksi
	ID string
	// TunnelID adalah ID tunnel yang terkait
	TunnelID string
	// Data adalah data yang dikirim melalui koneksi
	Data []byte
}

// NewConnection membuat instance Connection baru
func NewConnection(id string, tunnelID string) *Connection {
	return &Connection{
		ID:       id,
		TunnelID: tunnelID,
	}
}

// SetData mengatur data koneksi
func (c *Connection) SetData(data []byte) {
	c.Data = data
}
