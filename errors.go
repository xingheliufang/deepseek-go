package deepseek

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ErrorResponse represents the error response from the API.
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Error implements the error interface.
func (e ErrorResponse) Error() string {
	return fmt.Sprintf("API error: %d - %s", e.Code, e.Message)
}

// APIError represents a generic API error.
type APIError struct {
	StatusCode int
	ErrorMsg   string
}

// Error implements the error interface.
func (a APIError) Error() string {
	return fmt.Sprintf("API error %d: %s", a.StatusCode, a.ErrorMsg)
}

// UnauthorizedError represents a 401 Unauthorized error.
type UnauthorizedError struct {
	APIError
}

// NotFoundError represents a 404 Not Found error.
type NotFoundError struct {
	APIError
}

// ServerErrorCode represents internal server errors.
type ServerErrorCode struct {
	APIError
}

func HandleAPIError(resp *http.Response) error {
	var apiErr ErrorResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiErr); err != nil {
		switch resp.StatusCode {
		case http.StatusUnauthorized:
			return UnauthorizedError{APIError{StatusCode: resp.StatusCode, ErrorMsg: "401 Error are a result of unauthorized acess. Please make sure your API key is correct."}}
		case http.StatusNotFound:
			return NotFoundError{APIError{StatusCode: resp.StatusCode, ErrorMsg: apiErr.Message}}
		case http.StatusInternalServerError:
			return ServerErrorCode{APIError{StatusCode: resp.StatusCode, ErrorMsg: apiErr.Message}}
		default:
			return APIError{StatusCode: resp.StatusCode, ErrorMsg: ("Failed to decode error respons " + apiErr.Message)}
		}
	}
	return APIError{StatusCode: resp.StatusCode, ErrorMsg: "Failed to decode error response"}
}
