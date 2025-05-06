package port

// Logger adalah interface untuk logging
type Logger interface {
	// Debug mencatat pesan debug
	Debug(format string, args ...interface{})
	
	// Info mencatat pesan informasi
	Info(format string, args ...interface{})
	
	// Warn mencatat pesan peringatan
	Warn(format string, args ...interface{})
	
	// Error mencatat pesan error
	Error(format string, args ...interface{})
	
	// SetLevel mengubah level logging
	SetLevel(level string)
	
	// Close menutup logger
	Close() error
}
