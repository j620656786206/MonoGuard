package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/monoguard/api/internal/config"
	"github.com/monoguard/api/internal/handlers"
	"github.com/monoguard/api/internal/middleware"
	"github.com/monoguard/api/internal/repository"
	"github.com/monoguard/api/internal/services"
	"github.com/monoguard/api/pkg/database"
	"github.com/sirupsen/logrus"
)

// App represents the application
type App struct {
	config *config.Config
	logger *logrus.Logger
	server *http.Server
	db     *database.DB
	redis  *database.RedisClient
}

// New creates a new application instance
func New() (*App, error) {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	// Setup logger
	logger := setupLogger(cfg.App.LogLevel)
	logger.WithFields(logrus.Fields{
		"service":     cfg.App.Name,
		"version":     cfg.App.Version,
		"environment": cfg.App.Environment,
	}).Info("Starting application")

	// Connect to database
	db, err := database.New(&cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Run database migrations
	if err := db.AutoMigrate(); err != nil {
		return nil, fmt.Errorf("failed to run database migrations: %w", err)
	}

	// Connect to Redis (optional for development)
	var redis *database.RedisClient
	if cfg.App.Environment == "production" || cfg.Redis.Host != "" {
		redis, err = database.NewRedisClient(&cfg.Redis)
		if err != nil {
			logger.WithError(err).Warn("Failed to connect to Redis, continuing without cache")
			redis = nil
		}
	}

	// Set Gin mode
	gin.SetMode(cfg.Server.Mode)

	app := &App{
		config: cfg,
		logger: logger,
		db:     db,
		redis:  redis,
	}

	app.setupServer()

	return app, nil
}

// setupLogger configures the logger
func setupLogger(logLevel string) *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})

	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)

	return logger
}

