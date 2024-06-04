package handlers

import (
	"encoding/json"
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

func RootHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("Image Cache Proxy is running"))
	if err != nil {
		return
	}
}

type Layer struct {
	MediaType string `json:"mediaType"`
	Size      int64  `json:"size"`
	Digest    string `json:"digest"`
}

type Manifest struct {
	SchemaVersion int     `json:"schemaVersion"`
	MediaType     string  `json:"mediaType"`
	Config        Layer   `json:"config"`
	Layers        []Layer `json:"layers"`
}

func ProxyHandler(w http.ResponseWriter, r *http.Request) {
	imagePath := r.URL.Path

	// Fetch the image manifest
	manifestData, err := registry.FetchManifest(imagePath)
	if err != nil {
		logging.Error("Error fetching manifest: " + err.Error())
		http.Error(w, "Manifest not found", http.StatusNotFound)
		return
	}

	// Parse the manifest to get layer digests
	var manifest Manifest
	if err := json.Unmarshal(manifestData, &manifest); err != nil {
		logging.Error("Error parsing manifest: " + err.Error())
		http.Error(w, "Invalid manifest", http.StatusInternalServerError)
		return
	}

	// Fetch and cache each layer
	for _, layer := range manifest.Layers {
		cachedLayer, err := imageCache.Get(layer.Digest)
		if err != nil {
			logging.Info("Fetching layer from Docker Hub: " + layer.Digest)
			cachedLayer, err = registry.FetchLayer(layer.Digest, imagePath)
			if err != nil {
				logging.Error("Error fetching layer: " + err.Error())
				http.Error(w, "Layer not found", http.StatusNotFound)
				return
			}
			err := imageCache.Set(layer.Digest, cachedLayer)
			if err != nil {
				return
			}
		}

		// Serve the cached layer to the client (simplified for example)
		w.Header().Set("Content-Type", "application/octet-stream")
		_, err = w.Write(cachedLayer)
		if err != nil {
			return
		}
	}
}
