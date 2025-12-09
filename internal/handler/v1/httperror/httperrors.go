// Package httperror provides HTTPError struct methods for returning errors
package httperror

import (
	"net/http"
)

// HTTPError is a model of HTTP error
type HTTPError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// NotFound return StatusNotFound
func NotFound(message string) HTTPError {
	return HTTPError{Code: http.StatusNotFound, Message: message}
}

// InvalidRequest return StatusBadRequest
func InvalidRequest(message string) HTTPError {
	return HTTPError{Code: http.StatusBadRequest, Message: message}
}

// InternalServerError return StatusInternalServerError
func InternalServerError(message string) HTTPError {
	return HTTPError{Code: http.StatusInternalServerError, Message: message}
}
