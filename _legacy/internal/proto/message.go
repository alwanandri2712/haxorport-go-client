package proto

import (
	"encoding/json"
	"fmt"
	"time"
)

// MessageType mendefinisikan tipe pesan yang dikirim antara klien dan server
type MessageType string

const (
	// MessageTypeAuth adalah pesan untuk autentikasi
	MessageTypeAuth MessageType = "auth"
	// MessageTypeRegister adalah pesan untuk mendaftarkan tunnel
	MessageTypeRegister MessageType = "register"
	// MessageTypeUnregister adalah pesan untuk menghapus tunnel
	MessageTypeUnregister MessageType = "unregister"
	// MessageTypeData adalah pesan yang berisi data tunnel
	MessageTypeData MessageType = "data"
	// MessageTypePing adalah pesan ping untuk menjaga koneksi tetap hidup
	MessageTypePing MessageType = "ping"
	// MessageTypePong adalah respons terhadap pesan ping
	MessageTypePong MessageType = "pong"
	// MessageTypeError adalah pesan error
	MessageTypeError MessageType = "error"
)

// ProtocolVersion adalah versi protokol saat ini
const ProtocolVersion = "1.0.0"

// Message adalah struktur dasar untuk semua pesan yang dikirim antara klien dan server
type Message struct {
	// Type adalah tipe pesan
	Type MessageType `json:"type"`
	// Version adalah versi protokol
	Version string `json:"version"`
	// Timestamp adalah waktu pesan dibuat
	Timestamp int64 `json:"timestamp"`
	// Payload adalah data pesan yang sebenarnya
	Payload json.RawMessage `json:"payload"`
}

// NewMessage membuat pesan baru dengan tipe dan payload tertentu
func NewMessage(msgType MessageType, payload interface{}) (*Message, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshaling payload: %v", err)
	}

	return &Message{
		Type:      msgType,
		Version:   ProtocolVersion,
		Timestamp: time.Now().UnixNano() / int64(time.Millisecond),
		Payload:   payloadBytes,
	}, nil
}

// ParsePayload mem-parse payload pesan ke dalam struct yang diberikan
func (m *Message) ParsePayload(v interface{}) error {
	return json.Unmarshal(m.Payload, v)
}

// AuthPayload adalah payload untuk pesan autentikasi
type AuthPayload struct {
	// Token adalah token autentikasi
	Token string `json:"token"`
}

// RegisterPayload adalah payload untuk pesan pendaftaran tunnel
type RegisterPayload struct {
	// TunnelType adalah tipe tunnel (http, tcp)
	TunnelType string `json:"tunnel_type"`
	// Subdomain adalah subdomain yang diminta (opsional)
	Subdomain string `json:"subdomain,omitempty"`
	// LocalPort adalah port lokal yang akan di-tunnel
	LocalPort int `json:"local_port"`
	// RemotePort adalah port remote yang diminta (untuk TCP, opsional)
	RemotePort int `json:"remote_port,omitempty"`
	// Auth adalah informasi autentikasi untuk tunnel (opsional)
	Auth *TunnelAuth `json:"auth,omitempty"`
}

// TunnelAuth adalah informasi autentikasi untuk tunnel
type TunnelAuth struct {
	// Type adalah tipe autentikasi (basic, header)
	Type string `json:"type"`
	// Username adalah username untuk autentikasi basic
	Username string `json:"username,omitempty"`
	// Password adalah password untuk autentikasi basic
	Password string `json:"password,omitempty"`
	// HeaderName adalah nama header untuk autentikasi header
	HeaderName string `json:"header_name,omitempty"`
	// HeaderValue adalah nilai header untuk autentikasi header
	HeaderValue string `json:"header_value,omitempty"`
}

// UnregisterPayload adalah payload untuk pesan penghapusan tunnel
type UnregisterPayload struct {
	// TunnelID adalah ID tunnel yang akan dihapus
	TunnelID string `json:"tunnel_id"`
}

// DataPayload adalah payload untuk pesan data
type DataPayload struct {
	// TunnelID adalah ID tunnel yang terkait dengan data
	TunnelID string `json:"tunnel_id"`
	// ConnectionID adalah ID koneksi yang terkait dengan data
	ConnectionID string `json:"connection_id"`
	// Data adalah data yang dikirim
	Data []byte `json:"data"`
}

// ErrorPayload adalah payload untuk pesan error
type ErrorPayload struct {
	// Code adalah kode error
	Code string `json:"code"`
	// Message adalah pesan error
	Message string `json:"message"`
}

// RegisterResponsePayload adalah payload untuk respons terhadap pesan pendaftaran
type RegisterResponsePayload struct {
	// Success menunjukkan apakah pendaftaran berhasil
	Success bool `json:"success"`
	// TunnelID adalah ID tunnel yang dibuat
	TunnelID string `json:"tunnel_id"`
	// URL adalah URL publik untuk tunnel HTTP
	URL string `json:"url,omitempty"`
	// RemotePort adalah port remote untuk tunnel TCP
	RemotePort int `json:"remote_port,omitempty"`
	// Error adalah pesan error jika pendaftaran gagal
	Error string `json:"error,omitempty"`
}
