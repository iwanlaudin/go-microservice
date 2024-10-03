![Golang Logo](https://ddev.com/img/blog/2024/05/golang-nerd-banner.png)
# GO Microservice

## Structure Directory
```
myproject/
├── services/
│   ├── order/
│   │   ├── cmd/
│   │   │   └── main.go           # Entry point untuk Order service
│   │   ├── internal/
│   │   │   ├── api/
│   │   │   │   ├── handlers/     # HTTP request handlers untuk Order
│   │   │   │   ├── middleware/   # HTTP middleware khusus Order
│   │   │   │   └── routes/       # Order API route definitions
│   │   │   ├── models/           # Model data Order
│   │   │   ├── repository/       # Lapisan akses data Order
│   │   │   └── service/          # Logika bisnis Order
│   │   ├── config/               # Konfigurasi Order service
│   │   └── Dockerfile
│   └── ...                       # Folder untuk service lainnya
├── pkg/
│   ├── common/                   # Paket umum yang digunakan oleh semua services
│   │   ├── logger/
│   │   ├── config/               # Konfigurasi semua service
│   │   ├── database/
│   │   └── auth/
│   ├── rabbitmq/                 # Utilitas RabbitMQ
│   └── ...                       # Paket shared lainnya
├── scripts/
│   ├── build.sh                  # Script untuk build semua services
│   └── deploy.sh                 # Script untuk deploy semua services
├── deployments/
│   ├── docker-compose.yml        # Untuk menjalankan semua services secara lokal
│   └── kubernetes/               # Konfigurasi Kubernetes jika diperlukan
├── migrations/                   # Database migrations
├── test/
│   └── integration/              # Tes integrasi antar services
├── go.mod
├── go.sum
└── README.md
```