package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/alwanandri2712/haxorport-go-client/internal/domain/model"
)

// AuthService adalah interface untuk layanan autentikasi
type AuthService interface {
	// ValidateToken memvalidasi token autentikasi
	ValidateToken(token string) (bool, error)
	// ValidateTokenWithResponse memvalidasi token autentikasi dan mengembalikan respons lengkap
	ValidateTokenWithResponse(token string) (*model.AuthResponse, error)
}

// authService adalah implementasi AuthService
type authService struct {
	validationURL string
}

// NewAuthService membuat instance AuthService baru
func NewAuthService(validationURL string) AuthService {
	return &authService{
		validationURL: validationURL,
	}
}

// ValidateToken memvalidasi token autentikasi dengan mengirim request ke API validasi
func (s *authService) ValidateToken(token string) (bool, error) {
	// Gunakan ValidateTokenWithResponse dan hanya kembalikan status valid
	response, err := s.ValidateTokenWithResponse(token)
	if err != nil {
		return false, err
	}
	
	// Periksa apakah respons menunjukkan token valid
	return response.Status == "success" && response.Code == 200, nil
}

// ValidateTokenWithResponse memvalidasi token autentikasi dan mengembalikan respons lengkap
func (s *authService) ValidateTokenWithResponse(token string) (*model.AuthResponse, error) {
	// Jika token kosong, langsung return error
	if token == "" {
		return nil, fmt.Errorf("token tidak boleh kosong")
	}

	// Buat form data
	data := url.Values{}
	data.Set("token", token)

	// Buat request
	req, err := http.NewRequest("POST", s.validationURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("gagal membuat request: %v", err)
	}

	// Set header Content-Type
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// Tambahkan User-Agent untuk menghindari pemblokiran
	req.Header.Set("User-Agent", "HaxorportClient/1.0")

	// Kirim request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("gagal mengirim request: %v", err)
	}
	defer resp.Body.Close()

	// Periksa status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("validasi token gagal dengan status code: %d", resp.StatusCode)
	}

	// Baca seluruh body response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("gagal membaca response body: %v", err)
	}

	// Periksa apakah response adalah JSON yang valid
	if !json.Valid(respBody) {
		// Jika bukan JSON valid, coba lihat 100 karakter pertama untuk debugging
		preview := string(respBody)
		if len(preview) > 100 {
			preview = preview[:100] + "..."
		}
		return nil, fmt.Errorf("respons bukan JSON valid: %s", preview)
	}

	// Parse response
	var authResponse model.AuthResponse
	err = json.Unmarshal(respBody, &authResponse)
	if err != nil {
		return nil, fmt.Errorf("gagal parse response: %v", err)
	}

	return &authResponse, nil
}

// ValidateTokenWithURLEncoded memvalidasi token autentikasi dengan mengirim request dengan format application/x-www-form-urlencoded
func (s *authService) ValidateTokenWithURLEncoded(token string) (bool, error) {
	// Jika token kosong, langsung return false
	if token == "" {
		return false, fmt.Errorf("token tidak boleh kosong")
	}

	// Buat form data
	data := url.Values{}
	data.Set("token", token)

	// Buat request
	req, err := http.NewRequest("POST", s.validationURL, strings.NewReader(data.Encode()))
	if err != nil {
		return false, fmt.Errorf("gagal membuat request: %v", err)
	}

	// Set header Content-Type
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Kirim request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("gagal mengirim request: %v", err)
	}
	defer resp.Body.Close()

	// Periksa status code
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("validasi token gagal dengan status code: %d", resp.StatusCode)
	}

	// Parse response
	var result struct {
		Valid bool `json:"valid"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return false, fmt.Errorf("gagal parse response: %v", err)
	}

	return result.Valid, nil
}
