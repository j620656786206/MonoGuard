# MonoGuard API

A comprehensive Go API for monorepo architecture analysis and validation.

## Features

- **Project Management**: CRUD operations for projects
- **Dependency Analysis**: Comprehensive dependency analysis including duplicates, version conflicts, unused packages, and circular dependencies
- **Health Scoring**: Calculate health scores based on analysis results  
- **REST API**: Full REST API with consistent response formats
- **Database Integration**: PostgreSQL with GORM ORM
- **Redis Caching**: Redis integration for caching and session management
- **Structured Logging**: JSON structured logging with request tracing
- **Health Checks**: Health, readiness, and liveness endpoints
- **Configuration Management**: Environment-based configuration with validation

## Project Structure

```
apps/api/
├── cmd/server/          # Application entrypoint
├── internal/
│   ├── app/            # Application setup and initialization
│   ├── config/         # Configuration management
│   ├── handlers/       # HTTP request handlers
│   ├── middleware/     # HTTP middleware
│   ├── models/         # Database models
│   ├── repository/     # Data access layer
│   └── services/       # Business logic layer
├── pkg/
│   └── database/       # Database connection utilities
└── tests/              # Unit and integration tests
```

## Tech Stack

- **Framework**: Gin HTTP framework
- **Database**: PostgreSQL with GORM ORM
- **Cache**: Redis
- **Logging**: Logrus with structured JSON logging
- **Configuration**: Environment variables with validation
- **Testing**: Testify testing framework

## Getting Started

### Prerequisites

- Go 1.21+
- PostgreSQL 12+
- Redis 6+

### Installation

1. **Clone the repository:**
   ```bash
   git clone <repository-url>
   cd mono-guard/apps/api
   ```

2. **Install dependencies:**
   ```bash
   go mod tidy
   ```

3. **Set up environment:**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. **Start dependencies:**
   ```bash
   # Using Docker Compose from project root
   cd ../..
   docker-compose up -d postgres redis
   ```

5. **Run the application:**
   ```bash
   go run cmd/server/main.go
   ```

## Configuration

The API uses environment variables for configuration. See `.env.example` for all available options:

### Server Configuration
- `PORT`: Server port (default: 8080)
- `HOST`: Server host (default: 0.0.0.0)  
- `GIN_MODE`: Gin mode (debug/release/test)

### Database Configuration
- `DB_HOST`: PostgreSQL host
- `DB_PORT`: PostgreSQL port
- `DB_USER`: Database user
- `DB_PASSWORD`: Database password
- `DB_NAME`: Database name
- `DB_SSL_MODE`: SSL mode (disable/require/verify-ca/verify-full)

### Redis Configuration
- `REDIS_HOST`: Redis host
- `REDIS_PORT`: Redis port
- `REDIS_PASSWORD`: Redis password (optional)
- `REDIS_DB`: Redis database number

## API Endpoints

### Health Checks
- `GET /health` - Comprehensive health check
- `GET /ready` - Readiness check
- `GET /live` - Liveness check

### Projects
- `POST /api/v1/projects` - Create project
- `GET /api/v1/projects` - List projects (paginated)
- `GET /api/v1/projects/{id}` - Get project by ID
- `PUT /api/v1/projects/{id}` - Update project
- `DELETE /api/v1/projects/{id}` - Delete project
- `POST /api/v1/projects/{id}/analyze` - Trigger dependency analysis
- `GET /api/v1/projects/owner/{ownerId}` - Get projects by owner

### Analysis
- `GET /api/v1/projects/{projectId}/analyses/dependencies` - Get dependency analyses for project
- `GET /api/v1/projects/{projectId}/analyses/dependencies/latest` - Get latest dependency analysis
- `GET /api/v1/projects/{projectId}/health-score/latest` - Get latest health score
- `GET /api/v1/analysis/dependencies/{id}` - Get dependency analysis by ID
- `GET /api/v1/analysis/architecture/{id}` - Get architecture validation by ID

## API Response Format

All API endpoints return responses in a consistent format:

### Success Response
```json
{
  "success": true,
  "data": {...},
  "message": "Operation completed successfully",
  "timestamp": "2023-12-07T10:30:00Z"
}
```

### Error Response
```json
{
  "success": false,
  "message": "Error message",
  "timestamp": "2023-12-07T10:30:00Z",
  "error": {
    "code": "ERROR_CODE",
    "message": "Detailed error message",
    "details": {...}
  }
}
```

### Paginated Response
```json
{
  "success": true,
  "data": [...],
  "message": "Data retrieved successfully",
  "timestamp": "2023-12-07T10:30:00Z",
  "pagination": {
    "currentPage": 1,
    "totalPages": 5,
    "totalItems": 50,
    "itemsPerPage": 10,
    "hasNextPage": true,
    "hasPreviousPage": false
  }
}
```

## Dependency Analysis Features

The dependency analysis engine provides:

### Duplicate Detection
- Identifies duplicate dependencies across packages
- Calculates estimated waste and bundle impact
- Provides migration recommendations

### Version Conflict Analysis
- Detects semantic version conflicts
- Identifies breaking changes
- Suggests resolution strategies

### Unused Dependency Detection
- Scans source code for actual usage
- Calculates confidence scores
- Estimates size impact

### Circular Dependency Detection
- Identifies direct and indirect circular dependencies
- Assesses severity and impact
- Provides resolution guidance

### Bundle Impact Analysis
- Calculates total bundle size
- Identifies potential savings
- Provides detailed breakdown

## Development

### Running Tests
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test package
go test ./tests/unit/services -v
```

### Building
```bash
# Build binary
go build -o bin/api cmd/server/main.go

# Build for production
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/api cmd/server/main.go
```

### Database Migrations

The application automatically runs database migrations on startup using GORM's AutoMigrate feature. This creates/updates tables based on the model definitions.

For production environments, consider using a proper migration tool like [golang-migrate](https://github.com/golang-migrate/migrate).

## Logging

The API uses structured JSON logging with request tracing:

```json
{
  "level": "info",
  "msg": "Request completed",
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "method": "GET",
  "path": "/api/v1/projects",
  "status_code": 200,
  "latency_ms": 45,
  "client_ip": "127.0.0.1",
  "time": "2023-12-07T10:30:00Z"
}
```

## Error Handling

The API implements comprehensive error handling:

- **Validation Errors**: 422 Unprocessable Entity with field details
- **Not Found**: 404 Not Found with descriptive messages  
- **Server Errors**: 500 Internal Server Error with request tracking
- **Panic Recovery**: Graceful panic recovery with logging

## Security Features

- CORS configuration for cross-origin requests
- Request ID tracking for distributed tracing
- Input validation and sanitization
- Structured error responses (no sensitive data leakage)

## Performance Considerations

- Database connection pooling
- Redis caching for expensive operations
- Efficient database queries with proper indexing
- Request timeout configuration
- Memory-efficient dependency graph processing

## Deployment

The API is designed for containerized deployment:

```bash
# Build Docker image
docker build -t monoguard-api .

# Run with Docker
docker run -p 8080:8080 --env-file .env monoguard-api
```

For production deployment, see the main project's deployment documentation.

## Contributing

1. Follow Go best practices and conventions
2. Write tests for new functionality
3. Use structured logging for important events
4. Update documentation for API changes
5. Ensure proper error handling

## License

This project is part of the MonoGuard monorepo analysis tool.