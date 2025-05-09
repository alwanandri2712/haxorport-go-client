#!/bin/bash

# Colors for output formatting
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Functions for displaying colored messages
print_info() {
    echo -e "${YELLOW}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_debug() {
    echo -e "${BLUE}[DEBUG]${NC} $1"
}

print_header() {
    echo -e "\n${GREEN}==================================================${NC}"
    echo -e "${GREEN}$1${NC}"
    echo -e "${GREEN}==================================================${NC}"
}

# Detect platform
detect_platform() {
    # Detect OS
    case "$(uname -s)" in
        Linux*)     
            OS="linux" 
            CONFIG_DIR="$HOME/.haxorport"
            if [ "$(id -u)" -eq 0 ]; then
                # If root user
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
    
    # Detect architecture
    case "$(uname -m)" in
        x86_64|amd64)  ARCH="amd64" ;;
        arm64|aarch64) ARCH="arm64" ;;
        *)             ARCH="unknown" ;;
    esac
    
    print_info "Detected platform: $OS/$ARCH"
}

# Check if Go is installed
check_go() {
    if ! command -v go &> /dev/null; then
        print_error "Go is not installed. Please install Go first: https://golang.org/doc/install"
        exit 1
    fi
    
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    print_info "Go version: $GO_VERSION"
}

# Build application
build_app() {
    print_info "Downloading dependencies..."
    go mod download || print_error "Failed to download dependencies"
    
    print_info "Building application..."
    if [ "$OS" = "windows" ]; then
        go build -o bin/haxor.exe main.go || print_error "Build failed"
    else
        go build -o bin/haxor main.go || print_error "Build failed"
    fi
    
    # Give execution permission to binary (for Linux and macOS)
    if [ "$OS" != "windows" ]; then
        chmod +x bin/haxor
    fi
    
    print_success "Build successful!"
}

# Create default configuration
create_default_config() {
    CONFIG_FILE="$CONFIG_DIR/config.yaml"
    
    # Create config directory if it doesn't exist
    if [ ! -d "$CONFIG_DIR" ]; then
        print_info "Creating configuration directory: $CONFIG_DIR"
        mkdir -p "$CONFIG_DIR"
    fi
    
    # Create logs directory
    mkdir -p "$CONFIG_DIR/logs"
    
    # Check if config file exists
    if [ ! -f "$CONFIG_FILE" ]; then
        print_info "Creating default configuration file at $CONFIG_FILE"
        
        # Create default config file
        cat > "$CONFIG_FILE" << EOF
# Haxorport Client Configuration
server_address: control.haxorport.online
control_port: 443
data_port: 0

# Authentication Configuration
auth_enabled: true
auth_token: ""  # Add your token here
auth_validation_url: https://haxorport.online/AuthToken/validate

# TLS Configuration
tls_enabled: true
tls_cert: ""
tls_key: ""

# Base domain for tunnel subdomains
base_domain: "haxorport.online"

# Logging Configuration
log_level: warn
log_file: "$CONFIG_DIR/logs/haxor-client.log"
tunnels: []
EOF
        
        print_success "Default configuration created"
        print_info "IMPORTANT: You need to edit $CONFIG_FILE and add your auth token"
    else
        print_info "Configuration file already exists at $CONFIG_FILE"
    fi
}

# Main function
main() {
    print_header "Haxorport Client Setup"
    
    # Detect platform
    detect_platform
    
    # Check Go
    check_go
    
    # Create bin directory
    mkdir -p bin
    
    # Build application
    build_app
    
    # Create default configuration
    create_default_config
    
    # Show success message
    print_header "✅ Haxorport Client successfully built!"
    echo -e "Binary location: $(pwd)/bin/"
    echo -e ""
    
    # Check if haxor is already installed globally
    if command -v haxor &> /dev/null; then
        echo -e "Haxor is already installed globally. You can use:"
        echo -e "  haxor --help"
        echo -e ""
        echo -e "Or use the newly built binary:"
    fi
    
    echo -e "To run the newly built binary:"
    if [ "$OS" = "windows" ]; then
        echo -e "  ./bin/haxor.exe --help"
    else
        echo -e "  ./bin/haxor --help"
    fi
    echo -e ""
    echo -e "Example usage:"
    echo -e "  haxor http --port 80"
    echo -e "  haxor tcp --local-port 22 --remote-port 2222"
    echo -e "${GREEN}==================================================${NC}"
    echo -e "IMPORTANT: Before using, make sure to edit your config file and add your auth token:"
    echo -e "  $CONFIG_FILE"
    echo -e "${GREEN}==================================================${NC}"
}

# Run main function
main
