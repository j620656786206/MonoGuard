package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/monoguard/api/internal/config"
	"github.com/monoguard/api/pkg/database"
	"github.com/sirupsen/logrus"
)

// HealthHandler handles health check endpoints
type HealthHandler struct {
	db     *database.DB
	redis  *database.RedisClient
	config *config.Config
	logger *logrus.Logger
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(db *database.DB, redis *database.RedisClient, cfg *config.Config, logger *logrus.Logger) *HealthHandler {
	return &HealthHandler{
		db:     db,
		redis:  redis,
		config: cfg,
		logger: logger,
	}
}

// HealthStatus represents the health status of the service
type HealthStatus struct {
	Status      string                 `json:"status"`
	Service     string                 `json:"service"`
	Version     string                 `json:"version"`
	Environment string                 `json:"environment"`
	Timestamp   string                 `json:"timestamp"`
	Uptime      string                 `json:"uptime"`
	Checks      map[string]CheckResult `json:"checks"`
}

// CheckResult represents the result of a health check
type CheckResult struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Details interface{} `json:"details,omitempty"`
}

var startTime = time.Now()

// Health returns the health status of the service
func (h *HealthHandler) Health(c *gin.Context) {
	checks := make(map[string]CheckResult)
	overallStatus := "healthy"

	// Database health check
	if err := h.db.HealthCheck(); err != nil {
		checks["database"] = CheckResult{
			Status:  "unhealthy",
			Message: err.Error(),
		}
		overallStatus = "unhealthy"
		h.logger.WithError(err).Error("Database health check failed")
	} else {
		checks["database"] = CheckResult{
			Status:  "healthy",
			Message: "Database connection is active",
			Details: h.db.GetStats(),
		}
	}

	// Redis health check (optional)
	if h.redis != nil {
		if err := h.redis.HealthCheck(); err != nil {
			checks["redis"] = CheckResult{
				Status:  "unhealthy",
				Message: err.Error(),
			}
			overallStatus = "unhealthy"
			h.logger.WithError(err).Error("Redis health check failed")
		} else {
			checks["redis"] = CheckResult{
				Status:  "healthy",
				Message: "Redis connection is active",
			}
		}
	} else {
		checks["redis"] = CheckResult{
			Status:  "disabled",
			Message: "Redis is disabled for this environment",
		}
	}

	// Calculate uptime
	uptime := time.Since(startTime)

	healthStatus := HealthStatus{
		Status:      overallStatus,
		Service:     h.config.App.Name,
		Version:     h.config.App.Version,
		Environment: h.config.App.Environment,
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
		Uptime:      uptime.String(),
		Checks:      checks,
	}

	// Set appropriate HTTP status code
	statusCode := http.StatusOK
	if overallStatus == "unhealthy" {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, healthStatus)
}

// Ready returns the readiness status of the service
func (h *HealthHandler) Ready(c *gin.Context) {
	// Check if all critical services are ready
	if err := h.db.HealthCheck(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":    "not ready",
			"message":   "Database not ready",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
		return
	}

	// Only check Redis if it's enabled
	if h.redis != nil {
		if err := h.redis.HealthCheck(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":    "not ready",
				"message":   "Redis not ready",
				"timestamp": time.Now().UTC().Format(time.RFC3339),
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    "ready",
		"message":   "Service is ready to accept requests",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// Live returns the liveness status of the service
func (h *HealthHandler) Live(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "alive",
		"service":   h.config.App.Name,
		"version":   h.config.App.Version,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}