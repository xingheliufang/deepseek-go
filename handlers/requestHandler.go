package handlers

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

func HandleTimeout() (time.Duration, error) {
	timeoutStr := os.Getenv("DEEPSEEK_TIMEOUT")
	if timeoutStr == "" {
		return 5 * time.Minute, nil
	}
	duration, err := time.ParseDuration(timeoutStr)
	if err != nil {
		// Return error with context if parsing fails
		return 0, fmt.Errorf("invalid timeout duration %q: %w", timeoutStr, err)
	}
	return duration, nil
}

func HandleSendChatCompletionRequest(req *http.Request) (*http.Response, error) {
	timeout, err := HandleTimeout()
	if err != nil {
		return nil, fmt.Errorf("error getting timeout: %w", err)
	}
	client := &http.Client{
     Timeout: timeout,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}

	return resp, nil
}

func HandelNormalRequest(req *http.Request) (*http.Response, error) {
	timeout, err := HandleTimeout()
	if err != nil {
		return nil, fmt.Errorf("error getting timeout: %w", err)
	}

	client := &http.Client{
		Timeout: timeout,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}

	return resp, nil
}

