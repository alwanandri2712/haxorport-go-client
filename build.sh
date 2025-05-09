#!/bin/bash

# Script all-in-one untuk Haxorport Client
# Mendukung Linux, macOS, dan Windows (via WSL/Git Bash)
# Author: Haxorport Team

# Warna untuk output
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
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

print_debug() {
    echo -e "${BLUE}[DEBUG]${NC} $1"
}

# Deteksi OS dan arsitektur
detect_platform() {
    # Deteksi OS
    case "$(uname -s)" in
        Linux*)     
            OS="linux" 
            CONFIG_DIR="$HOME/.haxorport"
            if [ "$(id -u)" -eq 0 ]; then
                # Jika root user
                CONFIG_DIR="/etc/haxorport"
            fi
            ;;
        Darwin*)    
            OS="darwin" 
            CONFIG_DIR="$HOME/Library/Preferences/haxorport"
            ;;
        MINGW*|MSYS*) 
            OS="windows"
            CONFIG_DIR="$HOME/.haxorport/config"
            ;;
        *)          
            OS="unknown"
            CONFIG_DIR="$HOME/.haxorport"
            ;;
    esac
    
    # Deteksi arsitektur
    case "$(uname -m)" in
        x86_64|amd64)  ARCH="amd64" ;;
        arm64|aarch64) ARCH="arm64" ;;
        *)             ARCH="unknown" ;;
    esac
    
    CONFIG_FILE="$CONFIG_DIR/config.yaml"
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

