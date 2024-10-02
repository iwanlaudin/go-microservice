#!/bin/bash
# File: scripts/deploy.sh

set -e  # Exit immediately if a command exits with a non-zero status.

# Masuk ke direktori root proyek
cd "$(dirname "$0")/.."

# Fungsi untuk deploy service
deploy_service() {
    echo "Deploying $1 service..."
    scp "bin/$1" user@your-server.com:/path/to/services/$1/
    ssh user@your-server.com "systemctl restart $1-service"
}

# Deploy setiap service
deploy_service "order"
deploy_service "payment"
# Tambahkan service lain di sini

echo "All services have been deployed successfully."