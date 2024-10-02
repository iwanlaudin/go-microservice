#!/bin/bash

# File: scripts/run-all-services.sh

# Fungsi untuk menjalankan service
run_service() {
    cd $1
    go run cmd/main.go &
    cd -
}

# Masuk ke direktori root proyek
cd "$(dirname "$0")/.."

# Jalankan setiap service
run_service "services/order"
# run_service "services/payment"
# Tambahkan service lain di sini dengan format yang sama

echo "All services are running. Press Ctrl+C to stop all services."

# Tunggu input user untuk menghentikan semua service
trap "trap - SIGTERM && kill -- -$$" SIGINT SIGTERM EXIT

# Tunggu semua background process
wait