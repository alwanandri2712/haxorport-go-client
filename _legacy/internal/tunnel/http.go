package tunnel

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/haxorport/client/internal/config"
)

// HTTPTunnel adalah tunnel untuk protokol HTTP
type HTTPTunnel struct {
	*Tunnel
	server *http.Server
}

// NewHTTPTunnel membuat instance HTTPTunnel baru
func NewHTTPTunnel(id string, config config.TunnelConfig, client *Client) *HTTPTunnel {
	tunnel := NewTunnel(id, config, client)
	return &HTTPTunnel{
		Tunnel: tunnel,
	}
}

// Start memulai tunnel HTTP
func (t *HTTPTunnel) Start() error {
	// Buat URL target
	targetURL := fmt.Sprintf("http://localhost:%d", t.Config.LocalPort)
	target, err := url.Parse(targetURL)
	if err != nil {
		return fmt.Errorf("error parsing URL target: %v", err)
	}

	// Buat reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(target)

	// Modifikasi header
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		t.modifyRequest(req)
	}

	// Tangani error
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		t.client.logger.Error("Proxy error: %v", err)
		http.Error(w, "Error proxy", http.StatusBadGateway)
	}

	// Modifikasi respons
	proxy.ModifyResponse = func(resp *http.Response) error {
		t.modifyResponse(resp)
		return nil
	}

	// Buat server HTTP
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", t.Config.LocalPort),
		Handler: proxy,
	}

	t.server = server

	// Mulai server dalam goroutine
	go func() {
		t.client.logger.Info("Tunnel HTTP %s mendengarkan di %s", t.ID, server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			t.client.logger.Error("Error server HTTP: %v", err)
		}
	}()

	return nil
}

// Stop menghentikan tunnel HTTP
func (t *HTTPTunnel) Stop() {
	// Hentikan server HTTP
	if t.server != nil {
		t.server.Close()
	}

	// Hentikan tunnel dasar
	t.Tunnel.Stop()
}

// modifyRequest memodifikasi permintaan sebelum diteruskan ke target
func (t *HTTPTunnel) modifyRequest(req *http.Request) {
	// Tambahkan header X-Forwarded-*
	req.Header.Set("X-Forwarded-Host", req.Host)
	req.Header.Set("X-Forwarded-Proto", "http")
	
	// Tambahkan header kustom untuk mengidentifikasi tunnel
	req.Header.Set("X-Haxorport-Tunnel-ID", t.ID)
}

// modifyResponse memodifikasi respons sebelum dikirim kembali ke klien
func (t *HTTPTunnel) modifyResponse(resp *http.Response) {
	// Tambahkan header untuk mengidentifikasi haxorport
	resp.Header.Set("X-Powered-By", "Haxorport")
}

// ExtractSubdomain mengekstrak subdomain dari host
func ExtractSubdomain(host, baseDomain string) string {
	// Hapus port jika ada
	if colonIndex := strings.Index(host, ":"); colonIndex != -1 {
		host = host[:colonIndex]
	}

	// Periksa apakah host adalah subdomain dari baseDomain
	if !strings.HasSuffix(host, baseDomain) {
		return ""
	}

	// Ekstrak subdomain
	subdomain := strings.TrimSuffix(host, "."+baseDomain)
	return subdomain
}
