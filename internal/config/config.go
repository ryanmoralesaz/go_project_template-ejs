package config

import (
	"os"
	"path/filepath"
)

// Config holds application configuration
type Config struct {
	Port         string
	DataFilePath string
	StaticDir    string
}

// New creates a new configuration with defaults
func New() *Config {
	// Get current working directory
	wd, _ := os.Getwd()

	return &Config{
		Port:         getEnv("PORT", "3000"),
		DataFilePath: getEnv("DATA_FILE_PATH", filepath.Join(wd, "data", "users.json")),
		StaticDir:    getEnv("STATIC_DIR", filepath.Join(wd, "static")),
	}
}

// getEnv gets an environment variable with a fallback default
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}