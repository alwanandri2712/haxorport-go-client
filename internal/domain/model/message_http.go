package model

import (
	"encoding/json"
	"net/http"
)

// MessageTypeHTTPRequest adalah tipe pesan untuk permintaan HTTP
const MessageTypeHTTPRequest MessageType = "http_request"

// MessageTypeHTTPResponse adalah tipe pesan untuk respons HTTP
const MessageTypeHTTPResponse MessageType = "http_response"

// HTTPRequestPayload adalah payload untuk pesan permintaan HTTP
type HTTPRequestPayload struct {
	// Request adalah permintaan HTTP
	Request *HTTPRequest `json:"request"`
}

// HTTPResponsePayload adalah payload untuk pesan respons HTTP
type HTTPResponsePayload struct {
	// Response adalah respons HTTP
	Response *HTTPResponse `json:"response"`
}

// NewHTTPRequestMessage membuat pesan permintaan HTTP baru
func NewHTTPRequestMessage(request *HTTPRequest) (*Message, error) {
	payload := HTTPRequestPayload{
		Request: request,
	}
	return NewMessage(MessageTypeHTTPRequest, payload)
}

// NewHTTPResponseMessage membuat pesan respons HTTP baru
func NewHTTPResponseMessage(response *HTTPResponse) (*Message, error) {
	payload := HTTPResponsePayload{
		Response: response,
	}
	return NewMessage(MessageTypeHTTPResponse, payload)
}

// ParseHTTPRequestPayload mengurai payload permintaan HTTP
func (m *Message) ParseHTTPRequestPayload() (*HTTPRequest, error) {
	var payload HTTPRequestPayload
	if err := m.ParsePayload(&payload); err != nil {
		return nil, err
	}
	return payload.Request, nil
}

// ParseHTTPResponsePayload mengurai payload respons HTTP
func (m *Message) ParseHTTPResponsePayload() (*HTTPResponse, error) {
	var payload HTTPResponsePayload
	if err := m.ParsePayload(&payload); err != nil {
		return nil, err
	}
	return payload.Response, nil
}
