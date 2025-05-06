#!/bin/bash

# Installer untuk haxorport pada Linux Debian
# Script ini akan menginstal haxorport dan semua dependensinya

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

# Fungsi untuk memeriksa apakah command tersedia
check_command() {
    if ! command -v $1 &> /dev/null; then
        print_error "$1 tidak ditemukan. Menginstal $1..."
        return 1
    else
        print_info "$1 sudah terinstal."
        return 0
    fi
}

# Fungsi untuk menginstal dependensi
install_dependencies() {
    print_info "Memeriksa dan menginstal dependensi..."
    
    # Update package list
    print_info "Memperbarui daftar paket..."
    sudo apt-get update
    
    # Instal Go jika belum ada
    if ! check_command go; then
        print_info "Menginstal Go..."
        sudo apt-get install -y golang-go
        
        if ! check_command go; then
            print_error "Gagal menginstal Go. Silakan instal secara manual."
            exit 1
        else
            print_success "Go berhasil diinstal."
        fi
    fi
    
    # Instal Git jika belum ada
    if ! check_command git; then
        print_info "Menginstal Git..."
        sudo apt-get install -y git
        
        if ! check_command git; then
            print_error "Gagal menginstal Git. Silakan instal secara manual."
            exit 1
        else
            print_success "Git berhasil diinstal."
        fi
    fi
    
    # Instal build-essential jika belum ada
    if ! dpkg -l | grep -q build-essential; then
        print_info "Menginstal build-essential..."
        sudo apt-get install -y build-essential
        
        if ! dpkg -l | grep -q build-essential; then
            print_error "Gagal menginstal build-essential. Silakan instal secara manual."
            exit 1
        else
            print_success "build-essential berhasil diinstal."
        fi
    else
        print_info "build-essential sudah terinstal."
    fi
    
    print_success "Semua dependensi berhasil diinstal."
}

# Fungsi untuk mengkloning repositori
clone_repository() {
    print_info "Memeriksa repositori..."
    
    # Jika direktori sudah ada, update saja
    if [ -d "haxorport-go-client" ]; then
        print_info "Repositori sudah ada. Memperbarui..."
        cd haxorport-go-client
        git pull
        cd ..
    else
        print_info "Mengkloning repositori..."
        git clone https://github.com/alwanandri2712/haxorport-go-client.git
        
        if [ ! -d "haxorport-go-client" ]; then
            print_error "Gagal mengkloning repositori."
            exit 1
        else
            print_success "Repositori berhasil dikloning."
        fi
    fi
}

# Fungsi untuk mengkompilasi aplikasi
build_application() {
    print_info "Mengkompilasi aplikasi..."
    
    cd haxorport-go-client
    
    # Buat direktori bin jika belum ada
    mkdir -p bin
    
    # Jalankan script build jika ada
    if [ -f "scripts/build.sh" ]; then
        print_info "Menjalankan script build..."
        chmod +x scripts/build.sh
        ./scripts/build.sh
    else
        # Kompilasi manual jika script build tidak ada
        print_info "Mengkompilasi secara manual..."
        go build -o bin/haxor main.go
    fi
    
    # Periksa apakah binary berhasil dibuat
    if [ ! -f "bin/haxor" ]; then
        print_error "Gagal mengkompilasi aplikasi."
        exit 1
    else
        print_success "Aplikasi berhasil dikompilasi."
    fi
    
    # Pastikan script wrapper memiliki izin eksekusi
    chmod +x haxorport
    
    cd ..
}

# Fungsi untuk menginstal aplikasi
install_application() {
    print_info "Menginstal aplikasi..."
    
    # Buat direktori instalasi
    sudo mkdir -p /opt/haxorport
    
    # Salin file yang diperlukan
    sudo cp -r haxorport-go-client/bin /opt/haxorport/
    sudo cp haxorport-go-client/haxorport /opt/haxorport/
    
    # Buat symlink ke /usr/local/bin
    sudo ln -sf /opt/haxorport/haxorport /usr/local/bin/haxorport
    
    # Periksa apakah instalasi berhasil
    if [ ! -f "/usr/local/bin/haxorport" ]; then
        print_error "Gagal menginstal aplikasi."
        exit 1
    else
        print_success "Aplikasi berhasil diinstal."
    fi
}

# Fungsi untuk menampilkan informasi penggunaan
show_usage() {
    echo -e "\n${GREEN}Haxorport berhasil diinstal!${NC}"
    echo -e "\nCara penggunaan:"
    echo -e "  ${BLUE}haxorport http http://localhost:9090${NC} - Membuat HTTP tunnel ke localhost:9090"
    echo -e "  ${BLUE}haxorport tcp 22${NC} - Membuat TCP tunnel ke port 22"
    echo -e "  ${BLUE}haxorport --help${NC} - Menampilkan bantuan"
    echo -e "\nUntuk menghapus instalasi, jalankan: ${YELLOW}sudo rm -rf /opt/haxorport /usr/local/bin/haxorport${NC}"
}

# Main script
echo -e "${GREEN}=== Installer Haxorport untuk Debian ===${NC}"
echo -e "Installer ini akan menginstal Haxorport dan semua dependensinya.\n"

# Tanya konfirmasi
read -p "Lanjutkan instalasi? (y/n): " confirm
if [[ $confirm != [yY] ]]; then
    print_warning "Instalasi dibatalkan."
    exit 0
fi

# Jalankan fungsi-fungsi instalasi
install_dependencies
clone_repository
build_application
install_application
show_usage

exit 0
