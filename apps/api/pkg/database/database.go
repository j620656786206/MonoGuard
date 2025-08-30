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
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
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

// AutoMigrate runs database migrations using native SQL to avoid GORM compatibility issues
func (db *DB) AutoMigrate() error {
	log.Println("Running database migrations using native SQL...")
	
	// Test database connection first
	sqlDB, err := db.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying database connection: %w", err)
	}
	
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}
	log.Println("Database connection verified")

	// Create tables using native SQL to avoid GORM "insufficient arguments" issue
	tables := []struct {
		name string
		sql  string
	}{
		{
			name: "projects_simple",
			sql: `CREATE TABLE IF NOT EXISTS projects_simple (
				id VARCHAR(255) PRIMARY KEY,
				name VARCHAR(255) NOT NULL,
				description TEXT,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			)`,
		},
		{
			name: "projects", 
			sql: `CREATE TABLE IF NOT EXISTS projects (
				id VARCHAR(255) PRIMARY KEY,
				name VARCHAR(255) NOT NULL,
				description TEXT,
				repository_url VARCHAR(255),
				branch VARCHAR(255),
				status VARCHAR(255),
				health_score INTEGER DEFAULT 0,
				owner_id VARCHAR(255),
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			)`,
		},
	}

	for _, table := range tables {
		log.Printf("Creating table: %s", table.name)
		
		if err := db.Exec(table.sql).Error; err != nil {
			log.Printf("Failed to create table %s: %v", table.name, err)
			return fmt.Errorf("failed to create table %s: %w", table.name, err)
		}
		
		log.Printf("Successfully created/verified table: %s", table.name)
	}

	log.Println("Database migrations completed successfully using native SQL")
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