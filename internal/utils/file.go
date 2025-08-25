package utils

import (
	"os"
	"path/filepath"
	"strings"
)

// FileExists checks if a file exists
func FileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

// EnsureDir ensures that a directory exists
func EnsureDir(path string) error {
	return os.MkdirAll(path, 0755)
}

// GetUserHomeDir returns the user's home directory
func GetUserHomeDir() string {
	home, _ := os.UserHomeDir()
	return home
}

// JoinPaths joins multiple path components
func JoinPaths(elements ...string) string {
	return filepath.Join(elements...)
}

// CleanPath cleans a path and resolves any . or .. elements
func CleanPath(path string) string {
	return filepath.Clean(path)
}

// TrimString trims whitespace from a string
func TrimString(s string) string {
	return strings.TrimSpace(s)
}

// IsEmptyString checks if a string is empty after trimming
func IsEmptyString(s string) bool {
	return TrimString(s) == ""
}

// IsValidURL performs basic URL validation
func IsValidURL(url string) bool {
	if IsEmptyString(url) {
		return false
	}
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}

// IsValidAuthToken performs basic auth token validation
func IsValidAuthToken(token string) bool {
	if IsEmptyString(token) {
		return false
	}
	// Basic validation: should start with "sk-" prefix for Anthropic tokens
	return strings.HasPrefix(token, "sk-")
}