# MonoGuard Security Guide

## Security Architecture Overview

### Security-First Design Principles
1. **Zero Trust Architecture**: Never trust, always verify
2. **Defense in Depth**: Multiple layers of security controls
3. **Principle of Least Privilege**: Minimal necessary access rights
4. **Data Minimization**: Collect and store only essential data
5. **Secure by Default**: Secure configurations as defaults

### Threat Model
- **Adversaries**: Malicious actors seeking to access customer code or data
- **Assets**: Customer source code, dependency analysis data, user credentials
- **Attack Vectors**: API exploitation, container escape, supply chain attacks
- **Impact**: Data breach, service disruption, reputation damage

## Data Protection

### Source Code Security
**Critical Principle**: Never persist customer source code

```go
// Secure analysis processing
type SecureAnalyzer struct {
    tempDir     string
    maxLifetime time.Duration
}

func (a *SecureAnalyzer) AnalyzeProject(ctx context.Context, projectPath string) (*Analysis, error) {
    // Create isolated temporary workspace
    workspace, err := a.createSecureWorkspace()
    if err != nil {
        return nil, err
    }
    defer a.cleanupWorkspace(workspace) // Always cleanup
    
    // Set automatic cleanup timer
    timer := time.AfterFunc(a.maxLifetime, func() {
        a.forceCleanup(workspace)
    })
    defer timer.Stop()
    
    // Process in memory only
    analysis, err := a.processInMemory(ctx, projectPath)
    if err != nil {
        return nil, err
    }
    
    return analysis, nil
}

func (a *SecureAnalyzer) createSecureWorkspace() (string, error) {
    // Create temporary directory with restricted permissions
    tempDir, err := os.MkdirTemp("", "monoguard-*")
    if err != nil {
        return "", err
    }
    
    // Set strict permissions (owner only)
    if err := os.Chmod(tempDir, 0700); err != nil {
        os.RemoveAll(tempDir)
        return "", err
    }
    
    return tempDir, nil
}
```

### Data Encryption
```go
// Encrypt sensitive configuration data
func EncryptConfig(data []byte, key []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }
    
    nonce := make([]byte, gcm.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return nil, err
    }
    
    ciphertext := gcm.Seal(nonce, nonce, data, nil)
    return ciphertext, nil
}
```

### Memory Protection
```go
// Secure memory handling for sensitive data
type SecureString struct {
    data []byte
}

func NewSecureString(s string) *SecureString {
    data := make([]byte, len(s))
    copy(data, s)
    return &SecureString{data: data}
}

func (s *SecureString) Clear() {
    // Overwrite memory with random data
    if _, err := rand.Read(s.data); err != nil {
        // Fallback to zero fill
        for i := range s.data {
            s.data[i] = 0
        }
    }
}

func (s *SecureString) String() string {
    return string(s.data)
}

// Use finalizer to ensure cleanup
func (s *SecureString) setFinalizer() {
    runtime.SetFinalizer(s, (*SecureString).Clear)
}
```

## Authentication & Authorization

### OAuth2 Integration
```go
// Secure OAuth2 configuration
type OAuthConfig struct {
    ClientID     string
    ClientSecret *SecureString
    RedirectURL  string
    Scopes       []string
    State        string
}

func (c *OAuthConfig) GetAuthURL() string {
    params := url.Values{
        "client_id":     {c.ClientID},
        "redirect_uri":  {c.RedirectURL},
        "scope":         {strings.Join(c.Scopes, " ")},
        "response_type": {"code"},
        "state":         {c.State},
    }
    
    return fmt.Sprintf("https://github.com/login/oauth/authorize?%s", params.Encode())
}

// Validate state parameter to prevent CSRF
func (c *OAuthConfig) ValidateState(receivedState string) error {
    if subtle.ConstantTimeCompare([]byte(c.State), []byte(receivedState)) != 1 {
        return errors.New("invalid state parameter")
    }
    return nil
}
```

### JWT Token Management
```go
// Secure JWT implementation
type JWTManager struct {
    signingKey   *SecureString
    refreshKey   *SecureString
    accessExpiry time.Duration
    refreshExpiry time.Duration
}

func (j *JWTManager) GenerateTokenPair(userID string, roles []string) (*TokenPair, error) {
    // Access token with short expiry
    accessToken, err := j.generateAccessToken(userID, roles, j.accessExpiry)
    if err != nil {
        return nil, err
    }
    
    // Refresh token with longer expiry
    refreshToken, err := j.generateRefreshToken(userID, j.refreshExpiry)
    if err != nil {
        return nil, err
    }
    
    return &TokenPair{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
        ExpiresIn:    int64(j.accessExpiry.Seconds()),
    }, nil
}

func (j *JWTManager) ValidateToken(tokenString string, tokenType TokenType) (*Claims, error) {
    var key []byte
    switch tokenType {
    case AccessToken:
        key = j.signingKey.data
    case RefreshToken:
        key = j.refreshKey.data
    default:
        return nil, errors.New("invalid token type")
    }
    
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return key, nil
    })
    
    if err != nil {
        return nil, err
    }
    
    claims, ok := token.Claims.(*Claims)
    if !ok || !token.Valid {
        return nil, errors.New("invalid token")
    }
    
    return claims, nil
}
```

