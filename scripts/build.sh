#!/bin/bash

# Script untuk build Haxorport Client
# Mendukung Linux, macOS, dan Windows (via WSL/Git Bash)

# Fungsi untuk menampilkan pesan error dan keluar
error_exit() {
    echo "ERROR: $1" >&2
    exit 1
}

# Deteksi OS dan arsitektur
detect_platform() {
    # Deteksi OS
    case "$(uname -s)" in
        Linux*)     OS="linux" ;;
        Darwin*)    OS="darwin" ;;
        MINGW*|MSYS*) OS="windows" ;;
        *)          OS="unknown" ;;
    esac
    
    # Deteksi arsitektur
    case "$(uname -m)" in
        x86_64|amd64)  ARCH="amd64" ;;
        arm64|aarch64) ARCH="arm64" ;;
        *)             ARCH="unknown" ;;
    esac
    
    echo "Terdeteksi platform: $OS/$ARCH"
}

# Cek apakah Go terinstall
check_go() {
    if ! command -v go &> /dev/null; then
        error_exit "Go tidak terinstall. Silakan install Go terlebih dahulu: https://golang.org/doc/install"
    fi
    
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    echo "Go version: $GO_VERSION"
}

# Tentukan direktori root proyek
set_root_dir() {
    # Coba gunakan BASH_SOURCE jika tersedia
    if [ -n "${BASH_SOURCE[0]}" ]; then
        ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.."; pwd)"
    else
        # Fallback untuk shell lain
        SCRIPT_DIR="$(cd "$(dirname "$0")"; pwd)"
        ROOT_DIR="$(cd "$SCRIPT_DIR/.."; pwd)"
    fi
    
    cd "$ROOT_DIR" || error_exit "Gagal pindah ke direktori root proyek"
    echo "Direktori proyek: $ROOT_DIR"
}

# Dapatkan versi aplikasi
get_version() {
    VERSION_FILE="$ROOT_DIR/cmd/version.go"
    if [ -f "$VERSION_FILE" ]; then
        VERSION=$(grep 'const Version = "' "$VERSION_FILE" | cut -d'"' -f2)
        echo "Building Haxorport Client v$VERSION"
    else
        VERSION="dev"
        echo "File version.go tidak ditemukan, menggunakan versi: dev"
    fi
}

# Bersihkan direktori build
clean_build_dir() {
    echo "Membersihkan direktori build..."
    rm -rf "$ROOT_DIR/bin"
    mkdir -p "$ROOT_DIR/bin"
}

# Build untuk platform saat ini
build_current_platform() {
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
}

# Build untuk semua platform
build_all_platforms() {
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
    
    echo "✅ Build multi-platform selesai!"
}

# Main script
echo "=== Haxorport Client Build Script ==="

# Inisialisasi
detect_platform
check_go
set_root_dir
get_version
clean_build_dir

# Build aplikasi
build_current_platform

# Build untuk semua platform jika flag --all diberikan
if [ "$1" == "--all" ]; then
    build_all_platforms
fi

echo "\n=================================================="
echo "✅ Haxorport Client berhasil di-build!"
echo "=================================================="
echo "Lokasi binary: $ROOT_DIR/bin/"
echo ""
echo "Untuk menjalankan:"
echo "  $ROOT_DIR/bin/haxor --help"
echo "=================================================="
