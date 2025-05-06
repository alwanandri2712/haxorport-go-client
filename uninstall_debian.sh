#!/bin/bash

# Uninstaller untuk haxorport pada Linux Debian
# Script ini akan menghapus instalasi haxorport

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
    print_info "Menghapus instalasi haxorport..."
    
    # Hapus symlink
    if [ -L "/usr/local/bin/haxorport" ]; then
        sudo rm -f /usr/local/bin/haxorport
        print_info "Symlink dihapus."
    else
        print_warning "Symlink tidak ditemukan."
    fi
    
    # Hapus direktori instalasi
    if [ -d "/opt/haxorport" ]; then
        sudo rm -rf /opt/haxorport
        print_info "Direktori instalasi dihapus."
    else
        print_warning "Direktori instalasi tidak ditemukan."
    fi
    
    print_success "Haxorport berhasil dihapus dari sistem."
}

# Main script
echo -e "${RED}=== Uninstaller Haxorport untuk Debian ===${NC}"
echo -e "Uninstaller ini akan menghapus Haxorport dari sistem Anda.\n"

# Tanya konfirmasi
read -p "Lanjutkan penghapusan? (y/n): " confirm
if [[ $confirm != [yY] ]]; then
    print_warning "Penghapusan dibatalkan."
    exit 0
fi

# Jalankan fungsi penghapusan
uninstall_application

echo -e "\n${GREEN}Haxorport berhasil dihapus dari sistem.${NC}"
echo -e "Jika Anda ingin menghapus repositori lokal, jalankan: ${YELLOW}rm -rf haxorport-go-client${NC}"

exit 0
