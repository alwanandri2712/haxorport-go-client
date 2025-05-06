package model

// TunnelType mendefinisikan tipe tunnel
type TunnelType string

const (
	// TunnelTypeHTTP adalah tipe tunnel untuk protokol HTTP
	TunnelTypeHTTP TunnelType = "http"
	// TunnelTypeTCP adalah tipe tunnel untuk protokol TCP
	TunnelTypeTCP TunnelType = "tcp"
)

// AuthType mendefinisikan tipe autentikasi
type AuthType string

const (
	// AuthTypeBasic adalah tipe autentikasi basic
	AuthTypeBasic AuthType = "basic"
	// AuthTypeHeader adalah tipe autentikasi header
	AuthTypeHeader AuthType = "header"
)

// TunnelAuth adalah konfigurasi autentikasi untuk tunnel
type TunnelAuth struct {
	// Type adalah tipe autentikasi (basic, header)
	Type AuthType
	// Username adalah username untuk autentikasi basic
	Username string
	// Password adalah password untuk autentikasi basic
	Password string
	// HeaderName adalah nama header untuk autentikasi header
	HeaderName string
	// HeaderValue adalah nilai header untuk autentikasi header
	HeaderValue string
}

// TunnelConfig adalah konfigurasi untuk tunnel
type TunnelConfig struct {
	// Name adalah nama tunnel (opsional)
	Name string
	// Type adalah tipe tunnel (http, tcp)
	Type TunnelType
	// LocalPort adalah port lokal yang akan di-tunnel
	LocalPort int
	// Subdomain adalah subdomain yang diminta (untuk HTTP, opsional)
	Subdomain string
	// RemotePort adalah port remote yang diminta (untuk TCP, opsional)
	RemotePort int
	// Auth adalah konfigurasi autentikasi untuk tunnel (opsional)
	Auth *TunnelAuth
}

// Tunnel merepresentasikan tunnel yang dibuat oleh klien
type Tunnel struct {
	// ID adalah ID unik tunnel
	ID string
	// Config adalah konfigurasi tunnel
	Config TunnelConfig
	// URL adalah URL publik untuk tunnel HTTP
	URL string
	// RemotePort adalah port remote untuk tunnel TCP
	RemotePort int
	// Active menunjukkan apakah tunnel aktif
	Active bool
}

// NewTunnel membuat instance Tunnel baru
func NewTunnel(id string, config TunnelConfig) *Tunnel {
	return &Tunnel{
		ID:     id,
		Config: config,
		Active: false,
	}
}

// SetHTTPInfo mengatur informasi untuk tunnel HTTP
func (t *Tunnel) SetHTTPInfo(url string) {
	t.URL = url
	t.Active = true
}

// SetTCPInfo mengatur informasi untuk tunnel TCP
func (t *Tunnel) SetTCPInfo(remotePort int) {
	t.RemotePort = remotePort
	t.Active = true
}

// Deactivate menonaktifkan tunnel
func (t *Tunnel) Deactivate() {
	t.Active = false
}
