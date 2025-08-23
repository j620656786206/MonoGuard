# MonoGuard Deployment Guide

## Local Development Setup

### Prerequisites
- Go 1.21+
- Node.js 18+
- PostgreSQL 15+
- Redis 6+
- Docker & Docker Compose

### Quick Setup
```bash
# Clone and setup repository
git clone <repo>
cd mono-guard

# Start dependencies with Docker Compose
docker-compose up -d postgres redis

# Setup backend
cd backend
go mod download
go run cmd/api/main.go

# Setup frontend (in new terminal)
cd frontend
npm install
npm run dev

# Setup CLI (in new terminal)
cd cli
npm install
npm run build
```

### Development Environment Configuration

#### Docker Compose Development
```yaml
# docker-compose.dev.yml
version: '3.8'
services:
  postgres:
    image: postgres:15
    environment:
      POSTGRES_DB: monoguard_dev
      POSTGRES_USER: dev
      POSTGRES_PASSWORD: dev123
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./backend/migrations:/docker-entrypoint-initdb.d/

  redis:
    image: redis:6-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data

  backend:
    build: 
      context: ./backend
      dockerfile: Dockerfile.dev
    environment:
      DATABASE_URL: postgres://dev:dev123@postgres:5432/monoguard_dev
      REDIS_URL: redis://redis:6379
      JWT_SECRET: dev-secret-key
      LOG_LEVEL: debug
    ports:
      - "8080:8080"
    volumes:
      - ./backend:/app
    depends_on:
      - postgres
      - redis

  frontend:
    build:
      context: ./frontend
      dockerfile: Dockerfile.dev
    environment:
      NEXT_PUBLIC_API_URL: http://localhost:8080
      NODE_ENV: development
    ports:
      - "3000:3000"
    volumes:
      - ./frontend:/app
      - /app/node_modules
    depends_on:
      - backend

volumes:
  postgres_data:
  redis_data:
```

#### Environment Variables
```bash
# .env.development
DATABASE_URL=postgres://dev:dev123@localhost:5432/monoguard_dev
REDIS_URL=redis://localhost:6379
JWT_SECRET=dev-secret-key-change-in-production
LOG_LEVEL=debug

# Frontend (.env.local)
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXT_PUBLIC_ENVIRONMENT=development
```

## Production Deployment Strategy

### Container Strategy

#### Multi-stage Docker Builds
```dockerfile
# backend/Dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/api/main.go

FROM alpine:3.18
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations/
EXPOSE 8080
CMD ["./main"]
```

```dockerfile
# frontend/Dockerfile
FROM node:18-alpine AS deps
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci --only=production && npm cache clean --force

FROM node:18-alpine AS builder
WORKDIR /app
COPY . .
COPY --from=deps /app/node_modules ./node_modules
RUN npm run build

FROM node:18-alpine AS runner
WORKDIR /app
ENV NODE_ENV=production
RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs
COPY --from=builder /app/next.config.js ./
COPY --from=builder /app/public ./public
COPY --from=builder /app/.next/standalone ./
COPY --from=builder /app/.next/static ./.next/static
USER nextjs
EXPOSE 3000
ENV PORT=3000
CMD ["node", "server.js"]
```

### Kubernetes Deployment

#### Namespace and ConfigMap
```yaml
# k8s/namespace.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: monoguard
  labels:
    name: monoguard

---
apiVersion: v1
kind: ConfigMap
metadata:
  name: monoguard-config
  namespace: monoguard
data:
  LOG_LEVEL: "info"
  REDIS_URL: "redis://redis:6379"
  DATABASE_MAX_CONNECTIONS: "20"
```

#### Database Deployment
```yaml
# k8s/postgres.yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: postgres
  namespace: monoguard
spec:
  serviceName: postgres
  replicas: 1
  selector:
    matchLabels:
      app: postgres
  template:
    metadata:
      labels:
        app: postgres
    spec:
      containers:
      - name: postgres
        image: postgres:15
        env:
        - name: POSTGRES_DB
          value: "monoguard"
        - name: POSTGRES_USER
          valueFrom:
            secretKeyRef:
              name: postgres-secret
              key: username
        - name: POSTGRES_PASSWORD
          valueFrom:
            secretKeyRef:
              name: postgres-secret
              key: password
        ports:
        - containerPort: 5432
        volumeMounts:
        - name: postgres-storage
          mountPath: /var/lib/postgresql/data
        resources:
          requests:
            memory: "512Mi"
            cpu: "250m"
          limits:
            memory: "1Gi"
            cpu: "500m"
  volumeClaimTemplates:
  - metadata:
      name: postgres-storage
    spec:
      accessModes: ["ReadWriteOnce"]
      resources:
        requests:
          storage: 10Gi
```

