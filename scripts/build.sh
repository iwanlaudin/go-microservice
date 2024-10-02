#!/bin/bash
# File: scripts/build.sh

set -e  # Exit immediately if a command exits with a non-zero status.

# Masuk ke direktori root proyek
cd "$(dirname "$0")/.."

# Fungsi untuk membangun service
build_service() {
    echo "Building $1 service..."
    cd "services/$1"
    go build -o "../../bin/$1" cmd/main.go
    cd ../..
}

# Buat direktori bin jika belum ada
mkdir -p bin

# Build setiap service
build_service "order"
build_service "payment"
# Tambahkan service lain di sini

echo "All services have been built successfully."