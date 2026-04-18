package github

import (
	"net/http"

	"github.com/google/go-github/v84/github"
)

type Token string

func NewClient(httpClient *http.Client, token Token) *github.Client {
	return github.NewClient(httpClient).WithAuthToken(string(token))
}