#### API Deployment
```yaml
# k8s/api-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: monoguard-api
  namespace: monoguard
spec:
  replicas: 3
  selector:
    matchLabels:
      app: monoguard-api
  template:
    metadata:
      labels:
        app: monoguard-api
    spec:
      containers:
      - name: api
        image: monoguard/api:latest
        ports:
        - containerPort: 8080
        env:
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: monoguard-secrets
              key: database-url
        - name: JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: monoguard-secrets
              key: jwt-secret
        envFrom:
        - configMapRef:
            name: monoguard-config
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"

---
apiVersion: v1
kind: Service
metadata:
  name: monoguard-api-service
  namespace: monoguard
spec:
  selector:
    app: monoguard-api
  ports:
  - protocol: TCP
    port: 80
    targetPort: 8080
  type: ClusterIP
```

#### Frontend Deployment
```yaml
# k8s/frontend-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: monoguard-frontend
  namespace: monoguard
spec:
  replicas: 2
  selector:
    matchLabels:
      app: monoguard-frontend
  template:
    metadata:
      labels:
        app: monoguard-frontend
    spec:
      containers:
      - name: frontend
        image: monoguard/frontend:latest
        ports:
        - containerPort: 3000
        env:
        - name: NEXT_PUBLIC_API_URL
          value: "https://api.monoguard.com"
        - name: NODE_ENV
          value: "production"
        livenessProbe:
          httpGet:
            path: /
            port: 3000
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /
            port: 3000
          initialDelaySeconds: 5
          periodSeconds: 5
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "256Mi"
            cpu: "200m"

---
apiVersion: v1
kind: Service
metadata:
  name: monoguard-frontend-service
  namespace: monoguard
spec:
  selector:
    app: monoguard-frontend
  ports:
  - protocol: TCP
    port: 80
    targetPort: 3000
  type: ClusterIP
```

#### Ingress Configuration
```yaml
# k8s/ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: monoguard-ingress
  namespace: monoguard
  annotations:
    kubernetes.io/ingress.class: "nginx"
    cert-manager.io/cluster-issuer: "letsencrypt-prod"
    nginx.ingress.kubernetes.io/rate-limit: "100"
spec:
  tls:
  - hosts:
    - monoguard.com
    - api.monoguard.com
    secretName: monoguard-tls
  rules:
  - host: monoguard.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: monoguard-frontend-service
            port:
              number: 80
  - host: api.monoguard.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: monoguard-api-service
            port:
              number: 80
```

### Infrastructure Requirements

#### Environment Sizing
| Environment | API Instances | Workers | Database | Storage |
|-------------|---------------|---------|----------|---------|
| Development | 1 | 1 | Shared | 10GB |
| Staging | 1 | 1 | Small | 50GB |
| Production | 3 | 2 | Managed | 500GB |

#### Production Infrastructure
- **Compute**: Kubernetes cluster with auto-scaling (3-10 nodes)
- **Database**: Managed PostgreSQL with read replicas for reporting
- **Caching**: Redis cluster for session management and query caching
- **Storage**: Block storage for database, object storage for artifacts
- **CDN**: CloudFlare for static asset delivery and DDoS protection
- **Monitoring**: Prometheus + Grafana for metrics, ELK stack for logs

### CI/CD Pipeline

#### GitHub Actions Workflow
```yaml
# .github/workflows/deploy.yml
name: Deploy to Production

on:
  push:
    branches: [main]
  release:
    types: [published]

jobs:
  test:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:15
        env:
          POSTGRES_PASSWORD: postgres
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - name: Run backend tests
        run: |
          cd backend
          go mod download
          go test -v -race -coverprofile=coverage.out ./...
      
      - uses: actions/setup-node@v4
        with:
          node-version: '18'
      - name: Run frontend tests
        run: |
          cd frontend
          npm ci
          npm run test:coverage
          npm run build

  build-and-push:
    needs: test
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    steps:
      - uses: actions/checkout@v4
      
      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3
        
      - name: Login to Container Registry
        uses: docker/login-action@v3
        with:
          registry: ghcr.io
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}
          
      - name: Build and push API
        uses: docker/build-push-action@v5
        with:
          context: ./backend
          push: true
          tags: ghcr.io/monoguard/api:${{ github.sha }},ghcr.io/monoguard/api:latest
          
      - name: Build and push Frontend
        uses: docker/build-push-action@v5
        with:
          context: ./frontend
          push: true
          tags: ghcr.io/monoguard/frontend:${{ github.sha }},ghcr.io/monoguard/frontend:latest

  deploy:
    needs: build-and-push
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    steps:
      - uses: actions/checkout@v4
      
      - name: Deploy to Kubernetes
        uses: azure/k8s-deploy@v1
        with:
          manifests: |
            k8s/api-deployment.yaml
            k8s/frontend-deployment.yaml
          images: |
            ghcr.io/monoguard/api:${{ github.sha }}
            ghcr.io/monoguard/frontend:${{ github.sha }}
          kubeconfig: ${{ secrets.KUBE_CONFIG }}
```

