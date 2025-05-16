# Auth Service

A microservice for handling user authentication and authorization, built with Go and following clean architecture principles.

## Features

- User authentication and authorization
- JWT token management
- Role-based access control (RBAC)
- gRPC API
- PostgreSQL database
- Docker support
- Clean architecture

## Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose
- PostgreSQL 15 or higher

## Getting Started

1. Clone the repository:
```bash
git clone https://github.com/Denterry/FinancialAdviser.git
cd Backend/auth-service
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

## API Documentation

The service exposes a gRPC API on port 50051. You can find the protocol buffer definitions in the `internal/controller/grpc/proto` directory

## Development

### Generate Protocol Buffers
```bash
make proto-all
```

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

## Project Structure

```
├── cmd
│   └── app
│       └── main.go
├── config
│   └── config.go
├── docker-compose-integration-test.yml
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
├── integration-tests
│   ├── Dockerfile
│   └── integration_test.go
├── internal
│   ├── app
│   │   ├── app.go
│   │   └── migrate.go
│   ├── controller
│   │   └── grpc
│   │       ├── auth.go
│   │       └── proto
│   │           └── auth
│   │               └── v1
│   │                   └── auth_service.proto
│   ├── entity
│   │   └── user.go
│   ├── mapper
│   │   └── user.go
│   ├── repo
│   │   ├── contracts.go
│   │   └── persistent
│   │       └── user_postgres.go
│   └── usecase
│       ├── auth
│       │   ├── auth.go
│       │   ├── errors.go
│       │   └── token.go
│       ├── contracts.go
│       └── errors.go
├── Makefile
├── migrations
│   ├── 20250420095249_users_table.down.sql
│   └── 20250420095249_users_table.up.sql
├── pkg
│   ├── grpcserver
│   │   ├── options.go
│   │   └── server.go
│   ├── httpserver
│   │   ├── options.go
│   │   └── server.go
│   ├── logger
│   │   └── logger.go
│   ├── pb
│   │   └── auth
│   │       └── v1
│   │           ├── auth_service_grpc.pb.go
│   │           └── auth_service.pb.go
│   └── postgres
│       ├── options.go
│       └── postgres.go
└── README.md
```
