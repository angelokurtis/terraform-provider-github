package http

import (
	"net/http"
)

func NewClient() *http.Client {
	c := &http.Client{
		Transport: new(RequestLogger),
	}

	return c
}
