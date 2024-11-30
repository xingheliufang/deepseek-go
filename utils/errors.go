package Utils

import (
	"errors"
	"fmt"
)

var ErrUnauthorized = errors.New("unauthorized: invalid API key")
var ErrBadRequest = errors.New("bad request: check your inputs")

type ErrorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type RequestError struct {
	HTTPStatus     string
	HTTPStatusCode int
	Err            error
	Body           []byte
}

func (e *RequestError) Error() string {
	return fmt.Sprintf(
		"Request failed with status %s (%d): %v",
		e.HTTPStatus,
		e.HTTPStatusCode,
		e.Err,
	)
}
