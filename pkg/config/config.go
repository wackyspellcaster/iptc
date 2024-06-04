package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port            string
	DockerHubToken  string
	RegistryURL     string
	CacheDir        string
	CacheSize       int
	CacheExpiration time.Duration
}

func Load() (*Config, error) {

	port := getEnv("PORT", "8080")

	token, err := os.ReadFile("/etc/secrets/DOCKER_HUB_TOKEN")
	if err != nil {
		return nil, err
	}

	registryURL := getEnv("REGISTRY_URL", "https://registry-1.docker.io")
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
		RegistryURL:     registryURL,
		CacheDir:        cacheDir,
		CacheSize:       cacheSize,
		CacheExpiration: cacheExpiration,
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
