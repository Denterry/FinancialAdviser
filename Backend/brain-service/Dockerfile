FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o x-service ./cmd/app

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/x-service .
COPY --from=builder /app/config/config.yaml ./config/

EXPOSE 50051

CMD ["./x-service"]
