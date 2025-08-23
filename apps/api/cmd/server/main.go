package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Set gin mode
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.DebugMode)
	}

	// Create gin router
	router := gin.Default()

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "monoguard-api",
			"version": "0.1.0",
		})
	})

	// Placeholder API routes - to be implemented
	api := router.Group("/api/v1")
	{
		api.GET("/projects", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Projects endpoint - Coming Soon!"})
		})
		
		api.GET("/analysis", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Analysis endpoint - Coming Soon!"})
		})
		
		api.GET("/dependencies", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Dependencies endpoint - Coming Soon!"})
		})
	}

	// Get port from environment or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting MonoGuard API server on port %s", port)
	log.Fatal(router.Run(":" + port))
}