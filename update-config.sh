#!/bin/bash

# Script untuk memperbarui konfigurasi haxorport-go-client
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
        print_info "Membuat direktori konfigurasi..."
        mkdir -p "$CONFIG_DIR"
    else
        print_info "Membuat backup konfigurasi lama..."
        cp "$CONFIG_FILE" "$CONFIG_FILE.bak.$(date +%Y%m%d%H%M%S)"
        print_success "Backup konfigurasi disimpan"
    fi
}

# Perbarui konfigurasi
update_config() {
    print_info "Memperbarui konfigurasi..."
    
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
    
    # Tampilkan informasi konfigurasi
    print_info "Konfigurasi baru:"
    print_info "- Server: control.haxorport.online:443"
    print_info "- Auth: Enabled"
    print_info "- TLS: Enabled"
    print_info "- Log: $CONFIG_DIR/logs/haxor-client.log (level: warn)"
}

# Uji konfigurasi baru
test_config() {
    print_info "Menguji konfigurasi baru..."
    
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
    
    # Jalankan perintah version
    print_info "Menjalankan: $HAXOR_CMD version"
    $HAXOR_CMD version
    
    if [ $? -eq 0 ]; then
        print_success "Konfigurasi baru berhasil!"
        print_info "Anda sekarang dapat menggunakan perintah berikut untuk membuat tunnel:"
        print_info "$HAXOR_CMD http --port 80"
    else
        print_error "Konfigurasi baru gagal. Silakan jalankan './debug-connection.sh' untuk diagnosa lebih lanjut."
    fi
}

# Main script
echo -e "${GREEN}=== Haxorport Config Update Script ===${NC}"

# Inisialisasi
detect_config_path
update_config
test_config

# Tampilkan informasi
echo -e "\n${GREEN}==================================================${NC}"
echo -e "${GREEN}✅ Konfigurasi berhasil diperbarui!${NC}"
echo -e "${GREEN}==================================================${NC}"
echo -e "Anda sekarang dapat menggunakan perintah berikut:"
echo -e "  haxor http --port 80"
echo -e "  haxor tcp --local-port 22 --remote-port 2222"
echo -e "${GREEN}==================================================${NC}"
