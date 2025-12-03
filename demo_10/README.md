# Book Service - Go CRUD Application

A simple CRUD application for managing books built with Go, PostgreSQL, and Prometheus monitoring.

## Features

- **CRUD Operations**: Create, Read, Update, Delete books
- **PostgreSQL**: Persistent data storage
- **Prometheus**: Metrics collection and monitoring
- **Docker Compose**: Easy setup with containerization

## Prerequisites

- Docker and Docker Compose
- Go 1.21+ (for local development)

## Quick Start

### Using Docker Compose

1. Start the services:
```bash
docker-compose up -d
```

This will start:
- PostgreSQL on port 5432
- Prometheus on port 9090

2. Build and run the Go application:
```bash
go mod download
go run ./cmd/main.go
```

The application will be available at `http://localhost:8080`

### API Endpoints

#### Get all books
```bash
curl http://localhost:8080/api/books
```

#### Create a book
```bash
curl -X POST http://localhost:8080/api/books \
  -H "Content-Type: application/json" \
  -d '{
    "title": "The Go Programming Language",
    "author": "Alan Donovan and Brian Kernighan",
    "isbn": "978-0134190440",
    "pages": 400,
    "published": "2015-10-26T00:00:00Z"
  }'
```

#### Get a specific book
```bash
curl http://localhost:8080/api/books/1
```

#### Update a book
```bash
curl -X PUT http://localhost:8080/api/books/1 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Updated Title",
    "pages": 450
  }'
```

#### Delete a book
```bash
curl -X DELETE http://localhost:8080/api/books/1
```

#### Health check
```bash
curl http://localhost:8080/health
```

### Monitoring

View Prometheus metrics:
- http://localhost:9090 (Prometheus UI)
- http://localhost:8080/metrics (Application metrics)

### Environment Variables

- `PORT`: Server port (default: :8080)
- `DB_USER`: Database user (default: bookuser)
- `DB_PASSWORD`: Database password (default: bookpass)
- `DB_HOST`: Database host (default: localhost)
- `DB_PORT`: Database port (default: 5432)
- `DB_NAME`: Database name (default: bookdb)

## Project Structure

```
.
├── cmd/
│   └── main.go                 # Entry point
├── internal/
│   ├── handler/
│   │   └── book_handler.go     # HTTP handlers
│   ├── models/
│   │   └── book.go             # Data models
│   ├── repository/
│   │   └── book_repository.go  # Database operations
│   └── service/
│       └── book_service.go     # Business logic
├── pkg/
│   └── database/
│       └── db.go               # Database connection & setup
├── docker-compose.yml
├── Dockerfile
├── prometheus.yml
└── go.mod

```

## Development

### Local Setup

1. Install dependencies:
```bash
go mod download
```

2. Make sure PostgreSQL is running via Docker Compose

3. Run the application:
```bash
go run ./cmd/main.go
```

### Building

Build the application:
```bash
go build -o book-app ./cmd/main.go
```

Build Docker image:
```bash
docker build -t book-service .
```

## Cleaning Up

Stop and remove Docker containers:
```bash
docker-compose down
```

To also remove volumes (database data):
```bash
docker-compose down -v
```
