![Golang Logo](https://ddev.com/img/blog/2024/05/golang-nerd-banner.png)
# GO Microservice

## Structure Directory
```
myproject/
├── services/
│   ├── tickets/
│   │   ├── cmd/
│   │   │   └── main.go           # Entry point untuk Ticket service
│   │   ├── internal/
│   │   │   ├── api/
│   │   │   │   ├── handlers/     # HTTP request handlers untuk Ticket
│   │   │   │   ├── middleware/   # HTTP middleware khusus Ticket
│   │   │   │   └── routes/       # Ticket API route definitions
│   │   │   ├── models/           # Model data Ticket
│   │   │   ├── repository/       # Lapisan akses data Ticket
│   │   │   └── service/          # Logika bisnis Ticket
│   │   │   └── messaging/        # Logika messaging Ticket
│   │   │   └── dto/              # Folder baru untuk request dan response
│   │   │       ├── request/      # Definisi struktur request
│   │   │       └── response/     # Definisi struktur response
│   │   ├── config/               # Konfigurasi Ticket service
│   │   ├── migrations/           # Database migrations Ticket service
│   │   └── Dockerfile
│   └── ...                       # Folder untuk service lainnya
├── pkg/
│   ├── common/                   # Paket umum yang digunakan oleh semua services
│   │   ├── logger/
│   │   ├── config/               # Konfigurasi semua service
│   │   ├── database/
│   │   └── auth/
│   ├── rabbitmq/                 # Utilitas RabbitMQ
│   └── redis/                    # Utilitas Redis
│   └── email/                    # Utilitas SMTP
│   └── ...                       # Paket shared lainnya
├── scripts/
│   ├── build.sh                  # Script untuk build semua services
│   └── deploy.sh                 # Script untuk deploy semua services
├── deployments/
│   ├── docker-compose.yml        # Untuk menjalankan semua services secara lokal
│   └── kubernetes/               # Konfigurasi Kubernetes jika diperlukan
├── test/
│   └── integration/              # Tes integrasi antar services
├── go.mod
├── go.sum
└── README.md
```

## Database and Interaction Diagram
```
+-----------------------+              +---------------------+
|   Database API Event  |              |  Database API Tiket |
|-----------------------|              |---------------------|
|  events               |<-------------|  tickets            |
+-----------------------+              +---------------------+
       ^                                      |
       |                                      v
+-----------------------+           +---------------------+
|Database API Notifikasi|          +-----------------------+
|-----------------------|          |Database API Pembayaran|
|  notifications        |<---------|-----------------------|
+-----------------------+          |  payments             |
                                   +-----------------------+
                                           ^
                                           |
                               +--------------------------+
                               | Database API Autentikasi |
                               |--------------------------|
                               | users                    |
                               | auth_tokens              |
                               +--------------------------+

```

## Features

- **RESTful API**: Manage data with easily accessible CRUD endpoints.
- **Rate Limiting**: Controls the number of requests a client can make in a given time period to prevent abuse and ensure quality of service.
- **API Versioning**: Provides multiple versions of the API to allow changes without disrupting existing users.
- **Authentication and Authorization**: Uses JWT tokens for API security.
- **Queue Management**: Integration with RabbitMQ for asynchronous processing.
- **Configuration Management**: Easily customizable via environment variables.
- **Redis Caching**: Caches frequently accessed data using Redis to improve response times and reduce database load.
- **Email Notifications**: Sends email notifications to users for various events, such as account creation, password resets, and alerts.
- **SMTP Configuration**: Easily configurable SMTP settings to connect to your email server.
- **Graceful Shutdown**: Application can gracefully handle shutdown signals to ensure ongoing processes complete before termination.
- **Logging**: Structured logging with Zap for easier debugging and tracing.
- **Scalability**: Designed to accommodate user growth.
- **Data Storage**: Support for postgres database types.
- **Error Management**: Provides clear error response.

## How to Run an API Application
1. Clone Repository
```bash
git clone https://github.com/iwanlaudin/go-microservice.git
```
2. Install Dependencies
```base
go mod tidy
```
3. Run Application
```base
go run services/authentication/cmd/main.go
go run services/order/cmd/main.go
...
```
## Running Go Migrate
1. Make sure the migration is available in the migrations folder. This folder is usually located in `migrations`
2. Creating a New Migration File

To create a new migration file, you can use the following command:
```bash
migrate create -ext sql -dir services/authentication/migrations -seq nama_migrasi
```
3. Run the migration command To apply the migration, run the following command:
```base
migrate -database "postgres://postgres:root@localhost:5432/AuthDb?sslmode=disable" -path services/authentication/migrations up
```
To cancel the migration, use:
```base
migrate -database "postgres://postgres:root@localhost:5432/AuthDb?sslmode=disable" -path services/authentication/migrations down
```

## Testing the API
```base
curl -X GET http://localhost:8080/endpoint
```

# Contact
If you have any questions or feedback, please contact iwanlaudin01@gmail.com.

Please customize with your app details, such as app name, description, available endpoints, and license. If there are any other sections you would like to add, please let me know!