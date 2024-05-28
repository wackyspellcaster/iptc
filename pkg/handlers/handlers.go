package handlers

import (
	"iptc/pkg/auth"
	"iptc/pkg/cache"
	"iptc/pkg/logging"
	"iptc/pkg/registry"
	"net/http"
)

var (
	dockerHubToken string
	imageCache     *cache.Cache
)

func SetDockerHubToken(token string) {
	auth.SetDockerHubToken(token)
}

func SetCache(c *cache.Cache) {
	imageCache = c
}

// RootHandler acts as a health check endpoint and provides the status of the server
func RootHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("Image Cache Proxy is running"))
	if err != nil {
		return
	}
}

// ProxyHandler processes requests to fetch Docker images.
// It first attempts to retrieve the requested image from the local cache.
// If the image is not found in the cache, it fetches the image from Docker Hub, caches it locally, and then serves the image to the client.
func ProxyHandler(w http.ResponseWriter, r *http.Request) {
	imagePath := r.URL.Path
	cachedImage, err := imageCache.Get(imagePath)
	if err == nil {
		logging.Info("Serving cached image: " + imagePath)
		w.Header().Set("Content-Type", "application/octet-stream")
		_, err := w.Write(cachedImage)
		if err != nil {
			logging.Error(err.Error())
		}
		return
	}

	logging.Info("Fetching image from Docker Hub: " + imagePath)
	upstreamImage, err := registry.Fetch(imagePath)
	if err != nil {
		logging.Error("Error fetching image: " + err.Error())
		http.Error(w, "Image not found", http.StatusNotFound)
		return
	}

	if err := imageCache.Set(imagePath, upstreamImage); err != nil {
		logging.Error("Error caching image: " + err.Error())
		http.Error(w, "Failed to cache image", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	_, err = w.Write(upstreamImage)
	if err != nil {
		logging.Error(err.Error())
	}
}
