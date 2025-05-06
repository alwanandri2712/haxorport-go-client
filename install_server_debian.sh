#!/bin/bash

# Installer untuk haxorport-go-server pada Linux Debian
# Script ini akan menginstal haxorport-go-server dan semua dependensinya

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
    print_info "Memeriksa repositori server..."
    
    # Jika direktori sudah ada, update saja
    if [ -d "haxorport-go-server" ]; then
        print_info "Repositori server sudah ada. Memperbarui..."
        cd haxorport-go-server
        git pull
        cd ..
    else
        print_info "Mengkloning repositori server..."
        git clone https://github.com/alwanandri2712/haxorport-go-server.git
        
        if [ ! -d "haxorport-go-server" ]; then
            print_error "Gagal mengkloning repositori server."
            exit 1
        else
            print_success "Repositori server berhasil dikloning."
        fi
    fi
}

# Fungsi untuk mengkompilasi aplikasi server
build_application() {
    print_info "Mengkompilasi aplikasi server..."
    
    cd haxorport-go-server
    
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
        go build -o bin/haxor-server main.go
    fi
    
    # Periksa apakah binary berhasil dibuat
    if [ ! -f "bin/haxor-server" ]; then
        print_error "Gagal mengkompilasi aplikasi server."
        exit 1
    else
        print_success "Aplikasi server berhasil dikompilasi."
    fi
    
    cd ..
}

# Fungsi untuk menginstal aplikasi server
install_application() {
    print_info "Menginstal aplikasi server..."
    
    # Buat direktori instalasi
    sudo mkdir -p /opt/haxorport-server
    
    # Salin file yang diperlukan
    sudo cp -r haxorport-go-server/bin /opt/haxorport-server/
    
    # Buat direktori untuk log
    sudo mkdir -p /opt/haxorport-server/logs
    
    # Buat script service
    cat > haxorport-server.service << EOF
[Unit]
Description=Haxorport Server
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/opt/haxorport-server
ExecStart=/opt/haxorport-server/bin/haxor-server
Restart=on-failure
RestartSec=5
StandardOutput=append:/opt/haxorport-server/logs/haxor-server.log
StandardError=append:/opt/haxorport-server/logs/haxor-server-error.log

[Install]
WantedBy=multi-user.target
EOF
    
    # Salin file service ke systemd
    sudo cp haxorport-server.service /etc/systemd/system/
    
    # Reload systemd
    sudo systemctl daemon-reload
    
    # Aktifkan service
    sudo systemctl enable haxorport-server.service
    
    # Mulai service
    sudo systemctl start haxorport-server.service
    
    # Periksa status service
    sudo systemctl status haxorport-server.service
    
    # Periksa apakah instalasi berhasil
    if ! systemctl is-active --quiet haxorport-server.service; then
        print_warning "Service tidak berjalan. Silakan periksa log untuk informasi lebih lanjut."
    else
        print_success "Service berhasil dijalankan."
    fi
}

# Fungsi untuk menampilkan informasi penggunaan
show_usage() {
    echo -e "\n${GREEN}Haxorport Server berhasil diinstal!${NC}"
    echo -e "\nInformasi service:"
    echo -e "  Status: ${BLUE}sudo systemctl status haxorport-server.service${NC}"
    echo -e "  Start: ${BLUE}sudo systemctl start haxorport-server.service${NC}"
    echo -e "  Stop: ${BLUE}sudo systemctl stop haxorport-server.service${NC}"
    echo -e "  Restart: ${BLUE}sudo systemctl restart haxorport-server.service${NC}"
    echo -e "  Log: ${BLUE}sudo journalctl -u haxorport-server.service${NC}"
    echo -e "\nUntuk menghapus instalasi, jalankan: ${YELLOW}sudo systemctl stop haxorport-server.service && sudo systemctl disable haxorport-server.service && sudo rm -rf /opt/haxorport-server /etc/systemd/system/haxorport-server.service${NC}"
}

# Main script
echo -e "${GREEN}=== Installer Haxorport Server untuk Debian ===${NC}"
echo -e "Installer ini akan menginstal Haxorport Server dan semua dependensinya.\n"

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
