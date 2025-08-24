package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// Logger returns a Gin middleware for logging HTTP requests
func Logger(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate request ID
		requestID := uuid.New().String()
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)

		// Start timer
		startTime := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log after request
		param := gin.LogFormatterParams{
			Request:    c.Request,
			TimeStamp:  time.Now(),
			Latency:    time.Since(startTime),
			ClientIP:   c.ClientIP(),
			Method:     c.Request.Method,
			StatusCode: c.Writer.Status(),
			ErrorMessage: c.Errors.ByType(gin.ErrorTypePrivate).String(),
			BodySize:     c.Writer.Size(),
			Keys:         c.Keys,
		}

		if raw != "" {
			param.Path = path + "?" + raw
		} else {
			param.Path = path
		}

		// Create structured log entry
		logEntry := logger.WithFields(logrus.Fields{
			"request_id":   requestID,
			"method":       param.Method,
			"path":         param.Path,
			"status_code":  param.StatusCode,
			"latency_ms":   param.Latency.Milliseconds(),
			"client_ip":    param.ClientIP,
			"user_agent":   c.Request.UserAgent(),
			"body_size":    param.BodySize,
		})

		// Log based on status code
		if param.StatusCode >= 500 {
			logEntry.WithField("error", param.ErrorMessage).Error("Server error")
		} else if param.StatusCode >= 400 {
			logEntry.Warn("Client error")
		} else {
			logEntry.Info("Request completed")
		}
	}
}

// RequestID middleware adds request ID to context
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}