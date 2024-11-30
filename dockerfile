# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache \
    git \
    gcc \
    musl-dev

# Copy go mod files first for better layer caching
COPY go.mod go.sum ./

# Download dependencies and tidy
RUN go mod download && \
    go mod tidy

# Copy source code
COPY . .

# Run tidy again after code copy
RUN go mod tidy

# Build the application with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o realpay .

# Final stage
FROM alpine:3.19

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    tzdata

# Create non-root user
RUN adduser -D appuser

# Create and set permissions for logs directory
RUN mkdir -p /app/logs && \
    chown -R appuser:appuser /app/logs

# Copy binary and migrations from builder
COPY --from=builder /app/realpay .
COPY --from=builder /app/migrations ./migrations

# Set ownership for application files
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Set environment variables
ENV PORT=8082 \
    GIN_MODE=release

# Expose port
EXPOSE 8082

# Run the application
CMD ["./realpay"]
