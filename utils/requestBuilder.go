package Utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type RequestBuilder struct {
	authToken string
	baseURL   string
	path      string
	body      []byte
}

// NewRequestBuilder initializes a new RequestBuilder.
func NewRequestBuilder(authToken string) *RequestBuilder {
	return &RequestBuilder{
		authToken: authToken,
	}
}

// SetBaseURL sets the base URL for the request.
func (rb *RequestBuilder) SetBaseURL(baseURL string) *RequestBuilder {
	rb.baseURL = baseURL
	return rb
}

// SetPath sets the path for the request.
func (rb *RequestBuilder) SetPath(path string) *RequestBuilder {
	rb.path = path
	return rb
}

// SetBodyFromStruct sets the request body from a struct, marshaling it to JSON.
func (rb *RequestBuilder) SetBodyFromStruct(data interface{}) *RequestBuilder {
	body, err := json.Marshal(data)
	if err != nil {
		panic(fmt.Sprintf("failed to marshal body: %v", err)) // Use panic for debugging; replace with proper error handling.
	}
	rb.body = body
	return rb
}

// Build constructs the HTTP request.
func (rb *RequestBuilder) Build(ctx context.Context) (*http.Request, error) {
	if rb.baseURL == "" || rb.path == "" {
		return nil, fmt.Errorf("baseURL or path not set")
	}

	url := fmt.Sprintf("%s%s", rb.baseURL, rb.path)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(rb.body))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+rb.authToken)
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}
