package handlers

import "strconv"

// parseIntDefault parses a string to int with a default fallback
func parseIntDefault(s string, defaultValue int) (int, error) {
	if s == "" {
		return defaultValue, nil
	}
	return strconv.Atoi(s)
}