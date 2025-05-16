# Step 1: Modules caching
FROM golang:1.24.2-alpine3.21 as modules

WORKDIR /modules
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Step 2: Builder
FROM golang:1.24.2-alpine3.21 as builder

# fetch CA cert bundle so our HTTP/gRPC clients can validate TLS
RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=modules /go/pkg /go/pkg
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -tags migrate -ldflags="-s -w" -o ./app ./cmd/app

# Step 3: Final
FROM scratch

WORKDIR /

COPY --from=builder /app/config /config
COPY --from=builder /app/migrations /migrations
COPY --from=builder /app/app /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 50051
EXPOSE 9090

CMD ["/app"]
