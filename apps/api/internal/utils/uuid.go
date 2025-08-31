package utils

import (
	"time"

	"github.com/google/uuid"
)

// GenerateUUID generates a new UUID string
func GenerateUUID() string {
	return uuid.New().String()
}

// GenerateTimestampedID generates a UUID with a timestamp if provided
func GenerateTimestampedID() (string, time.Time) {
	return uuid.New().String(), time.Now()
}