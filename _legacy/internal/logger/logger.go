package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

// Level mendefinisikan level logging
type Level int

const (
	// LevelDebug adalah level untuk pesan debug
	LevelDebug Level = iota
	// LevelInfo adalah level untuk pesan informasi
	LevelInfo
	// LevelWarn adalah level untuk pesan peringatan
	LevelWarn
	// LevelError adalah level untuk pesan error
	LevelError
)

// String mengembalikan representasi string dari level
func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// ParseLevel mengkonversi string ke Level
func ParseLevel(level string) Level {
	switch strings.ToLower(level) {
	case "debug":
		return LevelDebug
	case "info":
		return LevelInfo
	case "warn", "warning":
		return LevelWarn
	case "error":
		return LevelError
	default:
		return LevelInfo
	}
}

// Logger adalah wrapper untuk log.Logger dengan dukungan level
type Logger struct {
	logger *log.Logger
	level  Level
	writer io.Writer
}

// NewLogger membuat instance Logger baru
func NewLogger(writer io.Writer, level string) *Logger {
	return &Logger{
		logger: log.New(writer, "", 0),
		level:  ParseLevel(level),
		writer: writer,
	}
}

// SetLevel mengubah level logging
func (l *Logger) SetLevel(level Level) {
	l.level = level
}

// SetLevelFromString mengubah level logging dari string
func (l *Logger) SetLevelFromString(level string) {
	l.level = ParseLevel(level)
}

// log mencatat pesan dengan level tertentu
func (l *Logger) log(level Level, format string, args ...interface{}) {
	if level < l.level {
		return
	}

	// Format waktu
	now := time.Now().Format("2006-01-02 15:04:05.000")
	
	// Format pesan
	var message string
	if len(args) > 0 {
		message = fmt.Sprintf(format, args...)
	} else {
		message = format
	}

	// Log pesan
	l.logger.Printf("[%s] %s %s", now, level.String(), message)
}

// Debug mencatat pesan debug
func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(LevelDebug, format, args...)
}

// Info mencatat pesan informasi
func (l *Logger) Info(format string, args ...interface{}) {
	l.log(LevelInfo, format, args...)
}

// Warn mencatat pesan peringatan
func (l *Logger) Warn(format string, args ...interface{}) {
	l.log(LevelWarn, format, args...)
}

// Error mencatat pesan error
func (l *Logger) Error(format string, args ...interface{}) {
	l.log(LevelError, format, args...)
}

// Close menutup writer jika implementasi io.Closer
func (l *Logger) Close() error {
	if closer, ok := l.writer.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

// DefaultLogger adalah logger default yang menulis ke stdout
var DefaultLogger = NewLogger(os.Stdout, "info")
