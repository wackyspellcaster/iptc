package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port            string
	DockerHubToken  string
	CacheDir        string
	CacheSize       int
	CacheExpiration time.Duration
}

// Load collects configuration values from environment variables and K8s secrets
func Load() (*Config, error) {
	// Default to port 8080 if not set
	port := getEnv("PORT", "8080")

	// Read dockerhub token from k8s secret
	token, err := os.ReadFile("/etc/secrets/DOCKER_HUB_TOKEN")
	if err != nil {
		return nil, err
	}

	// Get cache config values from envvars
	cacheDir := getEnv("CACHE_DIR", "cache")
	cacheSize, err := strconv.Atoi(getEnv("CACHE_SIZE", "100"))
	if err != nil {
		cacheSize = 100
	}
	cacheExpiration, err := time.ParseDuration(getEnv("CACHE_EXPIRATION", "24h"))
	if err != nil {
		cacheExpiration = 24 * time.Hour
	}

	return &Config{
		Port:            port,
		DockerHubToken:  string(token),
		CacheDir:        cacheDir,
		CacheSize:       cacheSize,
		CacheExpiration: cacheExpiration,
	}, nil
}

// getEnv returns an environment variable value from a given argument if found
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
