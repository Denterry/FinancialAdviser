# X Service

A microservice for handling Twitter posts, built with Go and following clean architecture principles.

## Features

- Tweet management (create, read, update, delete)
- gRPC API
- PostgreSQL database
- Docker support
- Clean architecture

## Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose
- PostgreSQL 15 or higher
- Make

## Getting Started

1. Clone the repository:
```bash
git clone https://github.com/Denterry/FinancialAdviser.git
cd Backend/x-service
```

2. Install dependencies:
```bash
go mod download
```

3. Set up the database:
```bash
make migrate-up
```

4. Run the service:
```bash
make run
```

Or using Docker:
```bash
make docker-build
make docker-run
```

## Project Structure

```
.
├── cmd/
│   └── app/
│       └── main.go
├── config/
│   └── config.yaml
├── internal/
│   ├── app/
│   │   └── app.go
│   ├── controller/
│   │   └── grpc/
│   │       └── auth.go
│   ├── entity/
│   │   └── user.go
│   ├── repo/
│   │   ├── contracts.go
│   │   └── persistent/
│   │       └── user_postgres.go
│   └── usecase/
│       └── auth/
│           ├── auth.go
│           └── token.go
├── migrations/
│   └── 000001_init.up.sql
├── pkg/
│   ├── config/
│   │   └── config.go
│   ├── grpcserver/
│   │   └── server.go
│   ├── logger/
│   │   └── logger.go
│   └── postgres/
│       └── postgres.go
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## API Documentation

The service exposes a gRPC API on port 50051. You can find the protocol buffer definitions in the `internal/controller/grpc/proto` directory.

## Development

### Running Tests

```bash
make test
```

### Linting

```bash
make lint
```

### Database Migrations

To create a new migration:
```bash
migrate create -ext sql -dir migrations -seq migration_name
```

To apply migrations:
```bash
make migrate-up
```

To rollback migrations:
```bash
make migrate-down
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.
