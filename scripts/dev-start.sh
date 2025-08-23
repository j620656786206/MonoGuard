#!/bin/bash

# MonoGuard Development Startup Script

echo "🚀 Starting MonoGuard Development Environment..."

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker first."
    exit 1
fi

# Check if .env file exists
if [ ! -f .env ]; then
    echo "📋 Creating .env file from template..."
    cp .env.example .env
    echo "✅ .env file created. Please update it with your configuration."
fi

# Start services with Docker Compose
echo "🐳 Starting Docker services..."
docker-compose up -d postgres redis

# Wait for services to be healthy
echo "⏳ Waiting for database to be ready..."
until docker-compose exec postgres pg_isready -U monoguard > /dev/null 2>&1; do
  sleep 1
done

echo "⏳ Waiting for Redis to be ready..."
until docker-compose exec redis redis-cli ping > /dev/null 2>&1; do
  sleep 1
done

echo "✅ Infrastructure services are ready!"

# Install dependencies if needed
if [ ! -d "node_modules" ]; then
    echo "📦 Installing dependencies..."
    pnpm install
fi

# Build shared types
echo "🔧 Building shared types..."
pnpm nx build shared-types

echo "🎉 Development environment is ready!"
echo ""
echo "Available services:"
echo "  🗄️  PostgreSQL:  localhost:5432"
echo "  🔴 Redis:       localhost:6379"
echo "  🌐 Adminer:     http://localhost:8081"
echo ""
echo "To start the applications:"
echo "  📱 Frontend:    pnpm dev:frontend"
echo "  🚀 API:         go run apps/api/cmd/server/main.go"
echo "  ⚡ CLI:         pnpm dev:cli"