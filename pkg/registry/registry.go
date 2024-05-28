package registry

import (
	"fmt"
	"io"
	"iptc/pkg/logging"
	"net/http"

	"iptc/pkg/auth"
)

// Fetch retrieves a Docker image from DockerHub
func Fetch(imagePath string) ([]byte, error) {
	url := "https://registry-1.docker.io" + imagePath
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	auth.AddAuthHeaders(req)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch image: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logging.Error(err.Error())
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch image: %s", resp.Status)
	}

	return io.ReadAll(resp.Body)
}