# Build untuk semua platform
build_all_platforms() {
    print_info "Building untuk semua platform..."
    
    # Linux (amd64)
    print_info "Building untuk linux/amd64..."
    GOOS=linux GOARCH=amd64 go build -o "bin/haxor-linux-amd64" main.go || print_info "⚠️ Build untuk linux/amd64 gagal"
    
    # Linux (arm64)
    print_info "Building untuk linux/arm64..."
    GOOS=linux GOARCH=arm64 go build -o "bin/haxor-linux-arm64" main.go || print_info "⚠️ Build untuk linux/arm64 gagal"
    
    # macOS (amd64)
    print_info "Building untuk darwin/amd64..."
    GOOS=darwin GOARCH=amd64 go build -o "bin/haxor-darwin-amd64" main.go || print_info "⚠️ Build untuk darwin/amd64 gagal"
    
    # macOS (arm64)
    print_info "Building untuk darwin/arm64..."
    GOOS=darwin GOARCH=arm64 go build -o "bin/haxor-darwin-arm64" main.go || print_info "⚠️ Build untuk darwin/arm64 gagal"
    
    # Windows (amd64)
    print_info "Building untuk windows/amd64..."
    GOOS=windows GOARCH=amd64 go build -o "bin/haxor-windows-amd64.exe" main.go || print_info "⚠️ Build untuk windows/amd64 gagal"
    
    print_success "Build multi-platform selesai!"
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

# Periksa koneksi internet
check_internet() {
    print_info "Memeriksa koneksi internet..."
    if ping -c 1 control.haxorport.online &> /dev/null; then
        print_success "Koneksi internet OK"
        return 0
    else
        print_info "Tidak dapat menjangkau control.haxorport.online, mencoba dengan IP..."
        if ping -c 1 8.8.8.8 &> /dev/null; then
            print_info "Koneksi internet OK, tetapi DNS mungkin bermasalah"
            return 0
        else
            print_error "Tidak ada koneksi internet. Periksa koneksi jaringan Anda."
            return 1
        fi
    fi
}

# Periksa firewall
check_firewall() {
    print_info "Memeriksa akses ke port 443..."
    if nc -z -w 5 control.haxorport.online 443 &> /dev/null; then
        print_success "Port 443 dapat diakses"
        return 0
    else
        print_info "Port 443 tidak dapat diakses. Firewall mungkin memblokir koneksi."
        return 1
    fi
}

# Perbarui konfigurasi
update_config() {
    print_info "Memperbarui konfigurasi..."
    
    # Periksa apakah direktori konfigurasi ada
    if [ ! -d "$CONFIG_DIR" ]; then
        print_info "Membuat direktori konfigurasi: $CONFIG_DIR"
        mkdir -p "$CONFIG_DIR"
    fi
    
    # Backup konfigurasi lama jika ada
    if [ -f "$CONFIG_FILE" ]; then
        print_info "Backup konfigurasi lama disimpan di $CONFIG_FILE.bak"
        cp "$CONFIG_FILE" "$CONFIG_FILE.bak"
    fi
    
    # Minta token dari pengguna
    print_info "Silakan masukkan token autentikasi Anda dari dashboard Haxorport:"
    read -p "Token: " USER_TOKEN
    
    if [ -z "$USER_TOKEN" ]; then
        print_error "Token tidak boleh kosong. Silakan jalankan script ini lagi dan masukkan token yang valid."
    fi
    
    # Buat direktori logs jika belum ada
    mkdir -p "$CONFIG_DIR/logs"
    
    # Perbarui konfigurasi
    cat > "$CONFIG_FILE" << EOF
# Konfigurasi Haxorport Client
server_address: control.haxorport.online
control_port: 443
data_port: 0

# Konfigurasi autentikasi
auth_enabled: true
auth_token: $USER_TOKEN
auth_validation_url: https://haxorport.online/AuthToken/validate

# Konfigurasi TLS
tls_enabled: true
tls_cert: ""
tls_key: ""

# Domain dasar untuk subdomain tunnel
base_domain: "haxorport.online"

# Konfigurasi logging
log_level: warn
log_file: "$CONFIG_DIR/logs/haxor-client.log"
tunnels: []
EOF

    print_success "Konfigurasi berhasil diperbarui!"
}

# Uji koneksi
test_connection() {
    print_info "Menguji koneksi ke server..."
    
    # Cek apakah haxor tersedia
    if command -v haxor &> /dev/null; then
        HAXOR_CMD="haxor"
    elif [ -f "./bin/haxor" ]; then
        HAXOR_CMD="./bin/haxor"
    else
        print_error "Haxor client tidak ditemukan. Pastikan Anda berada di direktori proyek atau haxor sudah terinstal."
    fi
    
    # Jalankan perintah version untuk menguji koneksi
    print_info "Menjalankan perintah version untuk menguji koneksi..."
    VERSION_OUTPUT=$($HAXOR_CMD version 2>&1)
    VERSION_EXIT_CODE=$?
    
    if [ $VERSION_EXIT_CODE -eq 0 ]; then
        print_success "Koneksi berhasil! Version command berjalan dengan baik."
        print_debug "Output: $VERSION_OUTPUT"
    else
        print_info "Perintah version gagal. Mencoba diagnosa lebih lanjut..."
        print_debug "Error: $VERSION_OUTPUT"
        
        # Coba dengan verbose logging
        print_info "Mencoba dengan level log debug..."
        DEBUG_CONFIG="$CONFIG_DIR/debug.yaml"
        
        # Buat konfigurasi debug sementara
        cp "$CONFIG_FILE" "$DEBUG_CONFIG"
        sed -i.bak 's/log_level: warn/log_level: debug/g' "$DEBUG_CONFIG"
        
        # Jalankan dengan konfigurasi debug
        DEBUG_OUTPUT=$($HAXOR_CMD -c "$DEBUG_CONFIG" version 2>&1)
        print_debug "Debug output: $DEBUG_OUTPUT"
        
        # Hapus file konfigurasi debug sementara
        rm -f "$DEBUG_CONFIG" "$DEBUG_CONFIG.bak"
        
        # Periksa masalah umum
        if echo "$DEBUG_OUTPUT" | grep -q "bad handshake"; then
            print_info "Terdeteksi masalah 'bad handshake'. Kemungkinan penyebab:"
            print_info "1. Server mungkin sedang down atau tidak tersedia"
            print_info "2. Token autentikasi mungkin tidak valid"
            print_info "3. Firewall mungkin memblokir koneksi WebSocket"
            print_info "4. Proxy mungkin mengintervensi koneksi WebSocket"
            
            # Coba periksa status server
            print_info "Memeriksa status server..."
            if curl -s -o /dev/null -w "%{http_code}" https://control.haxorport.online 2>/dev/null | grep -q "200"; then
                print_success "Server merespons dengan baik melalui HTTPS"
                print_info "Kemungkinan masalah adalah pada token atau konfigurasi WebSocket"
            else
                print_info "Server tidak merespons dengan baik melalui HTTPS. Server mungkin sedang down."
            fi
        fi
    fi
    
    print_info "Anda dapat mencoba menjalankan perintah HTTP tunnel:"
    print_info "$HAXOR_CMD http --port 80"
}

# Fungsi untuk menampilkan bantuan
show_help() {
    echo -e "${GREEN}Penggunaan:${NC} $0 [OPSI] [COMMAND]"
    echo -e ""
    echo -e "${YELLOW}Opsi:${NC}"
    echo -e "  --help          Menampilkan bantuan ini"
    echo -e "  --all           Build untuk semua platform"
    echo -e "  --run COMMAND   Jalankan aplikasi setelah build"
    echo -e ""
    echo -e "${YELLOW}Perintah:${NC}"
    echo -e "  build           Build aplikasi (default)"
    echo -e "  config          Perbarui konfigurasi"
    echo -e "  test            Uji koneksi ke server"
    echo -e "  fix             Perbaiki masalah koneksi"
    echo -e ""
    echo -e "${YELLOW}Contoh:${NC}"
    echo -e "  $0              Build aplikasi"
    echo -e "  $0 --all        Build untuk semua platform"
    echo -e "  $0 config       Perbarui konfigurasi"
    echo -e "  $0 fix          Perbaiki masalah koneksi"
    echo -e "  $0 --run http   Build dan jalankan HTTP tunnel"
    echo -e ""
}

# Main script
echo -e "${GREEN}=== Haxorport Client Tool ===${NC}"

# Inisialisasi
detect_platform

# Parse argumen
COMMAND="build"
BUILD_ALL=false
RUN_ARGS=""

while [[ $# -gt 0 ]]; do
    case $1 in
        --help)
            show_help
            exit 0
            ;;
        --all)
            BUILD_ALL=true
            shift
            ;;
        --run)
            if [[ $# -gt 1 ]]; then
                RUN_ARGS="${@:2}"
                break
            else
                print_error "--run memerlukan argumen tambahan"
            fi
            ;;
        build|config|test|fix)
            COMMAND=$1
            shift
            ;;
        *)
            print_error "Opsi tidak dikenal: $1. Gunakan --help untuk bantuan."
            ;;
    esac
