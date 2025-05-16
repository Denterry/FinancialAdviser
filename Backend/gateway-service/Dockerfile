# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git make

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o gateway-service ./cmd/app

# Final stage
FROM alpine:3.18

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Copy the binary and config files
COPY --from=builder /app/gateway-service .
COPY --from=builder /app/config/config.yaml ./config/

# Expose port
EXPOSE 8080

# Run the application
CMD ["./gateway-service"] 