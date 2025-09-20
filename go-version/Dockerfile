# Build stage
FROM golang:1.25.1 AS builder

WORKDIR /app

# Set GOBIN for tools
ENV GOBIN=/usr/local/bin

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Install migration tool
RUN go install -tags 'sqlite' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

# Production stage
FROM debian:bookworm-slim

# Install SQLite (if needed for runtime)
RUN apt-get update && apt-get install -y \
    sqlite3 \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy built binary from builder stage
COPY --from=builder /app/main .

# Copy migration tool from builder stage
COPY --from=builder /usr/local/bin/migrate /usr/local/bin/migrate

# Copy migrations
COPY --from=builder /app/migrations ./migrations

# Create database directory
RUN mkdir -p /app/data

# Expose port
EXPOSE 8080

# Default command
CMD ["./main"]