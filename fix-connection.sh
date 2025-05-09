#!/bin/bash

# Script untuk memperbaiki masalah koneksi di haxorport-go-client
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
        print_error "File konfigurasi tidak ditemukan di $CONFIG_FILE"
    fi
}

# Perbaiki konfigurasi
fix_config() {
    print_info "Memperbaiki konfigurasi koneksi..."
    
    # Backup konfigurasi lama
    cp "$CONFIG_FILE" "$CONFIG_FILE.bak"
    print_info "Backup konfigurasi lama disimpan di $CONFIG_FILE.bak"
    
    # Minta token dari pengguna
    print_info "Silakan masukkan token autentikasi Anda dari dashboard Haxorport:"
    read -p "Token: " USER_TOKEN
    
    if [ -z "$USER_TOKEN" ]; then
        print_error "Token tidak boleh kosong. Silakan jalankan script ini lagi dan masukkan token yang valid."
    fi
    
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
log_file: "logs/haxor-client.log"
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
    $HAXOR_CMD version
    
    print_info "Koneksi berhasil diuji. Jika tidak ada error, koneksi sudah berhasil diperbaiki."
    print_info "Anda sekarang dapat mencoba menjalankan perintah HTTP tunnel:"
    print_info "$HAXOR_CMD http --port 80"
}

# Main script
echo -e "${GREEN}=== Haxorport Connection Fix Script ===${NC}"

# Inisialisasi
detect_config_path
fix_config
test_connection

# Tampilkan informasi
echo -e "\n${GREEN}==================================================${NC}"
echo -e "${GREEN}✅ Konfigurasi koneksi berhasil diperbaiki!${NC}"
echo -e "${GREEN}==================================================${NC}"
echo -e "Jika masih mengalami masalah, coba periksa:"
echo -e "1. Koneksi internet Anda"
echo -e "2. Firewall atau proxy yang mungkin memblokir koneksi WebSocket"
echo -e "3. Status server Haxorport"
echo -e "${GREEN}==================================================${NC}"
