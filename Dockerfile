# Build stage
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Install git and build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set the GOPROXY environment variable
ENV GOPROXY=https://goproxy.io,direct
# Set environment variable allow bypassing the proxy for specified repos (optional)
ENV GOPRIVATE=git.mycompany.com,github.com/my/private

# Install swag for generating swagger docs
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copy go mod files first for better layer caching
COPY go.mod go.sum ./

COPY .env /app/

# Download dependencies
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Generate swagger documentation
RUN swag init -g cmd/main.go -o docs

# Build the application with security flags
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o main ./cmd/main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Copy migrations directory
COPY --from=builder /app/internal/database/migrations ./internal/database/migrations

# Copy .env for runtime configuration (optional if using env vars or docker-compose env_file)
COPY --from=builder /app/.env ./

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"]
