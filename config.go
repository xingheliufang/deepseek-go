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
	AuthToken  string
	BaseURL    string
	HTTPClient HTTPDoer
}

func DefaultConfig(AuthToken string) ClientConfig {
	return ClientConfig{
		AuthToken:  AuthToken,
		BaseURL:    BaseURL,
		HTTPClient: &http.Client{},
	}
}
