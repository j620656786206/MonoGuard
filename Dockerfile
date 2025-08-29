# Build stage
FROM golang:1.23-alpine AS builder

# Install dependencies
RUN apk add --no-cache git ca-certificates

# Set working directory
WORKDIR /app

# Copy the entire project
COPY . .

# Go to API directory and build
WORKDIR /app/apps/api
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates wget

# Create non-root user
RUN adduser -D -s /bin/sh appuser

# Set working directory
WORKDIR /home/appuser

# Copy binary from build stage
COPY --from=builder /app/apps/api/main .

# Set ownership
RUN chown appuser:appuser /home/appuser/main
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
CMD ["./main"]