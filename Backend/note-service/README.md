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

## API Documentation

The service exposes a gRPC API on port 9090. You can find the protocol buffer definitions in the `internal/controller/grpc/proto` directory

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
.
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
│   │       ├── admin.go
│   │       ├── proto
│   │       │   ├── admin
│   │       │   │   └── v1
│   │       │   │       └── admin.proto
│   │       │   └── tweets
│   │       │       └── v1
│   │       │           └── tweets.proto
│   │       ├── tweets.go
│   │       └── types.go
│   ├── entity
│   │   ├── errors.go
│   │   ├── tweet.go
│   │   └── types.go
│   ├── repo
│   │   ├── contracts.go
│   │   ├── persistent
│   │   │   ├── helpers.go
│   │   │   └── tweet_postgres.go
│   │   └── webapi
│   │       ├── factory.go
│   │       ├── helpers.go
│   │       ├── twitter_api.go
│   │       └── twitter_scraper.go
│   └── usecase
│       ├── admin
│       │   └── tweet.go
│       ├── contracts.go
│       └── tweet
│           └── tweet.go
├── Makefile
├── migrations
│   ├── 20250503231352_tweets_table.down.sql
│   ├── 20250503231352_tweets_table.up.sql
│   ├── 20250515112831_uuid_ext.down.sql
│   ├── 20250515112831_uuid_ext.up.sql
│   ├── 20250515130212_authors_table.down.sql
│   ├── 20250515130212_authors_table.up.sql
│   ├── 20250515130456_symbols_table.down.sql
│   ├── 20250515130456_symbols_table.up.sql
│   ├── 20250515131510_tweet_symbols_table.down.sql
│   ├── 20250515131510_tweet_symbols_table.up.sql
│   ├── 20250515132712_articles_table.down.sql
│   ├── 20250515132712_articles_table.up.sql
│   ├── 20250515133338_article_symbol_table.down.sql
│   ├── 20250515133338_article_symbol_table.up.sql
│   ├── 20250515133514_crawl_jobs_table.down.sql
│   ├── 20250515133514_crawl_jobs_table.up.sql
│   ├── 20250515133756_crawl_job_logs.down.sql
│   ├── 20250515133756_crawl_job_logs.up.sql
│   ├── 20250515133947_provider_failures_table.down.sql
│   ├── 20250515133947_provider_failures_table.up.sql
│   ├── 20250515134210_updated_at_trigger.down.sql
│   ├── 20250515134210_updated_at_trigger.up.sql
│   ├── 20250515134349_indexes.down.sql
│   ├── 20250515134349_indexes.up.sql
│   ├── 20250515134520_tweet_search_mv.down.sql
│   ├── 20250515134520_tweet_search_mv.up.sql
│   ├── 20250515134741_sentiment_daily_agg_table.down.sql
│   └── 20250515134741_sentiment_daily_agg_table.up.sql
├── pkg
│   ├── grpcserver
│   │   ├── options.go
│   │   └── server.go
│   ├── logger
│   │   └── logger.go
│   ├── pb
│   │   ├── admin
│   │   │   └── v1
│   │   │       ├── admin_grpc.pb.go
│   │   │       └── admin.pb.go
│   │   └── tweets
│   │       └── v1
│   │           ├── tweets_grpc.pb.go
│   │           └── tweets.pb.go
│   └── postgres
│       ├── options.go
│       └── postgres.go
└── README.md
```