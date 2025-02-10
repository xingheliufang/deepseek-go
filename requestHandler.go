package deepseek

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// HandleTimeout gets the timeout duration from the DEEPSEEK_TIMEOUT environment variable.
func HandleTimeout() (time.Duration, error) {
	if err := godotenv.Load(); err != nil {
		_ = err
	}

	timeoutStr := os.Getenv("DEEPSEEK_TIMEOUT")
	if timeoutStr == "" {
		return 5 * time.Minute, nil
	}
	duration, err := time.ParseDuration(timeoutStr)
	if err != nil {
		return 0, fmt.Errorf("invalid timeout duration %q: %w", timeoutStr, err)
	}
	return duration, nil
}

// checkTimeoutError checks if the error is a timeout error and returns a custom error message.
func checkTimeoutError(err error, timeout time.Duration) error {
	var urlErr *url.Error
	if errors.As(err, &urlErr) && urlErr.Timeout() {
		return fmt.Errorf(
			"request timed out after %s. You can increase the timeout by setting the DEEPSEEK_TIMEOUT environment variable. Original error: %w",
			timeout,
			err,
		)
	}
	return nil
}

// HandleSendChatCompletionRequest sends a request to the DeepSeek API and returns the response.
func HandleSendChatCompletionRequest(c Client, req *http.Request) (*http.Response, error) {
	// Check if c.Timeout is already set or not
	timeout := c.Timeout
	if timeout == 0 {
		var err error
		timeout, err = HandleTimeout()
		if err != nil {
			return nil, fmt.Errorf("error getting timeout: %w", err)
		}
	}
	client := &http.Client{Timeout: timeout}
	resp, err := client.Do(req)
	if err != nil {
		if timeoutErr := checkTimeoutError(err, timeout); timeoutErr != nil {
			return nil, timeoutErr
		}
		return nil, fmt.Errorf("error sending request: %w", err)
	}

	return resp, nil
}

// HandleNormalRequest sends a request to the DeepSeek API and returns the response.
func HandleNormalRequest(c Client, req *http.Request) (*http.Response, error) {
	// Check if c.Timeout is already set or not
	timeout := c.Timeout
	if timeout == 0 {
		var err error
		timeout, err = HandleTimeout()
		if err != nil {
			return nil, fmt.Errorf("error getting timeout: %w", err)
		}
	}
	client := &http.Client{Timeout: timeout}
	resp, err := client.Do(req)
	if err != nil {
		if timeoutErr := checkTimeoutError(err, timeout); timeoutErr != nil {
			return nil, timeoutErr
		}
		return nil, fmt.Errorf("error sending request: %w", err)
	}

	return resp, nil
}