done

# Eksekusi perintah
case $COMMAND in
    build)
        check_go
        clean_bin_dir
        if [ "$BUILD_ALL" = true ]; then
            build_all_platforms
        else
            build_app
        fi
        
        # Tampilkan informasi
        echo -e "\n${GREEN}==================================================${NC}"
        echo -e "${GREEN}✅ Haxorport Client berhasil di-build!${NC}"
        echo -e "${GREEN}==================================================${NC}"
        echo -e "Lokasi binary: $(pwd)/bin/"
        echo -e ""
        
        # Periksa apakah haxor sudah terinstall secara global
        if command -v haxor &> /dev/null; then
            echo -e "Haxor sudah terinstall secara global. Anda dapat menggunakan:"
            if [ "$OS" = "windows" ]; then
                echo -e "  haxor.exe --help"
            else
                echo -e "  haxor --help"
            fi
            echo -e ""
            echo -e "Atau gunakan binary yang baru di-build:"
        fi
        
        echo -e "Untuk menjalankan binary yang baru di-build:"
        if [ "$OS" = "windows" ]; then
            echo -e "  ./bin/haxor.exe --help"
        else
            echo -e "  ./bin/haxor --help"
        fi
        echo -e ""
        echo -e "Contoh penggunaan:"
        echo -e "  haxor http --port 80"
        echo -e "  haxor tcp --local-port 22 --remote-port 2222"
        echo -e "${GREEN}==================================================${NC}"
        
        # Jalankan aplikasi jika RUN_ARGS tidak kosong
        if [ -n "$RUN_ARGS" ]; then
            print_info "Menjalankan aplikasi dengan argumen: $RUN_ARGS"
            if [ "$OS" = "windows" ]; then
                ./bin/haxor.exe $RUN_ARGS
            else
                ./bin/haxor $RUN_ARGS
            fi
        fi
        ;;
    config)
        update_config
        test_connection
        ;;
    test)
        check_internet
        check_firewall
        test_connection
        ;;
    fix)
        check_internet
        check_firewall
        update_config
        test_connection
        ;;
esac
