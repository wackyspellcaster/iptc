package auth

import (
	"net/http"
)

// These functions manage auth to set the DockerHub token in the request headers
var dockerHubToken string

func SetDockerHubToken(token string) {
	dockerHubToken = token
}

func AddAuthHeaders(req *http.Request) {
	if dockerHubToken != "" {
		req.Header.Set("Authorization", "Bearer "+dockerHubToken)
	}
}
