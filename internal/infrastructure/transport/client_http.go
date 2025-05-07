package transport

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/haxorport/haxor-client/internal/domain/model"
)

// HandleHTTPRequestMessage menangani pesan permintaan HTTP dari server
func (c *Client) HandleHTTPRequestMessage(msg *model.Message) error {
	// Parse payload
	request, err := msg.ParseHTTPRequestPayload()
	if err != nil {
		c.logger.Error("Gagal mengurai payload permintaan HTTP: %v", err)
		return err
	}

	c.logger.Info("Menerima permintaan HTTP: %s %s", request.Method, request.URL)

	// Buat permintaan HTTP ke layanan lokal
	// Gunakan skema yang diterima dari server (http atau https)
	scheme := "http"
	if request.Scheme != "" {
		scheme = request.Scheme
	}
	targetURL := fmt.Sprintf("%s://localhost:%d%s", scheme, request.LocalPort, request.URL)
	c.logger.Info("Mengirim permintaan ke layanan lokal: %s", targetURL)
	httpReq, err := http.NewRequest(request.Method, targetURL, bytes.NewReader(request.Body))
	if err != nil {
		c.logger.Error("Gagal membuat permintaan HTTP: %v", err)
		return c.sendHTTPErrorResponse(request.ID, err)
	}

	// Salin header
	for key, values := range request.Headers {
		for _, value := range values {
			httpReq.Header.Add(key, value)
		}
	}

	// Tambahkan header X-Forwarded-*
	httpReq.Header.Set("X-Forwarded-Host", request.Headers.Get("Host"))
	httpReq.Header.Set("X-Forwarded-Proto", scheme) // Gunakan skema yang diterima dari server
	httpReq.Header.Set("X-Forwarded-For", request.RemoteAddr)

	// Kirim permintaan
	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		c.logger.Error("Gagal mengirim permintaan HTTP: %v", err)
		return c.sendHTTPErrorResponse(request.ID, err)
	}
	defer resp.Body.Close()

	// Baca body respons
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Error("Gagal membaca body respons: %v", err)
		return c.sendHTTPErrorResponse(request.ID, err)
	}

	// Buat respons HTTP
	httpResp := &model.HTTPResponse{
		ID:         request.ID,
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
		Body:       body,
	}

	// Kirim respons ke server
	return c.sendHTTPResponse(httpResp)
}

// sendHTTPResponse mengirim respons HTTP ke server
func (c *Client) sendHTTPResponse(response *model.HTTPResponse) error {
	// Buat pesan respons HTTP
	msg, err := model.NewHTTPResponseMessage(response)
	if err != nil {
		c.logger.Error("Gagal membuat pesan respons HTTP: %v", err)
		return err
	}

	// Kirim pesan ke server
	return c.sendMessage(msg)
}

// sendHTTPErrorResponse mengirim respons HTTP error ke server
func (c *Client) sendHTTPErrorResponse(requestID string, err error) error {
	// Buat respons HTTP error
	httpResp := &model.HTTPResponse{
		ID:         requestID,
		StatusCode: http.StatusInternalServerError,
		Headers:    http.Header{},
		Error:      err.Error(),
	}

	// Kirim respons ke server
	return c.sendHTTPResponse(httpResp)
}
