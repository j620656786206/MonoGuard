package database

import (
	"database/sql"
	"fmt"
	"log"
	"runtime"
	"runtime/debug"
	"strings"
	"time"

	"github.com/monoguard/api/internal/config"
	"github.com/monoguard/api/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)


// DB holds the database connection
type DB struct {
	*gorm.DB
}

// logVersionInfo logs GORM and Go version information for debugging
func logVersionInfo() {
	log.Printf("=== VERSION INFORMATION ===")
	log.Printf("Go Version: %s", runtime.Version())
	log.Printf("GOOS: %s, GOARCH: %s", runtime.GOOS, runtime.GOARCH)
	
	info, ok := debug.ReadBuildInfo()
	if ok {
		log.Printf("Main Module: %s", info.Main.Path)
		for _, mod := range info.Deps {
			if strings.Contains(mod.Path, "gorm") {
				log.Printf("GORM Module: %s@%s", mod.Path, mod.Version)
			}
		}
	}
	log.Printf("=== END VERSION INFO ===")
}

// New creates a new database connection
func New(cfg *config.DatabaseConfig) (*DB, error) {
	// Log version information for debugging
	logVersionInfo()
	
	// Configure GORM logger
	gormLogger := logger.Default
	gormLogger = gormLogger.LogMode(logger.Info)

	var db *gorm.DB
	var err error

	// PostgreSQL only
	dsn := cfg.GetDSN()
	log.Printf("Connecting to PostgreSQL with DSN: %s (password hidden)", 
		fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=%s TimeZone=UTC",
			cfg.Host, cfg.Port, cfg.User, cfg.DBName, cfg.SSLMode))
	
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
		// PostgreSQL optimizations
		DisableForeignKeyConstraintWhenMigrating: false,
		SkipDefaultTransaction: false,
		PrepareStmt: true,
	})
	
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying sql.DB for connection pooling
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying database connection: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxOpenConns(cfg.MaxOpen)
	sqlDB.SetMaxIdleConns(cfg.MaxIdle)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return &DB{db}, nil
}

// AutoMigrate runs database migrations using GORM
func (db *DB) AutoMigrate() error {
	log.Println("Running GORM database migrations...")
	
	// Test database connection first
	sqlDB, err := db.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying database connection: %w", err)
	}
	
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}
	log.Println("Database connection verified")

	// Run GORM migrations using the DB instance properly
	log.Printf("Running GORM AutoMigrate...")
	
	err = db.DB.AutoMigrate(
		&models.Project{},
		&models.DependencyAnalysis{},
		&models.ArchitectureValidation{},
		&models.FileProcessingResult{},
		&models.UploadedFile{},
		&models.PackageJsonFile{},
	)
	
	if err != nil {
		log.Printf("Migration error: %v", err)
		return fmt.Errorf("failed to migrate models: %w", err)
	}
	
	log.Printf("Successfully migrated all models")

	log.Println("GORM database migrations completed successfully")
	return nil
}

// HealthCheck checks if the database is accessible
func (db *DB) HealthCheck() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying database connection: %w", err)
	}

	ctx, cancel := createTimeoutContext(5 * time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	return nil
}

// GetStats returns database connection statistics
func (db *DB) GetStats() sql.DBStats {
	sqlDB, _ := db.DB.DB()
	return sqlDB.Stats()
}

// Close closes the database connection
func (db *DB) Close() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}