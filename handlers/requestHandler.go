package handlers

import (
	"fmt"
	"net/http"
	"time"
)

func HandleSendChatCompletionRequest(req *http.Request) (*http.Response, error) {
	client := &http.Client{
		Timeout: 240 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	return resp, nil
}

func HandelNormalRequest(req *http.Request) (*http.Response, error) {
	client := &http.Client{
		Timeout: 20 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	return resp, nil
}
