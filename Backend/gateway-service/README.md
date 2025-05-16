# Gateway Service

The Gateway Service is a crucial component of the Financial Adviser platform that acts as a single entry point for all client requests. It handles routing, authentication, and request/response transformation between the client and various microservices.

## Features

- **API Gateway**: Routes requests to appropriate microservices
- **Authentication**: Handles user authentication and authorization
- **Request Transformation**: Transforms HTTP requests to gRPC calls
- **Response Aggregation**: Combines responses from multiple services
- **Rate Limiting**: Prevents abuse through request rate limiting
- **CORS Handling**: Manages Cross-Origin Resource Sharing
- **Logging**: Comprehensive request/response logging
- **Circuit Breaking**: Prevents cascading failures
- **Load Balancing**: Distributes traffic across service instances

## Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose
- Redis
- Access to other microservices (Auth, Subscription, ML)

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/Denterry/FinancialAdviser.git
   cd Backend/gateway-service
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Set up environment variables:
   ```bash
   cp config/config.yaml.example config/config.yaml
   # Edit config/config.yaml with your settings
   ```

4. Build and run:
   ```bash
   make build
   make run
   ```

## Docker Deployment

1. Build the Docker image:
   ```bash
   make docker-build
   ```

2. Run with Docker Compose:
   ```bash
   make docker-run
   ```

## API Documentation

### Authentication Endpoints

- `POST /api/auth/register` - Register a new user
- `POST /api/auth/login` - Login user
- `POST /api/auth/refresh` - Refresh access token

### Subscription Endpoints

- `GET /api/subscription/plans` - Get available subscription plans
- `POST /api/subscription/subscribe` - Subscribe to a plan
- `GET /api/subscription/status` - Get subscription status
- `POST /api/subscription/cancel` - Cancel subscription

### ML Endpoints

- `POST /api/ml/analyze` - Analyze financial data
- `GET /api/ml/recommendations` - Get financial recommendations

## Configuration

The service can be configured using environment variables or the `config.yaml` file:

```yaml
server:
  port: ":8080"
  read_timeout: 5s
  write_timeout: 10s

services:
  auth:
    host: "auth-service"
    port: ":50051"
  subscription:
    host: "sub-service"
    port: ":50052"
  ml:
    host: "ml-service"
    port: ":50053"

redis:
  host: "redis"
  port: ":6379"
  password: ""
  db: 0

jwt:
  secret: "your-secret-key"
  access_token_ttl: 15m
  refresh_token_ttl: 720h
  signing_method: "HS256"
  signing_key: "your-signing-key"
  verification_key: "your-verification-key"

log:
  level: "info"
  file: "logs/gateway.log"
```

## Project Structure

```
.
├── cmd/
│   └── app/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── handler/
│   │   └── handler.go
│   ├── middleware/
│   │   └── middleware.go
│   └── service/
│       └── service.go
├── config/
│   └── config.yaml
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## Development

### Running Tests

```bash
make test
```

### Code Generation

```bash
make generate
```

### Installing Development Tools

```bash
make install-tools
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details. 