### Role-Based Access Control (RBAC)
```go
// RBAC implementation
type Permission string

const (
    ReadProject    Permission = "project:read"
    WriteProject   Permission = "project:write"
    DeleteProject  Permission = "project:delete"
    ManageUsers    Permission = "users:manage"
    ViewReports    Permission = "reports:view"
    RunAnalysis    Permission = "analysis:run"
)

type Role struct {
    Name        string       `json:"name"`
    Permissions []Permission `json:"permissions"`
}

var (
    ViewerRole = Role{
        Name:        "viewer",
        Permissions: []Permission{ReadProject, ViewReports},
    }
    
    DeveloperRole = Role{
        Name: "developer",
        Permissions: []Permission{
            ReadProject, WriteProject, ViewReports, RunAnalysis,
        },
    }
    
    AdminRole = Role{
        Name: "admin",
        Permissions: []Permission{
            ReadProject, WriteProject, DeleteProject,
            ViewReports, RunAnalysis, ManageUsers,
        },
    }
)

// Middleware for permission checking
func RequirePermission(permission Permission) gin.HandlerFunc {
    return func(c *gin.Context) {
        claims, exists := c.Get("claims")
        if !exists {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
            c.Abort()
            return
        }
        
        userClaims := claims.(*Claims)
        if !hasPermission(userClaims.Roles, permission) {
            c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
            c.Abort()
            return
        }
        
        c.Next()
    }
}
```

## API Security

### Rate Limiting
```go
// Implement rate limiting middleware
type RateLimiter struct {
    store   map[string]*TokenBucket
    mutex   sync.RWMutex
    cleanup time.Duration
}

type TokenBucket struct {
    tokens    int
    capacity  int
    refillRate int
    lastRefill time.Time
}

func NewRateLimiter() *RateLimiter {
    rl := &RateLimiter{
        store:   make(map[string]*TokenBucket),
        cleanup: time.Hour,
    }
    
    go rl.cleanupExpired()
    return rl
}

func (rl *RateLimiter) Allow(key string, capacity, refillRate int) bool {
    rl.mutex.Lock()
    defer rl.mutex.Unlock()
    
    bucket, exists := rl.store[key]
    if !exists {
        bucket = &TokenBucket{
            tokens:     capacity,
            capacity:   capacity,
            refillRate: refillRate,
            lastRefill: time.Now(),
        }
        rl.store[key] = bucket
    }
    
    // Refill tokens
    now := time.Now()
    elapsed := now.Sub(bucket.lastRefill)
    tokensToAdd := int(elapsed.Seconds()) * bucket.refillRate
    
    bucket.tokens = min(bucket.capacity, bucket.tokens+tokensToAdd)
    bucket.lastRefill = now
    
    if bucket.tokens > 0 {
        bucket.tokens--
        return true
    }
    
    return false
}

// Rate limiting middleware
func RateLimitMiddleware(capacity, refillRate int) gin.HandlerFunc {
    limiter := NewRateLimiter()
    
    return func(c *gin.Context) {
        key := c.ClientIP()
        if userID, exists := c.Get("user_id"); exists {
            key = fmt.Sprintf("user:%v", userID)
        }
        
        if !limiter.Allow(key, capacity, refillRate) {
            c.JSON(http.StatusTooManyRequests, gin.H{
                "error": "rate limit exceeded",
                "retry_after": 60,
            })
            c.Abort()
            return
        }
        
        c.Next()
    }
}
```

### Input Validation
```go
// Comprehensive input validation
type Validator struct {
    validate *validator.Validate
}

func NewValidator() *Validator {
    validate := validator.New()
    
    // Register custom validations
    validate.RegisterValidation("project_name", validateProjectName)
    validate.RegisterValidation("safe_path", validateSafePath)
    
    return &Validator{validate: validate}
}

func validateProjectName(fl validator.FieldLevel) bool {
    name := fl.Field().String()
    // Only allow alphanumeric, hyphens, underscores
    matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, name)
    return matched && len(name) >= 3 && len(name) <= 50
}

func validateSafePath(fl validator.FieldLevel) bool {
    path := fl.Field().String()
    // Prevent path traversal attacks
    if strings.Contains(path, "..") {
        return false
    }
    if strings.Contains(path, "~") {
        return false
    }
    return true
}

// Request validation middleware
func ValidateJSON(model interface{}) gin.HandlerFunc {
    return func(c *gin.Context) {
        if err := c.ShouldBindJSON(model); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{
                "error": "invalid request format",
                "details": err.Error(),
            })
            c.Abort()
            return
        }
        
        validator := NewValidator()
        if err := validator.validate.Struct(model); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{
                "error": "validation failed",
                "details": err.Error(),
            })
            c.Abort()
            return
        }
        
        c.Set("validated_model", model)
        c.Next()
    }
}
```

