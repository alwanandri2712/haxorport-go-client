package port

// Client adalah interface untuk berkomunikasi dengan server haxorport
type Client interface {
	// Connect menghubungkan ke server haxorport
	Connect() error
	
	// Close menutup koneksi ke server
	Close()
	
	// IsConnected mengembalikan status koneksi
	IsConnected() bool
	
	// RunWithReconnect menjalankan klien dengan reconnect otomatis
	RunWithReconnect()
}
