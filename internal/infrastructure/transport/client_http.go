package transport

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/alwanandri2712/haxorport-go-client/internal/domain/model"
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

	// Buat permintaan HTTP ke layanan lokal di komputer klien
	// Selalu gunakan HTTP untuk koneksi lokal, terlepas dari skema yang diterima dari server
	// Ini karena layanan lokal biasanya hanya mendukung HTTP
	scheme := "http"
	
	// Gunakan localhost di komputer klien, bukan di server
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

	// Kirim permintaan ke layanan lokal melalui koneksi balik
	c.logger.Info("Membuat koneksi HTTP ke layanan lokal dengan metode %s", request.Method)
	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		c.logger.Error("Gagal mengirim permintaan HTTP ke layanan lokal: %v", err)
		return c.sendHTTPErrorResponse(request.ID, err)
	}
	c.logger.Info("Berhasil terhubung ke layanan lokal, status: %d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	defer resp.Body.Close()

	// Baca body respons
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Error("Gagal membaca body respons: %v", err)
		return c.sendHTTPErrorResponse(request.ID, err)
	}

	// Periksa Content-Type untuk menentukan apakah ini adalah HTML
	contentType := resp.Header.Get("Content-Type")
	if strings.Contains(contentType, "text/html") {
		// Ganti URL lokal dengan URL tunnel dalam respons HTML
		localURLPrefix := fmt.Sprintf("http://localhost:%d", request.LocalPort)
		localURLPrefixSecure := fmt.Sprintf("https://localhost:%d", request.LocalPort)
		
		// Buat URL tunnel berdasarkan skema yang diterima
		tunnelScheme := "http"
		if request.Scheme == "https" {
			tunnelScheme = "https"
		}
		
		// Ekstrak hostname yang tepat dari header Host
		hostname := ""
		if host, ok := request.Headers["Host"]; ok && len(host) > 0 {
			hostname = host[0]
		}
		
		// Jika hostname masih kosong, gunakan X-Forwarded-Host
		if hostname == "" {
			if host, ok := request.Headers["X-Forwarded-Host"]; ok && len(host) > 0 {
				hostname = host[0]
			}
		}
		
		// Jika hostname masih kosong, ekstrak subdomain dari URL tunnel
		if hostname == "" {
			// Coba dapatkan subdomain dari URL yang diberikan oleh pengguna
			subdomain := c.GetSubdomain()
			if subdomain != "" {
				hostname = subdomain + ".haxorport.online"
			} else {
				// Fallback ke tunnel ID jika subdomain tidak tersedia
				hostname = request.TunnelID + ".haxorport.online"
			}
		}
		
		c.logger.Info(fmt.Sprintf("Menggunakan hostname: %s untuk penggantian URL", hostname))
		tunnelURLPrefix := fmt.Sprintf("%s://%s", tunnelScheme, hostname)
		
		// Ganti URL dalam body
		bodyStr := string(body)
		bodyStr = strings.ReplaceAll(bodyStr, localURLPrefix, tunnelURLPrefix)
		bodyStr = strings.ReplaceAll(bodyStr, localURLPrefixSecure, tunnelURLPrefix)
		
		// Ganti URL relatif dalam href dan src
		// Contoh: href="/path" menjadi href="https://subdomain.haxorport.online/path"
		bodyStr = strings.ReplaceAll(bodyStr, "href=\"/", "href=\""+tunnelURLPrefix+"/")
		bodyStr = strings.ReplaceAll(bodyStr, "src=\"/", "src=\""+tunnelURLPrefix+"/")
		
		// Update body dengan konten yang telah dimodifikasi
		body = []byte(bodyStr)
		
		c.logger.Info("URL lokal dalam respons HTML diganti dengan URL tunnel")
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
