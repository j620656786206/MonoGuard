package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORS returns a CORS middleware configured for the application
func CORS() gin.HandlerFunc {
	config := cors.DefaultConfig()
	
	// Allow all origins in development, configure properly for production
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	config.AllowHeaders = []string{
		"Origin",
		"Content-Length",
		"Content-Type",
		"Authorization",
		"X-Requested-With",
		"Accept",
		"Accept-Encoding",
		"Accept-Language",
		"Connection",
		"Host",
		"Referer",
		"User-Agent",
	}
	config.ExposeHeaders = []string{
		"Content-Length",
		"Content-Type",
		"X-Request-ID",
		"X-Response-Time",
	}
	config.AllowCredentials = true
	config.MaxAge = 86400 // 24 hours

	return cors.New(config)
}