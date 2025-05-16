# Subscription Service

The Subscription Service is a microservice responsible for managing user subscriptions, plans, and payments in the Financial Adviser platform. It provides functionality for creating and managing subscription plans, handling user subscriptions, processing payments, and managing subscription lifecycle.

## Features

- **Subscription Plan Management**
  - Create, read, update, and delete subscription plans
  - Define plan features, pricing, and duration
  - Support for different plan types (Basic, Pro, Enterprise)

- **User Subscription Management**
  - Create new subscriptions
  - View subscription details
  - Update subscription settings
  - Cancel subscriptions
  - Handle subscription renewals

- **Payment Processing**
  - Integration with Stripe for payment processing
  - Support for multiple payment methods
  - Handle payment webhooks
  - Process subscription renewals
  - Manage payment history

- **Subscription Lifecycle**
  - Handle subscription status changes
  - Manage grace periods
  - Process subscription expirations
  - Handle auto-renewal settings

## Prerequisites

- Go 1.21 or higher
- PostgreSQL 15 or higher
- Docker and Docker Compose (for containerized deployment)
- Stripe account (for payment processing)

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/FinancialAdviser.git
   cd FinancialAdviser/Backend/sub-service
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Set up environment variables:
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. Run database migrations:
   ```bash
   make migrate-up
   ```

## Running the Service

### Local Development

1. Start the service:
   ```bash
   go run cmd/app/main.go
   ```

2. Or use Docker Compose:
   ```bash
   docker-compose up -d
   ```

### Production Deployment

1. Build the Docker image:
   ```bash
   docker build -t subscription-service .
   ```

2. Run the container:
   ```bash
   docker run -d --name subscription-service \
     -p 8080:8080 \
     --env-file .env \
     subscription-service
   ```

## API Documentation

The service provides the following gRPC endpoints:

### Plan Management
- `CreatePlan`: Create a new subscription plan
- `GetPlan`: Retrieve plan details
- `ListPlans`: List available plans
- `UpdatePlan`: Update plan details
- `DeletePlan`: Delete a plan

### Subscription Management
- `CreateSubscription`: Create a new subscription
- `GetSubscription`: Get subscription details
- `ListSubscriptions`: List user subscriptions
- `UpdateSubscription`: Update subscription settings
- `CancelSubscription`: Cancel a subscription

### Payment Processing
- `ProcessPayment`: Process a subscription payment
- `GetPayment`: Get payment details
- `ListPayments`: List subscription payments

## Configuration

The service can be configured using environment variables:

- `SERVER_PORT`: Port to run the service on (default: 8080)
- `SERVER_HOST`: Host to bind the service to (default: localhost)
- `DB_HOST`: PostgreSQL host
- `DB_PORT`: PostgreSQL port
- `DB_USER`: Database user
- `DB_PASSWORD`: Database password
- `DB_NAME`: Database name
- `JWT_SECRET`: Secret key for JWT tokens
- `STRIPE_SECRET_KEY`: Stripe API secret key
- `STRIPE_WEBHOOK_SECRET`: Stripe webhook secret
- And more...

See `.env.example` for all available configuration options.

## Development

### Project Structure

```
.
├── cmd/
│   └── app/
│       └── main.go
├── config/
│   └── config.go
├── internal/
│   ├── controller/
│   │   └── grpc/
│   ├── entity/
│   ├── repo/
│   │   └── postgres/
│   └── usecase/
├── migrations/
│   ├── 000001_init_schema.up.sql
│   └── 000001_init_schema.down.sql
├── pkg/
│   └── pb/
│       └── subscription/
│           └── v1/
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
└── README.md
```

### Running Tests

```bash
make test
```

### Code Generation

Generate gRPC code:
```bash
make generate
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
