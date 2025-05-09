#!/bin/bash

# Installer untuk haxorport - Mendukung berbagai sistem operasi
# Script ini akan menginstal haxorport dan semua dependensinya

# Warna untuk output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Variabel global
REPO_URL="https://github.com/alwanandri2712/haxorport-go-client.git"
REPO_BRANCH="https-tunnel-support"
INSTALL_DIR="/opt/haxorport"
CONFIG_DIR="/etc/haxorport"
BIN_DIR="/usr/local/bin"
LOG_DIR="/var/log/haxorport"

# Untuk macOS
if [[ "$OSTYPE" == "darwin"* ]]; then
    INSTALL_DIR="$HOME/Library/Application Support/haxorport"
    CONFIG_DIR="$HOME/Library/Preferences/haxorport"
    LOG_DIR="$HOME/Library/Logs/haxorport"
    
    # Periksa apakah menggunakan Apple Silicon
    if [[ $(uname -m) == "arm64" ]]; then
        print_info "Terdeteksi Apple Silicon (M1/M2)"
        # Gunakan /opt/homebrew/bin untuk Apple Silicon jika ada
        if [ -d "/opt/homebrew/bin" ]; then
            BIN_DIR="/opt/homebrew/bin"
        else
            BIN_DIR="/usr/local/bin"
        fi
    else
        BIN_DIR="/usr/local/bin"
    fi
fi

# Untuk Windows (WSL)
if grep -q Microsoft /proc/version 2>/dev/null; then
    INSTALL_DIR="$HOME/.haxorport"
    CONFIG_DIR="$HOME/.haxorport/config"
    LOG_DIR="$HOME/.haxorport/logs"
    BIN_DIR="$HOME/.local/bin"
    mkdir -p "$BIN_DIR"
    export PATH="$PATH:$BIN_DIR"
    
    # Tambahkan BIN_DIR ke PATH secara permanen jika belum ada
    if ! grep -q "$BIN_DIR" "$HOME/.bashrc" 2>/dev/null; then
        print_info "Menambahkan $BIN_DIR ke PATH di .bashrc"
        echo "export PATH=\"$PATH:$BIN_DIR\"" >> "$HOME/.bashrc"
    fi
    
    # Jika menggunakan zsh, tambahkan juga ke .zshrc
    if [ -f "$HOME/.zshrc" ] && ! grep -q "$BIN_DIR" "$HOME/.zshrc" 2>/dev/null; then
        print_info "Menambahkan $BIN_DIR ke PATH di .zshrc"
        echo "export PATH=\"$PATH:$BIN_DIR\"" >> "$HOME/.zshrc"
    fi
fi

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
        print_error "$1 tidak ditemukan."
        return 1
    else
        print_info "$1 sudah terinstal."
        return 0
    fi
}

# Fungsi untuk menginstal dependensi berdasarkan OS
install_dependencies() {
    print_info "Memeriksa dan menginstal dependensi..."
    
    # Periksa Go
    if ! check_command go; then
        print_info "Go tidak ditemukan. Menginstal Go..."
        
        # Deteksi OS dan instal Go
        if [[ "$OSTYPE" == "linux-gnu"* ]]; then
            # Deteksi package manager
            if command -v apt-get &> /dev/null; then
                # Debian/Ubuntu
                sudo apt-get update
                sudo apt-get install -y golang-go
            elif command -v yum &> /dev/null; then
                # CentOS/RHEL
                sudo yum install -y golang
            elif command -v pacman &> /dev/null; then
                # Arch Linux
                sudo pacman -Sy go
            elif command -v zypper &> /dev/null; then
                # openSUSE
                sudo zypper install -y go
            else
                print_error "Package manager tidak dikenali. Silakan instal Go secara manual."
                exit 1
            fi
        elif [[ "$OSTYPE" == "darwin"* ]]; then
            # macOS
            if command -v brew &> /dev/null; then
                brew install go
            else
                print_error "Homebrew tidak ditemukan. Silakan instal Homebrew terlebih dahulu."
                print_info "Jalankan: /bin/bash -c \"\$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)\""
                
                # Tambahkan instruksi khusus untuk Apple Silicon
                if [[ $(uname -m) == "arm64" ]]; then
                    print_info "Untuk Apple Silicon (M1/M2), setelah menginstal Homebrew, jalankan:"
                    print_info "echo 'eval \"\$(/opt/homebrew/bin/brew shellenv)\"' >> ~/.zprofile"
                    print_info "eval \"\$(/opt/homebrew/bin/brew shellenv)\""
                fi
                
                exit 1
            fi
        else
            print_error "Sistem operasi tidak didukung untuk instalasi otomatis. Silakan instal Go secara manual."
            exit 1
        fi
    fi
    
    # Periksa Git
    if ! check_command git; then
        print_info "Git tidak ditemukan. Menginstal Git..."
        
        # Deteksi OS dan instal Git
        if [[ "$OSTYPE" == "linux-gnu"* ]]; then
            # Deteksi package manager
            if command -v apt-get &> /dev/null; then
                # Debian/Ubuntu
                sudo apt-get update
                sudo apt-get install -y git
            elif command -v yum &> /dev/null; then
                # CentOS/RHEL
                sudo yum install -y git
            elif command -v pacman &> /dev/null; then
                # Arch Linux
                sudo pacman -Sy git
            elif command -v zypper &> /dev/null; then
                # openSUSE
                sudo zypper install -y git
            else
                print_error "Package manager tidak dikenali. Silakan instal Git secara manual."
                exit 1
            fi
        elif [[ "$OSTYPE" == "darwin"* ]]; then
            # macOS
            if command -v brew &> /dev/null; then
                brew install git
            else
                print_error "Homebrew tidak ditemukan. Silakan instal Homebrew terlebih dahulu."
                exit 1
            fi
        else
            print_error "Sistem operasi tidak didukung untuk instalasi otomatis. Silakan instal Git secara manual."
            exit 1
        fi
    fi
    
    print_success "Semua dependensi berhasil diinstal."
}