### Monitoring and Observability

#### Health Checks
```go
// backend/internal/api/health.go
func (s *Server) healthHandler(c *gin.Context) {
    health := map[string]interface{}{
        "status":    "healthy",
        "timestamp": time.Now().UTC(),
        "version":   s.version,
        "uptime":    time.Since(s.startTime).String(),
    }
    
    // Check database connectivity
    if err := s.db.Ping(); err != nil {
        health["status"] = "unhealthy"
        health["database"] = "down"
        c.JSON(http.StatusServiceUnavailable, health)
        return
    }
    
    // Check Redis connectivity
    if err := s.redis.Ping().Err(); err != nil {
        health["status"] = "degraded"
        health["cache"] = "down"
    }
    
    c.JSON(http.StatusOK, health)
}
```

#### Prometheus Metrics
```yaml
# k8s/monitoring.yaml
apiVersion: v1
kind: ServiceMonitor
metadata:
  name: monoguard-api
  namespace: monoguard
spec:
  selector:
    matchLabels:
      app: monoguard-api
  endpoints:
  - port: metrics
    path: /metrics
    interval: 30s
```

#### Grafana Dashboard
```json
{
  "dashboard": {
    "title": "MonoGuard Metrics",
    "panels": [
      {
        "title": "Request Rate",
        "type": "graph",
        "targets": [
          {
            "expr": "rate(http_requests_total[5m])",
            "legendFormat": "{{method}} {{endpoint}}"
          }
        ]
      },
      {
        "title": "Response Time",
        "type": "graph",
        "targets": [
          {
            "expr": "histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))",
            "legendFormat": "95th percentile"
          }
        ]
      },
      {
        "title": "Error Rate",
        "type": "singlestat",
        "targets": [
          {
            "expr": "rate(http_requests_total{status!~\"2..\"}[5m]) / rate(http_requests_total[5m])",
            "legendFormat": "Error Rate"
          }
        ]
      }
    ]
  }
}
```

### Backup and Disaster Recovery

#### Database Backup Strategy
```bash
# Automated daily backups
#!/bin/bash
# backup-script.sh

DB_NAME="monoguard"
BACKUP_DIR="/backups/postgresql"
DATE=$(date +%Y%m%d_%H%M%S)

# Create backup
pg_dump -h $DB_HOST -U $DB_USER -d $DB_NAME | \
  gzip > $BACKUP_DIR/monoguard_backup_$DATE.sql.gz

# Upload to cloud storage
aws s3 cp $BACKUP_DIR/monoguard_backup_$DATE.sql.gz \
  s3://monoguard-backups/postgresql/

# Cleanup local backups older than 7 days
find $BACKUP_DIR -name "*.sql.gz" -mtime +7 -delete
```

#### Disaster Recovery Plan
1. **RTO (Recovery Time Objective)**: 4 hours
2. **RPO (Recovery Point Objective)**: 1 hour
3. **Backup Frequency**: Daily full backup, hourly incremental
4. **Failover Strategy**: Multi-region deployment with DNS failover
5. **Data Replication**: Database read replicas in secondary region

### Security Deployment Considerations

#### Secret Management
```yaml
# k8s/secrets.yaml (encrypted with sealed-secrets)
apiVersion: bitnami.com/v1alpha1
kind: SealedSecret
metadata:
  name: monoguard-secrets
  namespace: monoguard
spec:
  encryptedData:
    database-url: <encrypted-value>
    jwt-secret: <encrypted-value>
    oauth-client-secret: <encrypted-value>
```

#### Network Security
- **Network Policies**: Restrict pod-to-pod communication
- **TLS Termination**: All external traffic encrypted with TLS 1.3
- **Rate Limiting**: Protect against DDoS attacks
- **WAF**: Web Application Firewall for additional protection

#### Container Security
- **Image Scanning**: Scan all images for vulnerabilities
- **Non-root Containers**: Run containers with non-root user
- **Read-only Filesystem**: Mount root filesystem as read-only
- **Security Contexts**: Apply appropriate security policies