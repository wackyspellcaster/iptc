package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"iptc/pkg/cache"
	"iptc/pkg/config"
	"iptc/pkg/handlers"
	"iptc/pkg/logging"
	"iptc/pkg/registry"
)

func main() {
	logging.Init()
	logger := logging.GetLogger()

	cfg, err := config.Load()
	if err != nil {
		logger.Fatalf("Failed to load configuration: %v", err)
	}

	newCache, err := cache.NewCache(cfg.CacheDir, cfg.CacheSize, cfg.CacheExpiration)
	if err != nil {
		logger.Fatalf("Failed to initialize newCache: %v", err)
	}

	registry.SetConfig(cfg)
	handlers.SetCache(newCache)
	handlers.SetDockerHubToken(cfg.DockerHubToken)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      setupRouter(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Infof("Starting server on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Could not listen on port %s: %v", cfg.Port, err)
		}
	}()

	<-quit
	logger.Info("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalf("Could not gracefully shutdown the server: %v", err)
	}
	logger.Info("Server stopped")
}

func setupRouter() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.RootHandler)
	mux.HandleFunc("/v2/", handlers.ProxyHandler)
	return mux
}
