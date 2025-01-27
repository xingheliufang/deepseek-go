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
	ResponseBody  string // Raw JSON response body
}

func (e APIError) Error() string {
	if e.APICode != 0 {
		return fmt.Sprintf("HTTP %d (Code %d): %s", e.StatusCode, e.APICode, e.Message)
	}
	return fmt.Sprintf("HTTP %d: %s \n%v", e.StatusCode, e.Message, e.ResponseBody)
}

func HandleAPIError(resp *http.Response) error {
	defer func() { _ = resp.Body.Close() }()

	body, _ := io.ReadAll(resp.Body)
	responseBody := string(body)

	var apiResponse struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	err := json.Unmarshal(body, &apiResponse)

	baseError := APIError{
		StatusCode:   resp.StatusCode,
		APICode:      apiResponse.Code,
		Message:      apiResponse.Message,
		ResponseBody: responseBody,
	}

	if err == nil && apiResponse.Code != 0 {
		return baseError
	}

	// Handle cases where the error response couldn't be parsed
	baseError.ResponseBody = responseBody
	switch resp.StatusCode {
	case http.StatusBadRequest:
		baseError.Message = "Bad request"
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

	baseError.OriginalError = fmt.Errorf("failed to decode: %w (body: %s)", err, responseBody)
	return baseError
}
