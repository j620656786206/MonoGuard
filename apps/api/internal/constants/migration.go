package constants

import "sync"

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