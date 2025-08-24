package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Recovery returns a middleware that recovers from panics and logs the error
func Recovery(logger *logrus.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		requestID, _ := c.Get("request_id")
		
		logger.WithFields(logrus.Fields{
			"request_id": requestID,
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"client_ip":  c.ClientIP(),
			"panic":      recovered,
		}).Error("Panic recovered")

		c.JSON(http.StatusInternalServerError, gin.H{
			"success":   false,
			"message":   "Internal server error",
			"timestamp": getCurrentTimestamp(),
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "An unexpected error occurred",
			},
		})
		
		c.Abort()
	})
}