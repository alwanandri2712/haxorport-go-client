#!/bin/bash
# Build and package haxorport-go-client for Linux distribution
# Author: alwanandri2712
# Usage: bash build-package.sh
set -e

APP_NAME="haxorport"
DIST_DIR="dist"
OUTPUT_TAR="${APP_NAME}-go-client-linux-amd64.tar.gz"

# Clean previous build
echo "[INFO] Cleaning previous build..."
rm -rf "$DIST_DIR" "$OUTPUT_TAR"
mkdir -p "$DIST_DIR"

# Build binary for Linux amd64
echo "[INFO] Building binary for Linux amd64..."
GOOS=linux GOARCH=amd64 go build -o "$DIST_DIR/$APP_NAME" main.go

# Copy config and installer script
echo "[INFO] Copying config.yaml and install.sh..."
cp config.yaml "$DIST_DIR/"
cp install.sh "$DIST_DIR/"

# (Optional) Copy additional scripts or docs if needed
# cp README.md "$DIST_DIR/"

# Create tar.gz package
echo "[INFO] Creating tar.gz package..."
tar -czvf "$OUTPUT_TAR" -C "$DIST_DIR" .

# Success message
echo "[SUCCESS] Package created: $OUTPUT_TAR"
