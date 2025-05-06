#!/bin/bash

# Script untuk build Haxorport Client

set -e

# Direktori root proyek
ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

# Versi
VERSION=$(grep 'const Version = "' cmd/version.go | cut -d'"' -f2)
echo "Building Haxorport Client v$VERSION"

# Bersihkan direktori build
echo "Cleaning build directory..."
rm -rf bin
mkdir -p bin

# Build untuk platform saat ini
echo "Building for current platform..."
go build -o bin/haxor main.go

# Build untuk platform lain jika diperlukan
if [ "$1" == "--all" ]; then
    echo "Building for all platforms..."

    # Linux (amd64)
    echo "Building for Linux (amd64)..."
    GOOS=linux GOARCH=amd64 go build -o bin/haxor-linux-amd64 main.go

    # Linux (arm64)
    echo "Building for Linux (arm64)..."
    GOOS=linux GOARCH=arm64 go build -o bin/haxor-linux-arm64 main.go

    # Windows (amd64)
    echo "Building for Windows (amd64)..."
    GOOS=windows GOARCH=amd64 go build -o bin/haxor-windows-amd64.exe main.go

    # macOS (amd64)
    echo "Building for macOS (amd64)..."
    GOOS=darwin GOARCH=amd64 go build -o bin/haxor-darwin-amd64 main.go

    # macOS (arm64)
    echo "Building for macOS (arm64)..."
    GOOS=darwin GOARCH=arm64 go build -o bin/haxor-darwin-arm64 main.go
fi

echo "Build completed successfully!"
echo "Binary location: $ROOT_DIR/bin/haxor"
