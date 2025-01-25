package deepseek

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type APIError struct {
	StatusCode    int    // HTTP status code
	APICode       int    // Business error code from API response
	Message       string // Human-readable error message
	OriginalError error  // Wrapped error for debugging
}

func (e APIError) Error() string {
	if e.APICode != 0 {
		return fmt.Sprintf("HTTP %d (Code %d): %s", e.StatusCode, e.APICode, e.Message)
	}
	return fmt.Sprintf("HTTP %d: %s", e.StatusCode, e.Message)
}

func HandleAPIError(resp *http.Response) error {
	defer func() { _ = resp.Body.Close() }()

	var apiResponse struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	body, _ := io.ReadAll(resp.Body)
	err := json.Unmarshal(body, &apiResponse)

	baseError := APIError{
		StatusCode: resp.StatusCode,
		APICode:    apiResponse.Code,
		Message:    apiResponse.Message,
	}

	if err == nil && apiResponse.Code != 0 {
		return baseError
	}

	// Handle cases couldn't parse the error response
	switch resp.StatusCode {
	case http.StatusUnauthorized:
		baseError.Message = "Invalid authentication credentials"
	case http.StatusPaymentRequired:
		baseError.Message = "Insufficient account balance"
	case http.StatusTooManyRequests:
		baseError.Message = "Rate limit exceeded"
	case http.StatusNotFound:
		baseError.Message = "Requested resource not found"
	case http.StatusInternalServerError:
		baseError.Message = "Internal server error"
	default:
		baseError.Message = fmt.Sprintf("Unexpected API response (HTTP %d)", resp.StatusCode)
	}

	baseError.OriginalError = fmt.Errorf("failed to decode: %w (body: %s)", err, string(body))
	return baseError
}
