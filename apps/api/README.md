# MonoGuard API

Go-based API service for MonoGuard monorepo analysis platform.

## Directory Structure

```
apps/api/
├── cmd/
│   └── server/           # Application entry points
├── internal/
│   ├── config/          # Configuration management
│   ├── handlers/        # HTTP handlers
│   ├── middleware/      # HTTP middleware
│   ├── models/          # Data models
│   ├── services/        # Business logic
│   └── repository/      # Data access layer
├── pkg/
│   └── database/        # Database utilities
├── migrations/          # Database migrations
├── scripts/             # Build and deployment scripts
└── docs/                # API documentation
```

## Development

### Prerequisites

- Go 1.21+
- PostgreSQL
- Redis

### Setup

1. Copy environment variables:
```bash
cp .env.example .env
```

2. Install dependencies:
```bash
go mod tidy
```

3. Run the server:
```bash
go run cmd/server/main.go
```

### Available Endpoints

- `GET /health` - Health check
- `GET /api/v1/projects` - Projects (placeholder)
- `GET /api/v1/analysis` - Analysis (placeholder)
- `GET /api/v1/dependencies` - Dependencies (placeholder)

## Building

```bash
go build -o bin/server cmd/server/main.go
```

## Docker

```bash
docker build -t monoguard-api .
docker run -p 8080:8080 monoguard-api
```