### SQL Injection Prevention
```go
// Use parameterized queries with GORM
func (r *ProjectRepository) GetProjectsByUser(userID string) ([]Project, error) {
    var projects []Project
    
    // GORM automatically handles parameterization
    err := r.db.Where("user_id = ? AND deleted_at IS NULL", userID).Find(&projects).Error
    if err != nil {
        return nil, err
    }
    
    return projects, nil
}

// For raw queries, always use parameters
func (r *ProjectRepository) GetProjectStats(projectID string) (*ProjectStats, error) {
    var stats ProjectStats
    
    query := `
        SELECT 
            COUNT(*) as total_packages,
            COUNT(CASE WHEN status = 'error' THEN 1 END) as error_count,
            AVG(health_score) as avg_health_score
        FROM project_analyses 
        WHERE project_id = $1 AND created_at > $2
    `
    
    err := r.db.Raw(query, projectID, time.Now().AddDate(0, -1, 0)).Scan(&stats).Error
    return &stats, err
}
```

## Infrastructure Security

### Container Security
```dockerfile
# Secure Dockerfile practices
FROM golang:1.21-alpine AS builder

# Create non-root user
RUN adduser -D -s /bin/sh monoguard

# Use specific versions
RUN apk add --no-cache ca-certificates=20230506-r0

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/api/main.go

# Production stage
FROM scratch

# Copy CA certificates for HTTPS
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy user information
COPY --from=builder /etc/passwd /etc/passwd

# Copy binary
COPY --from=builder /app/main /main

# Use non-root user
USER monoguard

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["/main", "healthcheck"]

ENTRYPOINT ["/main"]
```

### Kubernetes Security Policies
```yaml
# Security contexts and policies
apiVersion: v1
kind: Pod
metadata:
  name: monoguard-api
spec:
  securityContext:
    runAsNonRoot: true
    runAsUser: 65534  # nobody user
    fsGroup: 65534
    seccompProfile:
      type: RuntimeDefault
  containers:
  - name: api
    image: monoguard/api:latest
    securityContext:
      allowPrivilegeEscalation: false
      readOnlyRootFilesystem: true
      runAsNonRoot: true
      capabilities:
        drop:
        - ALL
    volumeMounts:
    - name: tmp
      mountPath: /tmp
      readOnly: false
    - name: var-run
      mountPath: /var/run
      readOnly: false
    resources:
      limits:
        memory: "512Mi"
        cpu: "500m"
      requests:
        memory: "256Mi"
        cpu: "250m"
  volumes:
  - name: tmp
    emptyDir: {}
  - name: var-run
    emptyDir: {}
```

### Network Security
```yaml
# Network policies
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: monoguard-api-policy
  namespace: monoguard
spec:
  podSelector:
    matchLabels:
      app: monoguard-api
  policyTypes:
  - Ingress
  - Egress
  ingress:
  - from:
    - namespaceSelector:
        matchLabels:
          name: ingress-nginx
    ports:
    - protocol: TCP
      port: 8080
  - from:
    - podSelector:
        matchLabels:
          app: monoguard-frontend
    ports:
    - protocol: TCP
      port: 8080
  egress:
  - to:
    - podSelector:
        matchLabels:
          app: postgres
    ports:
    - protocol: TCP
      port: 5432
  - to:
    - podSelector:
        matchLabels:
          app: redis
    ports:
    - protocol: TCP
      port: 6379
  # Allow HTTPS outbound for OAuth
  - to: []
    ports:
    - protocol: TCP
      port: 443
```

## Compliance and Auditing