# Fungsi untuk menyiapkan repositori
setup_repository() {
    print_info "Menyiapkan repositori..."
    
    # Periksa apakah script dijalankan dari dalam repositori
    CURRENT_DIR=$(basename "$PWD")
    if [ "$CURRENT_DIR" = "haxorport-go-client" ]; then
        print_info "Script dijalankan dari dalam repositori. Menggunakan repositori saat ini."
        
        # Pastikan kita berada di branch yang benar
        CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
        if [ "$CURRENT_BRANCH" != "$REPO_BRANCH" ]; then
            print_info "Checkout ke branch $REPO_BRANCH..."
            git checkout $REPO_BRANCH
            git pull origin $REPO_BRANCH
        else
            print_info "Sudah berada di branch $REPO_BRANCH. Memperbarui..."
            git pull origin $REPO_BRANCH
        fi
        
        # Gunakan direktori saat ini
        REPO_DIR="$PWD"
    else
        # Jika direktori repositori sudah ada di direktori saat ini, gunakan itu
        if [ -d "haxorport-go-client" ]; then
            print_info "Repositori sudah ada. Memperbarui..."
            cd haxorport-go-client
            git fetch --all
            
            # Checkout ke branch yang ditentukan
            print_info "Checkout ke branch $REPO_BRANCH..."
            git checkout $REPO_BRANCH
            git pull origin $REPO_BRANCH
            
            REPO_DIR="$PWD"
            cd ..
        else
            print_info "Mengkloning repositori..."
            # Clone langsung dari branch yang ditentukan
            git clone -b $REPO_BRANCH $REPO_URL
            
            if [ ! -d "haxorport-go-client" ]; then
                print_error "Gagal mengkloning repositori."
                exit 1
            else
                print_success "Repositori berhasil dikloning."
                REPO_DIR="$PWD/haxorport-go-client"
            fi
        fi
    fi
}

# Fungsi untuk mengkompilasi aplikasi
build_application() {
    print_info "Mengkompilasi aplikasi..."
    
    # Masuk ke direktori repositori
    cd "$REPO_DIR"
    
    # Buat direktori bin jika belum ada
    mkdir -p bin
    
    # Jalankan script build
    print_info "Menjalankan script build..."
    chmod +x scripts/build.sh
    ./scripts/build.sh
    
    # Periksa apakah binary berhasil dibuat
    if [ ! -f "bin/haxor" ]; then
        print_error "Gagal mengkompilasi aplikasi."
        exit 1
    else
        print_success "Aplikasi berhasil dikompilasi."
    fi
    
    cd ..
}

