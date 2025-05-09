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

# Detect configuration directory
if [ "$(id -u)" -eq 0 ]; then
    # If root user
    CONFIG_DIR="/root/.haxorport"
else
    CONFIG_DIR="$HOME/.haxorport"
fi

CONFIG_FILE="$CONFIG_DIR/config.yaml"

# Check if config file exists
if [ ! -f "$CONFIG_FILE" ]; then
    print_error "Configuration file not found at $CONFIG_FILE"
    exit 1
fi

print_info "Found configuration file at $CONFIG_FILE"

# Backup the current config
print_info "Creating backup of current configuration..."
cp "$CONFIG_FILE" "${CONFIG_FILE}.bak"
print_success "Backup created at ${CONFIG_FILE}.bak"

# Fix configuration
print_info "Updating configuration to fix websocket handshake issue..."

# Ask for token from user
print_info "Please enter your authentication token from the Haxorport dashboard:"
read -p "Token: " USER_TOKEN

if [ -z "$USER_TOKEN" ]; then
    print_error "Token cannot be empty. Please run this script again and enter a valid token."
    exit 1
fi

# Create logs directory if it doesn't exist
mkdir -p "$CONFIG_DIR/logs"

# Update configuration
cat > "$CONFIG_FILE" << EOF
# Haxorport Client Configuration
server_address: control.haxorport.online
control_port: 443
data_port: 0

# Authentication Configuration
auth_enabled: true
auth_token: $USER_TOKEN
auth_validation_url: https://haxorport.online/AuthToken/validate

# TLS Configuration
tls_enabled: true
tls_cert: ""
tls_key: ""

# Base domain for tunnel subdomains
base_domain: "haxorport.online"

# Logging Configuration
log_level: debug
log_file: "$CONFIG_DIR/logs/haxor-client.log"
tunnels: []
EOF

print_success "Configuration updated successfully!"

# Test connection
print_info "Testing connection to server..."
print_info "Running: haxor version"
haxor version

if [ $? -eq 0 ]; then
    print_success "Connection successful!"
else
    print_info "Connection test failed. Let's try to diagnose the issue..."
    
    # Check internet connection
    print_info "Checking internet connection..."
    if ping -c 1 google.com &> /dev/null; then
        print_success "Internet connection OK"
    else
        print_error "No internet connection. Please check your network connection."
        exit 1
    fi
    
    # Check access to haxorport server
    print_info "Checking access to haxorport server..."
    if ping -c 1 control.haxorport.online &> /dev/null; then
        print_success "Can reach haxorport server"
    else
        print_error "Cannot reach haxorport server. The server may be down or blocked by your firewall."
        exit 1
    fi
    
    # Check if port 443 is accessible
    print_info "Checking if port 443 is accessible..."
    if nc -z -w 5 control.haxorport.online 443 &> /dev/null; then
        print_success "Port 443 is accessible"
    else
        print_error "Port 443 is not accessible. Firewall may be blocking the connection."
        exit 1
    fi
    
    # Check TLS connection
    print_info "Checking TLS connection to server..."
    if curl -s --head https://control.haxorport.online | grep "HTTP" > /dev/null; then
        print_success "TLS connection to server is working"
    else
        print_error "Cannot establish TLS connection to server."
        exit 1
    fi
    
    # Check token validation
    print_info "Checking token validation..."
    VALIDATION_RESULT=$(curl -s -X POST -H "Content-Type: application/json" -d "{\"token\":\"$USER_TOKEN\"}" https://haxorport.online/AuthToken/validate)
    print_debug "Validation result: $VALIDATION_RESULT"
    
    if echo "$VALIDATION_RESULT" | grep -q "valid"; then
        print_success "Token validation successful"
    else
        print_error "Token validation failed. Please check your token and try again."
        exit 1
    fi
fi

print_info "You can now try running: haxor http --port 80"
