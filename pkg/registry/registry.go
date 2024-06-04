package registry

import (
	"fmt"
	"io"
	"iptc/pkg/auth"
	"iptc/pkg/config"
	"iptc/pkg/logging"
	"net/http"
	"os"
	"path/filepath"
)

var cfg *config.Config

func SetConfig(c *config.Config) {
	cfg = c
}

func FetchManifest(imagePath string) ([]byte, error) {
	url := fmt.Sprintf("%s/v2/%s/manifests/latest", cfg.RegistryURL, imagePath)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.docker.distribution.manifest.v2+json")
	auth.AddAuthHeaders(req)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch manifest: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logging.Error(err.Error())
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch manifest: %s", resp.Status)
	}

	return io.ReadAll(resp.Body)
}

func FetchLayer(digest, imagePath string) ([]byte, error) {
	url := fmt.Sprintf("%s/v2/%s/blobs/%s", cfg.RegistryURL, imagePath, digest)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	auth.AddAuthHeaders(req)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch layer: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logging.Error(err.Error())
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch layer: %s", resp.Status)
	}

	layerPath := filepath.Join("cache", digest)
	out, err := os.Create(layerPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create layer file: %w", err)
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			logging.Error(err.Error())
		}
	}(out)

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to write layer to file: %w", err)
	}

	return os.ReadFile(layerPath)
}
