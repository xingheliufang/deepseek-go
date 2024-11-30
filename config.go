package deepseek

import (
	"net/http"
)

const BaseURL string = "https://api.deepseek.com/v1"

type APIType string
type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

type ClientConfig struct {
	authToken  string
	BaseURL    string
	HTTPClient HTTPDoer
}

func DefaultConfig(authToken string) ClientConfig {
	return ClientConfig{
		authToken:  authToken,
		BaseURL:    BaseURL,
		HTTPClient: &http.Client{},
	}
}
