#!/bin/bash

# Script untuk menjalankan Haxorport Client

set -e

# Direktori root proyek
ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

# Build terlebih dahulu
echo "Building Haxorport Client..."
./scripts/build.sh

# Jalankan aplikasi
echo "Running Haxorport Client..."
./bin/haxor "$@"
