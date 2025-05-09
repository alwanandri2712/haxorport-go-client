#!/bin/bash

# Script untuk build dan install Haxorport Client
# Mendukung Linux, macOS, dan Windows (via WSL)

# Fungsi untuk menampilkan pesan error dan keluar
error_exit() {
    echo "ERROR: $1" >&2
    exit 1
}

# Cek apakah Go terinstall
if ! command -v go &> /dev/null; then
    error_exit "Go tidak terinstall. Silakan install Go terlebih dahulu: https://golang.org/doc/install"
fi

# Cek versi Go
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
GO_MAJOR=$(echo $GO_VERSION | cut -d. -f1)
GO_MINOR=$(echo $GO_VERSION | cut -d. -f2)

if [ "$GO_MAJOR" -lt 1 ] || ([ "$GO_MAJOR" -eq 1 ] && [ "$GO_MINOR" -lt 16 ]); then
    error_exit "Go versi 1.16 atau lebih baru diperlukan. Versi terinstall: $GO_VERSION"
fi

# Deteksi OS dan arsitektur
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

case $ARCH in
    x86_64) ARCH="amd64" ;;
    amd64) ARCH="amd64" ;;
    arm64) ARCH="arm64" ;;
    aarch64) ARCH="arm64" ;;
    *) error_exit "Arsitektur tidak didukung: $ARCH" ;;
esac

case $OS in
    linux) OS="linux" ;;
    darwin) OS="darwin" ;;
    *) 
        if grep -q Microsoft /proc/version 2>/dev/null; then
            OS="linux" # WSL terdeteksi
            echo "Windows Subsystem for Linux (WSL) terdeteksi, menggunakan konfigurasi Linux"
        else
            error_exit "OS tidak didukung: $OS"
        fi
        ;;
esac

# Direktori root proyek
ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR" || error_exit "Gagal pindah ke direktori proyek"

# Cek file version.go
VERSION_FILE="$ROOT_DIR/cmd/version.go"
if [ -f "$VERSION_FILE" ]; then
    VERSION=$(grep 'const Version = "' "$VERSION_FILE" | cut -d'"' -f2)
    echo "Building Haxorport Client v$VERSION"
else
    VERSION="dev"
    echo "File version.go tidak ditemukan, menggunakan versi: dev"
fi

# Bersihkan direktori build
echo "Membersihkan direktori build..."
rm -rf "$ROOT_DIR/bin"
mkdir -p "$ROOT_DIR/bin"

# Build untuk platform saat ini
echo "Building untuk $OS/$ARCH..."

if [ "$OS" = "windows" ]; then
    OUTPUT="$ROOT_DIR/bin/haxor.exe"
else
    OUTPUT="$ROOT_DIR/bin/haxor"
fi

echo "Mengunduh dependensi..."
go mod download || error_exit "Gagal mengunduh dependensi"

echo "Building aplikasi..."
GOOS=$OS GOARCH=$ARCH go build -o "$OUTPUT" "$ROOT_DIR/main.go" || error_exit "Build gagal"

echo "✅ Build berhasil!"
echo "Binary location: $OUTPUT"

# Build untuk semua platform jika flag --all diberikan
if [ "$1" = "--all" ]; then
    echo "Building untuk semua platform..."
    
    # Linux (amd64)
    echo "Building untuk linux/amd64..."
    GOOS=linux GOARCH=amd64 go build -o "$ROOT_DIR/bin/haxor-linux-amd64" "$ROOT_DIR/main.go" || echo "⚠️ Build untuk linux/amd64 gagal"
    
    # Linux (arm64)
    echo "Building untuk linux/arm64..."
    GOOS=linux GOARCH=arm64 go build -o "$ROOT_DIR/bin/haxor-linux-arm64" "$ROOT_DIR/main.go" || echo "⚠️ Build untuk linux/arm64 gagal"
    
    # macOS (amd64)
    echo "Building untuk darwin/amd64..."
    GOOS=darwin GOARCH=amd64 go build -o "$ROOT_DIR/bin/haxor-darwin-amd64" "$ROOT_DIR/main.go" || echo "⚠️ Build untuk darwin/amd64 gagal"
    
    # macOS (arm64)
    echo "Building untuk darwin/arm64..."
    GOOS=darwin GOARCH=arm64 go build -o "$ROOT_DIR/bin/haxor-darwin-arm64" "$ROOT_DIR/main.go" || echo "⚠️ Build untuk darwin/arm64 gagal"
    
    # Windows (amd64)
    echo "Building untuk windows/amd64..."
    GOOS=windows GOARCH=amd64 go build -o "$ROOT_DIR/bin/haxor-windows-amd64.exe" "$ROOT_DIR/main.go" || echo "⚠️ Build untuk windows/amd64 gagal"
fi

# Tanya pengguna apakah ingin menginstall binary ke sistem
if [ "$OS" != "windows" ] && [ -z "$INSTALL_SKIP" ]; then
    echo ""
    echo "Apakah Anda ingin menginstall haxorport ke sistem? (y/n)"
    read -r INSTALL_CHOICE
    
    if [ "$INSTALL_CHOICE" = "y" ] || [ "$INSTALL_CHOICE" = "Y" ]; then
        # Tentukan lokasi instalasi berdasarkan OS
        if [ "$OS" = "darwin" ]; then
            INSTALL_DIR="/usr/local/bin"
        else
            INSTALL_DIR="/usr/local/bin"
        fi
        
        echo "Menginstall ke $INSTALL_DIR..."
        if [ -w "$INSTALL_DIR" ]; then
            cp "$OUTPUT" "$INSTALL_DIR/haxor" || error_exit "Gagal menyalin binary ke $INSTALL_DIR"
            chmod +x "$INSTALL_DIR/haxor" || error_exit "Gagal mengatur permission executable"
            echo "✅ Instalasi berhasil! Anda dapat menjalankan 'haxor' dari terminal."
        else
            echo "Memerlukan akses sudo untuk menginstall ke $INSTALL_DIR"
            sudo cp "$OUTPUT" "$INSTALL_DIR/haxor" || error_exit "Gagal menyalin binary ke $INSTALL_DIR"
            sudo chmod +x "$INSTALL_DIR/haxor" || error_exit "Gagal mengatur permission executable"
            echo "✅ Instalasi berhasil! Anda dapat menjalankan 'haxor' dari terminal."
        fi
        
        # Buat konfigurasi default jika belum ada
        CONFIG_DIR="$HOME/.haxorport"
        if [ ! -d "$CONFIG_DIR" ]; then
            mkdir -p "$CONFIG_DIR" || error_exit "Gagal membuat direktori konfigurasi"
            
            if [ -f "$ROOT_DIR/config.example.yaml" ]; then
                cp "$ROOT_DIR/config.example.yaml" "$CONFIG_DIR/config.yaml" || error_exit "Gagal menyalin file konfigurasi contoh"
                echo "✅ File konfigurasi contoh disalin ke $CONFIG_DIR/config.yaml"
            fi
        fi
    fi
fi

echo ""
echo "=================================================="
echo "✅ Haxorport Client berhasil di-build!"
echo "=================================================="
echo "Lokasi binary: $OUTPUT"
echo ""
echo "Untuk menjalankan:"
echo "  $OUTPUT --help"
echo ""
echo "Untuk mengatur token autentikasi:"
echo "  $OUTPUT auth-token YOUR_TOKEN"
echo ""
echo "Untuk membuat tunnel HTTP:"
echo "  $OUTPUT http --port 8080"
echo "=================================================="
