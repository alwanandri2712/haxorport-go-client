#!/bin/bash

# Script untuk debug koneksi di haxorport-go-client
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

# Deteksi OS dan path konfigurasi
detect_config_path() {
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
    
    CONFIG_FILE="$CONFIG_DIR/config.yaml"
    print_info "Terdeteksi platform: $OS"
    print_info "File konfigurasi: $CONFIG_FILE"
    
    # Periksa apakah file konfigurasi ada
    if [ ! -f "$CONFIG_FILE" ]; then
        print_info "File konfigurasi tidak ditemukan di $CONFIG_FILE"
        print_info "Mencoba menggunakan konfigurasi lokal..."
        CONFIG_FILE="./config.yaml"
        if [ ! -f "$CONFIG_FILE" ]; then
            print_error "File konfigurasi tidak ditemukan di $CONFIG_FILE"
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

# Buat konfigurasi debug
create_debug_config() {
    print_info "Membuat konfigurasi debug..."
    
    # Minta token dari pengguna
    print_info "Silakan masukkan token autentikasi Anda dari dashboard Haxorport:"
    read -p "Token: " USER_TOKEN
    
    if [ -z "$USER_TOKEN" ]; then
        print_error "Token tidak boleh kosong. Silakan jalankan script ini lagi dan masukkan token yang valid."
    fi
    
    # Buat direktori logs jika belum ada
    mkdir -p "debug"
    
    # Perbarui konfigurasi
    DEBUG_CONFIG="./debug/debug.yaml"
    cat > "$DEBUG_CONFIG" << EOF
# Konfigurasi Haxorport Client (Debug)
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
log_level: debug
log_file: "debug/haxor-debug.log"
tunnels: []
EOF

    print_success "Konfigurasi debug berhasil dibuat di $DEBUG_CONFIG"
    return 0
}

# Jalankan dengan konfigurasi debug
run_debug() {
    print_info "Menjalankan dengan konfigurasi debug..."
    
    # Cek apakah haxor tersedia
    if command -v haxor &> /dev/null; then
        HAXOR_CMD="haxor"
    elif [ -f "./bin/haxor" ]; then
        HAXOR_CMD="./bin/haxor"
    else
        print_info "Binary haxor tidak ditemukan, mencoba dengan go run..."
        if command -v go &> /dev/null; then
            HAXOR_CMD="go run main.go"
        else
            print_error "Tidak dapat menemukan binary haxor atau Go. Pastikan Anda berada di direktori proyek atau haxor sudah terinstal."
        fi
    fi
    
    # Jalankan dengan konfigurasi debug
    DEBUG_CONFIG="./debug/debug.log"
    print_info "Menjalankan: $HAXOR_CMD -c ./debug/debug.yaml version"
    $HAXOR_CMD -c ./debug/debug.yaml version
    
    # Jika berhasil, coba HTTP tunnel
    if [ $? -eq 0 ]; then
        print_success "Koneksi berhasil! Mencoba membuat HTTP tunnel..."
        print_info "Menjalankan: $HAXOR_CMD -c ./debug/debug.yaml http --port 80"
        $HAXOR_CMD -c ./debug/debug.yaml http --port 80
    else
        print_info "Koneksi gagal. Memeriksa log debug..."
        if [ -f "debug/haxor-debug.log" ]; then
            print_debug "Isi log debug:"
            cat debug/haxor-debug.log
        fi
    fi
}

# Periksa perbedaan dengan haxorport
check_haxorport() {
    print_info "Memeriksa apakah haxorport terinstall..."
    if command -v haxorport &> /dev/null; then
        print_success "haxorport ditemukan!"
        print_info "Menjalankan: haxorport version"
        haxorport version
        
        print_info "Mencoba membuat HTTP tunnel dengan haxorport..."
        print_info "Menjalankan: haxorport http --port 80"
        haxorport http --port 80 &
        HAXORPORT_PID=$!
        
        # Tunggu sebentar
        sleep 5
        
        # Hentikan proses
        kill $HAXORPORT_PID 2>/dev/null
        
        print_info "Jika haxorport berhasil tapi haxor gagal, kemungkinan ada perbedaan konfigurasi."
        print_info "Coba periksa konfigurasi haxorport dengan: haxorport config show"
    else
        print_info "haxorport tidak ditemukan."
    fi
}

# Periksa versi websocket
check_websocket_version() {
    print_info "Memeriksa versi library websocket..."
    if command -v go &> /dev/null; then
        print_info "Versi Go:"
        go version
        
        print_info "Versi library websocket:"
        go list -m github.com/gorilla/websocket
    else
        print_info "Go tidak terinstall, tidak dapat memeriksa versi websocket."
    fi
}

# Main script
echo -e "${GREEN}=== Haxorport Connection Debug Script ===${NC}"

# Inisialisasi
detect_config_path
check_internet
check_firewall
create_debug_config
run_debug
check_haxorport
check_websocket_version

# Tampilkan informasi
echo -e "\n${GREEN}==================================================${NC}"
echo -e "${GREEN}Debug koneksi selesai!${NC}"
echo -e "${GREEN}==================================================${NC}"
echo -e "Jika masih mengalami masalah, coba periksa:"
echo -e "1. Koneksi internet dan firewall Anda"
echo -e "2. Validitas token autentikasi"
echo -e "3. Status server Haxorport"
echo -e "4. Log debug di debug/haxor-debug.log"
echo -e "${GREEN}==================================================${NC}"
