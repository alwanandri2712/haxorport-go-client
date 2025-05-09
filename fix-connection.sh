#!/bin/bash

# Script untuk memperbaiki masalah koneksi di haxorport-go-client
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
        print_error "File konfigurasi tidak ditemukan di $CONFIG_FILE"
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

# Perbaiki konfigurasi
fix_config() {
    print_info "Memperbaiki konfigurasi koneksi..."
    
    # Periksa koneksi internet dan firewall
    check_internet
    check_firewall
    
    # Backup konfigurasi lama
    cp "$CONFIG_FILE" "$CONFIG_FILE.bak"
    print_info "Backup konfigurasi lama disimpan di $CONFIG_FILE.bak"
    
    # Minta token dari pengguna
    print_info "Silakan masukkan token autentikasi Anda dari dashboard Haxorport:"
    read -p "Token: " USER_TOKEN
    
    if [ -z "$USER_TOKEN" ]; then
        print_error "Token tidak boleh kosong. Silakan jalankan script ini lagi dan masukkan token yang valid."
    fi
    
    # Buat direktori logs jika belum ada
    mkdir -p "$(dirname "$CONFIG_FILE")/logs"
    
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
        DEBUG_CONFIG="$(dirname "$CONFIG_FILE")/debug.yaml"
        
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
    print_info "Jika masih mengalami masalah, coba hubungi support Haxorport."
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
