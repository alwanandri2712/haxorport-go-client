package model

import (
	"net/http"
)

// HTTPRequest merepresentasikan permintaan HTTP yang dikirim dari server ke client
type HTTPRequest struct {
	// ID adalah ID unik permintaan
	ID string `json:"id"`
	// TunnelID adalah ID tunnel yang terkait dengan permintaan
	TunnelID string `json:"tunnel_id"`
	// Method adalah metode HTTP (GET, POST, dll)
	Method string `json:"method"`
	// URL adalah URL permintaan
	URL string `json:"url"`
	// Headers adalah header permintaan
	Headers http.Header `json:"headers"`
	// Body adalah body permintaan
	Body []byte `json:"body,omitempty"`
	// LocalPort adalah port lokal yang akan dihubungi oleh client
	LocalPort int `json:"local_port"`
	// RemoteAddr adalah alamat remote dari client HTTP
	RemoteAddr string `json:"remote_addr"`
	// Scheme adalah skema protokol (http atau https)
	Scheme string `json:"scheme,omitempty"`
}

// HTTPResponse merepresentasikan respons HTTP yang dikirim dari client ke server
type HTTPResponse struct {
	// ID adalah ID permintaan yang terkait dengan respons
	ID string `json:"id"`
	// StatusCode adalah kode status HTTP (200, 404, dll)
	StatusCode int `json:"status_code"`
	// Headers adalah header respons
	Headers http.Header `json:"headers"`
	// Body adalah body respons
	Body []byte `json:"body,omitempty"`
	// Error adalah error yang terjadi (jika ada)
	Error string `json:"error,omitempty"`
}
