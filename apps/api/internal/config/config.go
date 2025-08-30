package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig   `validate:"required"`
	Database DatabaseConfig `validate:"required"`
	Redis    RedisConfig    `validate:"required"`
	App      AppConfig      `validate:"required"`
	Auth     AuthConfig     `validate:"required"`
	Upload   UploadConfig   `validate:"required"`
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port         string `validate:"required"`
	Host         string `validate:"required"`
	Mode         string `validate:"required,oneof=debug release test"`
	ReadTimeout  int    `validate:"min=1"`
	WriteTimeout int    `validate:"min=1"`
	IdleTimeout  int    `validate:"min=1"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string `validate:"required"`
	Port     int    `validate:"required,min=1,max=65535"`
	User     string `validate:"required"`
	Password string `validate:"required"`
	DBName   string `validate:"required"`
	SSLMode  string `validate:"required,oneof=disable require verify-ca verify-full"`
	MaxOpen  int    `validate:"min=1"`
	MaxIdle  int    `validate:"min=1"`
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string `validate:"required"`
	Port     int    `validate:"required,min=1,max=65535"`
	Password string
	DB       int    `validate:"min=0"`
}

// AppConfig holds application-specific configuration
type AppConfig struct {
	Name        string `validate:"required"`
	Version     string `validate:"required"`
	Environment string `validate:"required,oneof=development staging production"`
	LogLevel    string `validate:"required,oneof=trace debug info warn error fatal panic"`
}

// AuthConfig holds authentication configuration
type AuthConfig struct {
	JWTSecret     string `validate:"required,min=32"`
	TokenDuration int    `validate:"required,min=1"`
}

// UploadConfig holds file upload configuration
type UploadConfig struct {
	Directory string `validate:"required"`
	MaxSize   int64  `validate:"min=1"`
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	config := &Config{
		Server: ServerConfig{
			Port:         getEnvString("PORT", "8080"),
			Host:         getEnvString("HOST", "0.0.0.0"),
			Mode:         getEnvString("GIN_MODE", "debug"),
			ReadTimeout:  getEnvInt("READ_TIMEOUT", 30),
			WriteTimeout: getEnvInt("WRITE_TIMEOUT", 30),
			IdleTimeout:  getEnvInt("IDLE_TIMEOUT", 120),
		},
		Database: DatabaseConfig{
			Host:     getEnvString("DB_HOST", "localhost"),
			Port:     getEnvInt("DB_PORT", 5432),
			User:     getEnvString("DB_USER", "postgres"),
			Password: getEnvString("DB_PASSWORD", "postgres"),
			DBName:   getEnvString("DB_NAME", "monoguard"),
			SSLMode:  getEnvString("DB_SSL_MODE", "disable"),
			MaxOpen:  getEnvInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdle:  getEnvInt("DB_MAX_IDLE_CONNS", 10),
		},
		Redis: RedisConfig{
			Host:     getEnvString("REDIS_HOST", "localhost"),
			Port:     getEnvInt("REDIS_PORT", 6379),
			Password: getEnvString("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
		},
		App: AppConfig{
			Name:        getEnvString("APP_NAME", "MonoGuard API"),
			Version:     getEnvString("APP_VERSION", "0.1.0"),
			Environment: getEnvString("APP_ENV", "development"),
			LogLevel:    getEnvString("LOG_LEVEL", "info"),
		},
		Auth: AuthConfig{
			JWTSecret:     getEnvString("JWT_SECRET", generateDefaultSecret()),
			TokenDuration: getEnvInt("JWT_DURATION_HOURS", 24),
		},
		Upload: UploadConfig{
			Directory: getEnvString("UPLOAD_DIR", "./uploads"),
			MaxSize:   int64(getEnvInt("UPLOAD_MAX_SIZE_MB", 50)) * 1024 * 1024, // Convert MB to bytes
		},
	}

	// Validate configuration
	validator := validator.New()
	if err := validator.Struct(config); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	// Additional validation for PostgreSQL (production or when explicitly configured)
	if config.Database.Host != "sqlite" && config.Database.Host != "" && config.Database.Host != "localhost" {
		if config.Database.Host == "" || config.Database.User == "" || config.Database.Password == "" || config.Database.DBName == "" {
			return nil, fmt.Errorf("PostgreSQL configuration incomplete: host=%s, user=%s, dbname=%s (password length=%d)",
				config.Database.Host, config.Database.User, config.Database.DBName, len(config.Database.Password))
		}
		fmt.Printf("PostgreSQL config validated: host=%s, user=%s, dbname=%s, port=%d, sslmode=%s\n",
			config.Database.Host, config.Database.User, config.Database.DBName, config.Database.Port, config.Database.SSLMode)
	}

	return config, nil
}

// GetDSN returns the database connection string
func (c *DatabaseConfig) GetDSN() string {
	// For Railway PostgreSQL, ensure proper URL encoding and timezone
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=UTC",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
	
	return dsn
}

// GetRedisAddr returns the Redis address
func (c *RedisConfig) GetRedisAddr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// Helper functions for environment variable parsing
func getEnvString(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}


// generateDefaultSecret generates a default JWT secret for development
func generateDefaultSecret() string {
	return "your-super-secret-jwt-key-change-this-in-production-32-chars-minimum"
}