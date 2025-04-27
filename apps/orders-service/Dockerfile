# Dockerfile
FROM golang:1.24 AS builder

WORKDIR /app

# Cache Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build the app
RUN CGO_ENABLED=0 GOOS=linux go build -o orders-service ./cmd/server

# Final image
FROM alpine:3.21

RUN apk add --no-cache ca-certificates

WORKDIR /root/

COPY --from=builder /app/orders-service .

EXPOSE 8080

CMD ["./orders-service"]