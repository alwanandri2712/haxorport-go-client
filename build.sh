#!/bin/bash

# Script build sederhana untuk Haxorport Client
# Mendukung Linux, macOS, dan Windows (via WSL/Git Bash)
# Author: Haxorport Team

# Warna untuk output
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Fungsi untuk menampilkan pesan dengan warna
print_info() {
    echo -e "${YELLOW}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
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
    
    print_info "Terdeteksi platform: $OS/$ARCH"
}

# Cek apakah Go terinstall
check_go() {
    if ! command -v go &> /dev/null; then
        print_error "Go tidak terinstall. Silakan install Go terlebih dahulu: https://golang.org/doc/install"
    fi
    
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    print_info "Go version: $GO_VERSION"
}

# Bersihkan direktori bin
clean_bin_dir() {
    print_info "Membersihkan direktori bin..."
    mkdir -p bin
}

# Build aplikasi
build_app() {
    print_info "Mengunduh dependensi..."
    go mod download || print_error "Gagal mengunduh dependensi"
    
    print_info "Building aplikasi..."
    if [ "$OS" = "windows" ]; then
        go build -o bin/haxor.exe main.go || print_error "Build gagal"
    else
        go build -o bin/haxor main.go || print_error "Build gagal"
    fi
    
    # Beri izin eksekusi pada binary (untuk Linux dan macOS)
    if [ "$OS" != "windows" ]; then
        chmod +x bin/haxor
    fi
    
    print_success "Build berhasil!"
}

# Jalankan aplikasi setelah build (opsional)
run_app() {
    if [ "$1" = "--run" ]; then
        print_info "Menjalankan aplikasi..."
        if [ "$OS" = "windows" ]; then
            ./bin/haxor.exe "${@:2}"
        else
            ./bin/haxor "${@:2}"
        fi
    fi
}

# Main script
echo -e "${GREEN}=== Haxorport Client Build Script ===${NC}"

# Inisialisasi
detect_platform
check_go
clean_bin_dir

# Build aplikasi
build_app

# Tampilkan informasi
echo -e "\n${GREEN}==================================================${NC}"
echo -e "${GREEN}✅ Haxorport Client berhasil di-build!${NC}"
echo -e "${GREEN}==================================================${NC}"
echo -e "Lokasi binary: $(pwd)/bin/"
echo -e ""
echo -e "Untuk menjalankan:"
if [ "$OS" = "windows" ]; then
    echo -e "  ./bin/haxor.exe --help"
else
    echo -e "  ./bin/haxor --help"
fi
echo -e "${GREEN}==================================================${NC}"

# Jalankan aplikasi jika flag --run diberikan
run_app "$@"
