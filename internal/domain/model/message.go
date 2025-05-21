package model

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

// Message adalah struktur dasar untuk semua pesan yang dikirim antara klien dan server
type Message struct {
	// Type adalah tipe pesan
	Type MessageType `json:"type"`
	// Version adalah versi protokol
	Version string `json:"version"`
	// Timestamp adalah waktu pesan dibuat (dalam milidetik sejak epoch)
	Timestamp int64 `json:"timestamp"`
	// Payload adalah data pesan yang sebenarnya
	Payload json.RawMessage `json:"payload,omitempty"`
}

// NewMessage membuat pesan baru dengan tipe dan payload tertentu
func NewMessage(msgType MessageType, payload interface{}) (*Message, error) {
	var payloadJSON json.RawMessage
	var err error

	if payload != nil {
		payloadJSON, err = json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("gagal mengkonversi payload ke JSON: %v", err)
		}
	}

	return &Message{
		Type:      msgType,
		Version:   "1.0.0", // Versi protokol saat ini
		Timestamp: time.Now().UnixNano() / int64(time.Millisecond),
		Payload:   payloadJSON,
	}, nil
}

// ParsePayload mem-parse payload pesan ke dalam struct yang diberikan
func (m *Message) ParsePayload(v interface{}) error {
	if m.Payload == nil {
		return nil
	}
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
	// LocalAddr adalah alamat lokal yang spesifik untuk forwarding (opsional)
	LocalAddr string `json:"local_addr,omitempty"`
	// LocalPort adalah port lokal yang akan di-tunnel
	LocalPort int `json:"local_port"`
	// RemotePort adalah port remote yang diminta (untuk TCP, opsional)
	RemotePort int `json:"remote_port,omitempty"`
	// Auth adalah informasi autentikasi untuk tunnel (opsional)
	Auth *TunnelAuth `json:"auth,omitempty"`
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
