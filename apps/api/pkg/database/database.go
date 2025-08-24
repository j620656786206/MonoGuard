package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/monoguard/api/internal/config"
	"github.com/monoguard/api/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB holds the database connection
type DB struct {
	*gorm.DB
}

// New creates a new database connection
func New(cfg *config.DatabaseConfig) (*DB, error) {
	// Configure GORM logger
	gormLogger := logger.Default
	gormLogger = gormLogger.LogMode(logger.Info)

	var db *gorm.DB
	var err error

	// Use SQLite for development, PostgreSQL for production
	if cfg.Host == "sqlite" || cfg.Host == "" {
		// SQLite mode for development
		dbFile := cfg.DBName
		if dbFile == "" {
			dbFile = "monoguard.db"
		}
		db, err = gorm.Open(sqlite.Open(dbFile), &gorm.Config{
			Logger: gormLogger,
			NowFunc: func() time.Time {
				return time.Now().UTC()
			},
		})
	} else {
		// PostgreSQL mode
		dsn := cfg.GetDSN()
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: gormLogger,
			NowFunc: func() time.Time {
				return time.Now().UTC()
			},
		})
	}
	
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying sql.DB for connection pooling
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying database connection: %w", err)
	}

	// Configure connection pool (only for PostgreSQL)
	if cfg.Host != "sqlite" && cfg.Host != "" {
		sqlDB.SetMaxOpenConns(cfg.MaxOpen)
		sqlDB.SetMaxIdleConns(cfg.MaxIdle)
		sqlDB.SetConnMaxLifetime(time.Hour)
	}

	return &DB{db}, nil
}

// AutoMigrate runs database migrations
func (db *DB) AutoMigrate() error {
	log.Println("Running database migrations...")
	
	// Migrate models one by one to identify problematic model
	models := []interface{}{
		&models.Project{},
		&models.DependencyAnalysis{},
		&models.ArchitectureValidation{},
		&models.HealthScore{},
	}
	
	for i, model := range models {
		log.Printf("Migrating model %d...", i+1)
		err := db.DB.AutoMigrate(model)
		if err != nil {
			return fmt.Errorf("failed to migrate model %d: %w", i+1, err)
		}
		log.Printf("Successfully migrated model %d", i+1)
	}

	log.Println("Database migrations completed successfully")
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