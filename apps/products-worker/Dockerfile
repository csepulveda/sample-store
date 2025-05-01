# Dockerfile
FROM golang:1.24 AS builder

WORKDIR /app

# Cache Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build the app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} go build -o products-worker ./cmd/worker

# Final image
FROM alpine:3.21

RUN apk add --no-cache ca-certificates

WORKDIR /root/

COPY --from=builder /app/products-worker .

CMD ["./products-worker"]