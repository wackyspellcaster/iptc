package auth

import (
	"net/http"
)

var dockerHubToken string

func SetDockerHubToken(token string) {
	dockerHubToken = token
}

func AddAuthHeaders(req *http.Request) {
	if dockerHubToken != "" {
		req.Header.Set("Authorization", "Bearer "+dockerHubToken)
	}
}
