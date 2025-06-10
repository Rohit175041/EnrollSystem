# -------- Stage 1: Build the Go binary --------
FROM golang:1.23-alpine AS builder

# Install Git and timezone (optional), and CA certificates for Go modules
RUN apk add --no-cache git tzdata ca-certificates

WORKDIR /app

# Only copy go.mod and go.sum first (improves cache usage)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go binary with stripped debug info
RUN go build -ldflags="-s -w" -o enrollsystem ./cmd/students-api/main.go


# -------- Stage 2: Minimal runtime image --------
FROM alpine:latest

# Copy only what's needed
RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copy binary and config
COPY --from=builder /app/enrollsystem .
COPY config ./config

# Use a non-root user for security
RUN adduser -D appuser
USER appuser

EXPOSE 8080

# Start the app
CMD ["./enrollsystem", "-config", "./config/local.yaml"]