# Fungsi untuk menginstal aplikasi
install_application() {
    print_info "Menginstal aplikasi..."
    
    cd "$REPO_DIR"
    
    # Buat direktori instalasi
    mkdir -p "$INSTALL_DIR"
    mkdir -p "$CONFIG_DIR"
    mkdir -p "$LOG_DIR"
    
    # Salin file yang diperlukan
    cp -r bin "$INSTALL_DIR/"
    
    # Salin konfigurasi produksi
    cp config.prod.yaml "$CONFIG_DIR/config.yaml"
    
    # Sesuaikan konfigurasi log untuk OS
    if [[ "$OSTYPE" == "darwin"* ]] || grep -q Microsoft /proc/version 2>/dev/null; then
        # Untuk macOS dan WSL, gunakan path relatif ke home
        sed -i.bak "s|log_file:.*|log_file: \"$LOG_DIR/haxor-client.log\"|g" "$CONFIG_DIR/config.yaml"
        rm -f "$CONFIG_DIR/config.yaml.bak" 2>/dev/null
    else
        # Untuk Linux
        sed -i "s|log_file:.*|log_file: \"$LOG_DIR/haxor-client.log\"|g" "$CONFIG_DIR/config.yaml"
    fi
    
    # Buat script wrapper untuk haxorport
    cat > haxorport << EOF
#!/bin/bash
"$INSTALL_DIR/bin/haxor" --config "$CONFIG_DIR/config.yaml" "\$@"
EOF
    
    # Beri izin eksekusi pada script wrapper
    chmod +x haxorport
    
    # Salin script wrapper ke direktori bin
    if [[ "$OSTYPE" == "darwin"* ]] || grep -q Microsoft /proc/version 2>/dev/null; then
        # Untuk macOS dan WSL, tidak perlu sudo
        cp haxorport "$BIN_DIR/"
    else
        # Untuk Linux
        sudo cp haxorport "$BIN_DIR/"
    fi
    
    # Periksa apakah instalasi berhasil
    if [ ! -f "$BIN_DIR/haxorport" ]; then
        print_error "Gagal menginstal aplikasi."
        exit 1
    else
        print_success "Aplikasi berhasil diinstal."
    fi
    
    cd ..
}

# Fungsi untuk memperbarui aplikasi
update_application() {
    print_info "Memperbarui aplikasi..."
    
    # Update repositori
    setup_repository
    
    # Build ulang aplikasi
    build_application
    
    # Instal ulang aplikasi
    install_application
    
    print_success "Aplikasi berhasil diperbarui."
}

# Fungsi untuk menampilkan informasi penggunaan
show_usage() {
    echo -e "\n${GREEN}Haxorport berhasil diinstal!${NC}"
    echo -e "\nCara penggunaan:"
    echo -e "  ${BLUE}haxorport http http://localhost:8080${NC} - Membuat HTTP tunnel ke localhost:8080"
    echo -e "  ${BLUE}haxorport tcp 22${NC} - Membuat TCP tunnel ke port 22"
    echo -e "  ${BLUE}haxorport --help${NC} - Menampilkan bantuan"
    
    echo -e "\nKonfigurasi:"
    echo -e "  File konfigurasi: ${YELLOW}$CONFIG_DIR/config.yaml${NC}"
    echo -e "  Log file: ${YELLOW}$LOG_DIR/haxor-client.log${NC}"
    
    echo -e "\nUntuk memperbarui aplikasi, jalankan: ${YELLOW}$0 --update${NC}"
    
    if [[ "$OSTYPE" == "darwin"* ]] || grep -q Microsoft /proc/version 2>/dev/null; then
        echo -e "Untuk menghapus instalasi, jalankan: ${YELLOW}rm -rf $INSTALL_DIR $CONFIG_DIR $BIN_DIR/haxorport${NC}"
    else
        echo -e "Untuk menghapus instalasi, jalankan: ${YELLOW}sudo rm -rf $INSTALL_DIR $CONFIG_DIR $BIN_DIR/haxorport${NC}"
    fi
}

# Main script
echo -e "${GREEN}=== Installer Haxorport Multi-Platform ===${NC}"
echo -e "Installer ini akan menginstal Haxorport dan semua dependensinya.\n"
echo -e "Terdeteksi sistem operasi: ${YELLOW}$OSTYPE${NC}"

# Periksa apakah ini adalah permintaan update
if [ "$1" == "--update" ]; then
    update_application
    show_usage
    exit 0
fi

# Tanya konfirmasi
read -p "Lanjutkan instalasi? (y/n): " confirm
if [[ $confirm != [yY] ]]; then
    print_warning "Instalasi dibatalkan."
    exit 0
fi

# Jalankan fungsi-fungsi instalasi
install_dependencies
setup_repository
build_application
install_application
show_usage

exit 0
