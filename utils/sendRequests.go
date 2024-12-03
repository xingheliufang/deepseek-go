package Utils

import (
	"fmt"
	"net/http"
	"time"
)

func SendRequest(req *http.Request) (*http.Response, error) {
	client := &http.Client{
		Timeout: 120 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	return resp, nil
}