// setupServer configures the HTTP server and routes
func (a *App) setupServer() {
	router := gin.New()

	// Middleware
	router.Use(middleware.RequestID())
	router.Use(middleware.Logger(a.logger))
	router.Use(middleware.Recovery(a.logger))
	router.Use(middleware.CORS())

	// Health check routes
	healthHandler := handlers.NewHealthHandler(a.db, a.redis, a.config, a.logger)
	router.GET("/health", healthHandler.Health)
	router.GET("/ready", healthHandler.Ready)
	router.GET("/live", healthHandler.Live)
	router.GET("/debug/db", healthHandler.DebugDB) // Temporary diagnostic endpoint

	// Initialize repositories
	projectRepo := repository.NewProjectRepository(a.db)
	analysisRepo := repository.NewAnalysisRepository(a.db)

	// Initialize services
	analyzer := services.NewDependencyAnalyzer(a.logger)
	projectService := services.NewProjectService(projectRepo, analysisRepo, analyzer, a.logger)
	circularDetectorService := services.NewCircularDetectorService(projectRepo, analysisRepo, a.logger)
	layerValidatorService := services.NewLayerValidatorService(projectRepo, analysisRepo, a.logger)
	uploadService := services.NewUploadService(a.db, a.logger)
	basicEngine := services.NewBasicAnalysisEngine(a.logger)
	
	integratedAnalysis := services.NewIntegratedAnalysisService(
		basicEngine,
		circularDetectorService,
		layerValidatorService,
		uploadService,
		projectRepo,
		analysisRepo,
		a.logger,
	)

	// Initialize handlers
	projectHandler := handlers.NewProjectHandler(projectService, a.logger)
	analysisHandler := handlers.NewAnalysisHandler(analysisRepo, integratedAnalysis, a.logger)
	circularHandler := handlers.NewCircularHandler(circularDetectorService, a.logger)
	architectureHandler := handlers.NewArchitectureHandler(layerValidatorService, a.logger)
	uploadHandler := handlers.NewUploadHandler(uploadService, a.logger, a.config.Upload.Directory)
	githubHandler := handlers.NewGitHubHandler(integratedAnalysis, uploadService, a.logger)

	// API routes (simplified for MVP - no session management)
	v1 := router.Group("/api/v1")
	{

		// Project routes
		projects := v1.Group("/projects")
		{
			projects.POST("", projectHandler.CreateProject)
			projects.GET("", projectHandler.GetProjects)
			projects.GET("/:id", projectHandler.GetProject)
			projects.PUT("/:id", projectHandler.UpdateProject)
			projects.DELETE("/:id", projectHandler.DeleteProject)
			projects.POST("/:id/analyze", projectHandler.TriggerAnalysis)
			
			// Project-specific analysis routes
			projects.GET("/:id/analyses/dependencies", analysisHandler.GetProjectDependencyAnalyses)
			projects.GET("/:id/analyses/dependencies/latest", analysisHandler.GetLatestDependencyAnalysis)
			projects.GET("/:id/health-score/latest", analysisHandler.GetLatestHealthScore)
		}

		// Owner-specific routes
		v1.GET("/projects/owner/:ownerId", projectHandler.GetProjectsByOwner)

		// Circular dependency routes
		v1.POST("/projects/:id/detect-circular-deps", circularHandler.DetectCircularDependencies)
		v1.GET("/projects/:id/circular-dependencies/:analysis_id", circularHandler.GetCircularDependencyAnalysis)

		// Architecture validation routes
		v1.POST("/projects/:id/validate-architecture", architectureHandler.ValidateArchitecture)
		v1.GET("/projects/:id/architecture/graph", architectureHandler.GetArchitectureGraph)
		v1.GET("/projects/:id/architecture/violations", architectureHandler.GetArchitectureViolations)
		v1.POST("/projects/:id/config/validate", architectureHandler.ValidateConfig)

		// Analysis routes
		analysis := v1.Group("/analysis")
		{
			analysis.GET("/dependencies/:id", analysisHandler.GetDependencyAnalysis)
			analysis.GET("/architecture/:id", analysisHandler.GetArchitectureValidation)
			analysis.POST("/comprehensive/:uploadId", analysisHandler.StartComprehensiveAnalysis)
			analysis.POST("/upload", analysisHandler.AnalyzeUploadedFiles)
			analysis.POST("/github", githubHandler.AnalyzeGitHubRepository)
		}

		// Upload routes
		upload := v1.Group("/upload")
		{
			upload.POST("", uploadHandler.UploadFiles)
			upload.GET("/:id", uploadHandler.GetUploadResult)
			upload.GET("/files/:filename", uploadHandler.GetUploadedFile)
			upload.POST("/cleanup", uploadHandler.CleanupOldFiles)
		}
	}

	// Setup HTTP server
	a.server = &http.Server{
		Addr:         fmt.Sprintf("%s:%s", a.config.Server.Host, a.config.Server.Port),
		Handler:      router,
		ReadTimeout:  time.Duration(a.config.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(a.config.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(a.config.Server.IdleTimeout) * time.Second,
	}
}

// Start starts the application
func (a *App) Start() error {
	// Start server in goroutine
	go func() {
		a.logger.WithField("address", a.server.Addr).Info("Starting HTTP server")
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.logger.WithError(err).Fatal("Failed to start server")
		}
	}()

	// Wait for shutdown signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	
	<-quit
	a.logger.Info("Shutting down server...")

	return a.Shutdown()
}

// Shutdown gracefully shuts down the application
func (a *App) Shutdown() error {
	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown HTTP server
	if err := a.server.Shutdown(ctx); err != nil {
		a.logger.WithError(err).Error("Failed to shutdown server gracefully")
		return err
	}

	// Close database connection
	if err := a.db.Close(); err != nil {
		a.logger.WithError(err).Error("Failed to close database connection")
	}

	// Close Redis connection
	if err := a.redis.Close(); err != nil {
		a.logger.WithError(err).Error("Failed to close Redis connection")
	}

	a.logger.Info("Application shutdown complete")
	return nil
}