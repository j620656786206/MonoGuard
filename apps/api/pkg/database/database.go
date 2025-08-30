package database

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/monoguard/api/internal/config"
	"github.com/monoguard/api/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// IsMigrationMode indicates if we're currently in migration mode
var (
	isMigrationMode bool
	migrationMutex  sync.RWMutex
)

// SetMigrationMode sets the migration mode flag
func SetMigrationMode(mode bool) {
	migrationMutex.Lock()
	defer migrationMutex.Unlock()
	isMigrationMode = mode
}

// IsMigrationMode returns the current migration mode status
func IsMigrationMode() bool {
	migrationMutex.RLock()
	defer migrationMutex.RUnlock()
	return isMigrationMode
}

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
		log.Printf("Connecting to PostgreSQL with DSN: %s (password hidden)", 
			fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=%s TimeZone=UTC",
				cfg.Host, cfg.Port, cfg.User, cfg.DBName, cfg.SSLMode))
		
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: gormLogger,
			NowFunc: func() time.Time {
				return time.Now().UTC()
			},
			// Railway PostgreSQL optimizations
			DisableForeignKeyConstraintWhenMigrating: true,
			SkipDefaultTransaction: false,
			PrepareStmt: false,
			// Disable automatic pluralization to prevent table name issues
			NamingStrategy: nil,
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
	
	// Set migration mode to disable hooks globally
	SetMigrationMode(true)
	defer SetMigrationMode(false)
	
	// Test database connection first
	sqlDB, err := db.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying database connection: %w", err)
	}
	
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}
	log.Println("Database connection verified")

	// Create a migration session without hooks for extra safety
	migrationDB := db.DB.Session(&gorm.Session{
		SkipHooks: true,
	})

	// Migrate models one by one to identify problematic model
	models := []interface{}{
		&models.Project{},
		&models.DependencyAnalysis{},
		&models.ArchitectureValidation{},
		&models.HealthScore{},
		&models.UploadedFile{},
		&models.FileProcessingResult{},
		&models.PackageJsonFile{},
		&models.PackageJSONAnalysis{},
	}

	for i, model := range models {
		modelName := fmt.Sprintf("%T", model)
		log.Printf("Migrating model %d/%d: %s", i+1, len(models), modelName)
		
		if err := migrationDB.AutoMigrate(model); err != nil {
			log.Printf("Migration error for %s: %v", modelName, err)
			return fmt.Errorf("failed to migrate %s: %w", modelName, err)
		}
		log.Printf("Successfully migrated: %s", modelName)
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