### Audit Logging
```go
// Comprehensive audit logging
type AuditLogger struct {
    logger *slog.Logger
}

type AuditEvent struct {
    Timestamp time.Time `json:"timestamp"`
    UserID    string    `json:"user_id"`
    Action    string    `json:"action"`
    Resource  string    `json:"resource"`
    Result    string    `json:"result"`
    IPAddress string    `json:"ip_address"`
    UserAgent string    `json:"user_agent"`
    Details   map[string]interface{} `json:"details,omitempty"`
}

func (a *AuditLogger) LogEvent(event AuditEvent) {
    a.logger.Info("audit_event",
        "user_id", event.UserID,
        "action", event.Action,
        "resource", event.Resource,
        "result", event.Result,
        "ip_address", event.IPAddress,
        "user_agent", event.UserAgent,
        "details", event.Details,
    )
}

// Audit middleware
func AuditMiddleware(auditor *AuditLogger) gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        
        // Process request
        c.Next()
        
        // Log audit event
        event := AuditEvent{
            Timestamp: start,
            UserID:    getUserID(c),
            Action:    fmt.Sprintf("%s %s", c.Request.Method, c.Request.URL.Path),
            Resource:  getResourceID(c),
            Result:    getResult(c.Writer.Status()),
            IPAddress: c.ClientIP(),
            UserAgent: c.Request.UserAgent(),
            Details:   getAuditDetails(c),
        }
        
        auditor.LogEvent(event)
    }
}
```

### GDPR Compliance
```go
// Data export for GDPR compliance
func (s *UserService) ExportUserData(userID string) (*UserDataExport, error) {
    export := &UserDataExport{
        UserID:    userID,
        Timestamp: time.Now(),
    }
    
    // Export user profile
    user, err := s.userRepo.GetByID(userID)
    if err != nil {
        return nil, err
    }
    export.Profile = user
    
    // Export projects
    projects, err := s.projectRepo.GetByUserID(userID)
    if err != nil {
        return nil, err
    }
    export.Projects = projects
    
    // Export analysis history (metadata only, no source code)
    analyses, err := s.analysisRepo.GetMetadataByUserID(userID)
    if err != nil {
        return nil, err
    }
    export.AnalysisHistory = analyses
    
    return export, nil
}

// Data deletion for GDPR compliance
func (s *UserService) DeleteUserData(userID string) error {
    tx := s.db.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()
    
    // Delete in reverse dependency order
    if err := s.analysisRepo.DeleteByUserID(tx, userID); err != nil {
        tx.Rollback()
        return err
    }
    
    if err := s.projectRepo.DeleteByUserID(tx, userID); err != nil {
        tx.Rollback()
        return err
    }
    
    if err := s.userRepo.Delete(tx, userID); err != nil {
        tx.Rollback()
        return err
    }
    
    return tx.Commit().Error
}
```

## Security Monitoring and Incident Response

### Security Monitoring
```go
// Security event detection
type SecurityMonitor struct {
    logger    *slog.Logger
    alerter   Alerter
    thresholds map[string]int
}

func (s *SecurityMonitor) MonitorFailedLogins(userID, ip string) {
    key := fmt.Sprintf("failed_login:%s:%s", userID, ip)
    count := s.incrementCounter(key, time.Hour)
    
    if count >= s.thresholds["failed_login"] {
        s.alerter.SendAlert(Alert{
            Type:     "security",
            Severity: "high",
            Message:  fmt.Sprintf("Multiple failed login attempts for user %s from IP %s", userID, ip),
            Details: map[string]interface{}{
                "user_id": userID,
                "ip":      ip,
                "count":   count,
            },
        })
    }
}

func (s *SecurityMonitor) MonitorSuspiciousActivity(userID string, actions []string) {
    // Check for privilege escalation attempts
    if containsPrivilegedActions(actions) {
        s.alerter.SendAlert(Alert{
            Type:     "security",
            Severity: "critical",
            Message:  fmt.Sprintf("Potential privilege escalation attempt by user %s", userID),
            Details: map[string]interface{}{
                "user_id": userID,
                "actions": actions,
            },
        })
    }
}
```

### Incident Response Procedures
1. **Detection**: Automated monitoring alerts + manual reporting
2. **Assessment**: Classify incident severity and impact
3. **Containment**: Isolate affected systems and prevent spread
4. **Eradication**: Remove threat and fix vulnerabilities
5. **Recovery**: Restore services and validate security
6. **Lessons Learned**: Document incident and improve defenses

### Security Testing
```bash
# Security testing checklist

# 1. Dependency scanning
go mod tidy
govulncheck ./...

# 2. Static analysis
gosec ./...
semgrep --config=auto .

# 3. Container scanning
docker scan monoguard/api:latest

# 4. Dynamic testing
curl -X POST http://localhost:8080/api/v1/projects \
  -H "Content-Type: application/json" \
  -d '{"name": "test<script>alert(1)</script>"}'

# 5. Authentication testing
curl -H "Authorization: Bearer invalid-token" \
  http://localhost:8080/api/v1/projects

# 6. Rate limiting testing
for i in {1..100}; do
  curl http://localhost:8080/api/v1/health &
done
```