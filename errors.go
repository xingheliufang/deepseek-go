package deepseek

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e ErrorResponse) Error() string {
	return fmt.Sprintf("API error: %d - %s", e.Code, e.Message)
}

type APIError struct {
	StatusCode int
	ErrorMsg   string
}

// @implements APIError
func (a APIError) Error() string {
	return fmt.Sprintf("API error %d: %s", a.StatusCode, a.ErrorMsg)
}

// Tries to handle errors listed on: https://api-docs.deepseek.com/quick_start/error_codes
func HandleAPIError(resp *http.Response) error {
	var apiErr ErrorResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiErr); err != nil {
		switch resp.StatusCode {
		case http.StatusUnauthorized:
			return APIError{StatusCode: resp.StatusCode, ErrorMsg: "401 Error are a result of unauthorized acess. Please make sure your API key is correct."}
		case http.StatusPaymentRequired:
			return APIError{StatusCode: resp.StatusCode, ErrorMsg: "You have run out of balance. Please check your account's balance, and go to the Top up page to add funds. https://platform.deepseek.com/top_up"}
		case http.StatusTooManyRequests:
			return APIError{StatusCode: resp.StatusCode, ErrorMsg: "You are sending requests too quickly."}
		case http.StatusNotFound:
			return APIError{StatusCode: resp.StatusCode, ErrorMsg: apiErr.Message}
		case http.StatusInternalServerError:
			return APIError{StatusCode: resp.StatusCode, ErrorMsg: apiErr.Message}
		default:
			return APIError{StatusCode: resp.StatusCode, ErrorMsg: ("Failed to decode error respons " + apiErr.Message)}
		}
	}
	return APIError{StatusCode: resp.StatusCode, ErrorMsg: "Failed to decode error response"}
}
