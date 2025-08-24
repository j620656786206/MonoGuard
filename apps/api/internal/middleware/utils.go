package middleware

import (
	"time"
)

// getCurrentTimestamp returns current timestamp in ISO format
func getCurrentTimestamp() string {
	return time.Now().UTC().Format(time.RFC3339)
}