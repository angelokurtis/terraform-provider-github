package github

import (
	"net/http"

	"github.com/google/go-github/v84/github"
)

type AuthToken string

func NewClient(httpClient *http.Client, authToken AuthToken) *github.Client {
	return github.NewClient(httpClient).WithAuthToken(string(authToken))
}
