package port

import "github.com/alwanandri2712/haxorport-go-client/internal/domain/model"

// ConfigRepository adalah interface untuk mengakses dan memanipulasi konfigurasi
type ConfigRepository interface {
	// Load memuat konfigurasi dari penyimpanan
	Load(path string) (*model.Config, error)
	
	// Save menyimpan konfigurasi ke penyimpanan
	Save(config *model.Config, path string) error
	
	// GetDefaultPath mengembalikan path default untuk file konfigurasi
	GetDefaultPath() (string, error)
}
