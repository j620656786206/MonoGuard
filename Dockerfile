# Multi-stage build for Go API in monorepo
FROM golang:1.21-alpine AS builder

# Set working directory
WORKDIR /app

# Copy entire monorepo (needed for go.mod resolution)
COPY . .

# Change to API directory and build
WORKDIR /app/apps/api

# Install dependencies
RUN apk add --no-cache git ca-certificates

# Download dependencies
RUN go mod download

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates wget

# Create non-root user
RUN adduser -D -s /bin/sh appuser

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/apps/api/main .

# Change ownership to non-root user
RUN chown -R appuser:appuser /root/
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the binary
CMD ["./main"]