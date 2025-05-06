#!/bin/bash

# Uninstaller untuk haxorport-go-server pada Linux Debian
# Script ini akan menghapus instalasi haxorport-go-server

# Warna untuk output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Fungsi untuk menampilkan pesan dengan warna
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Fungsi untuk menghapus instalasi
uninstall_application() {
    print_info "Menghapus instalasi haxorport-go-server..."
    
    # Hentikan dan nonaktifkan service
    if systemctl is-active --quiet haxorport-server.service; then
        print_info "Menghentikan service..."
        sudo systemctl stop haxorport-server.service
    fi
    
    if systemctl is-enabled --quiet haxorport-server.service; then
        print_info "Menonaktifkan service..."
        sudo systemctl disable haxorport-server.service
    fi
    
    # Hapus file service
    if [ -f "/etc/systemd/system/haxorport-server.service" ]; then
        print_info "Menghapus file service..."
        sudo rm -f /etc/systemd/system/haxorport-server.service
        sudo systemctl daemon-reload
    else
        print_warning "File service tidak ditemukan."
    fi
    
    # Hapus direktori instalasi
    if [ -d "/opt/haxorport-server" ]; then
        print_info "Menghapus direktori instalasi..."
        sudo rm -rf /opt/haxorport-server
    else
        print_warning "Direktori instalasi tidak ditemukan."
    fi
    
    print_success "Haxorport Server berhasil dihapus dari sistem."
}

# Main script
echo -e "${RED}=== Uninstaller Haxorport Server untuk Debian ===${NC}"
echo -e "Uninstaller ini akan menghapus Haxorport Server dari sistem Anda.\n"

# Tanya konfirmasi
read -p "Lanjutkan penghapusan? (y/n): " confirm
if [[ $confirm != [yY] ]]; then
    print_warning "Penghapusan dibatalkan."
    exit 0
fi

# Jalankan fungsi penghapusan
uninstall_application

echo -e "\n${GREEN}Haxorport Server berhasil dihapus dari sistem.${NC}"
echo -e "Jika Anda ingin menghapus repositori lokal, jalankan: ${YELLOW}rm -rf haxorport-go-server${NC}"

exit 0
