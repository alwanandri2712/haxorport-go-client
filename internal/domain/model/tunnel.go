package model


type TunnelType string

const (

	TunnelTypeHTTP TunnelType = "http"

	TunnelTypeTCP TunnelType = "tcp"
)


type AuthType string

const (

	AuthTypeBasic AuthType = "basic"

	AuthTypeHeader AuthType = "header"
)


type TunnelAuth struct {

	Type AuthType

	Username string

	Password string

	HeaderName string

	HeaderValue string
}


type TunnelConfig struct {

	Name string

	Type TunnelType

	LocalPort int

	LocalAddr string

	Subdomain string

	RemotePort int

	Auth *TunnelAuth
}


type Tunnel struct {

	ID string

	Config TunnelConfig

	URL string

	RemotePort int

	Active bool
}


func NewTunnel(id string, config TunnelConfig) *Tunnel {
	return &Tunnel{
		ID:     id,
		Config: config,
		Active: false,
	}
}


func (t *Tunnel) SetHTTPInfo(url string) {
	t.URL = url
	t.Active = true
}


func (t *Tunnel) SetTCPInfo(remotePort int) {
	t.RemotePort = remotePort
	t.Active = true
}


func (t *Tunnel) Deactivate() {
	t.Active = false
}
