#!/bin/bash

# MonoGuard Production Deployment Script

set -e  # Exit on any error

echo "🚀 Starting MonoGuard Production Deployment..."

# Check required environment variables
required_vars=("DB_PASSWORD" "JWT_SECRET" "NEXTAUTH_SECRET")
for var in "${required_vars[@]}"; do
    if [ -z "${!var}" ]; then
        echo "❌ Environment variable $var is required but not set"
        exit 1
    fi
done

# Build production images
echo "🔨 Building production Docker images..."
docker-compose -f docker-compose.prod.yml build --no-cache

# Run database migrations (if needed)
echo "🗄️  Running database migrations..."
# Add migration commands here when implemented

# Start production services
echo "🐳 Starting production services..."
docker-compose -f docker-compose.prod.yml up -d

# Wait for services to be healthy
echo "⏳ Waiting for services to start..."
sleep 30

# Health checks
echo "🏥 Running health checks..."
if ! curl -f http://localhost:8080/health > /dev/null 2>&1; then
    echo "❌ API health check failed"
    exit 1
fi

if ! curl -f http://localhost:3000 > /dev/null 2>&1; then
    echo "❌ Frontend health check failed"
    exit 1
fi

echo "✅ All services are healthy!"
echo ""
echo "🎉 MonoGuard is now running in production mode!"
echo ""
echo "Available endpoints:"
echo "  📱 Frontend: http://localhost:3000"
echo "  🚀 API:      http://localhost:8080"
echo ""
echo "To stop: docker-compose -f docker-compose.prod.yml down"
echo "To view logs: docker-compose -f docker-compose.prod.yml logs -f"