package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/monoguard/api/internal/config"
	"github.com/monoguard/api/internal/handlers"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestHealthHandler_Live(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	
	cfg := &config.Config{
		App: config.AppConfig{
			Name:    "test-service",
			Version: "1.0.0",
		},
	}
	
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel) // Suppress logs during testing
	
	handler := handlers.NewHealthHandler(nil, nil, cfg, logger)
	
	router := gin.New()
	router.GET("/live", handler.Live)

	// Execute
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/live", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "alive")
	assert.Contains(t, w.Body.String(), "test-service")
	assert.Contains(t, w.Body.String(), "1.0.0")
}