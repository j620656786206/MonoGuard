package database

import (
	"context"
	"time"
)

// createTimeoutContext creates a context with timeout
func createTimeoutContext(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}