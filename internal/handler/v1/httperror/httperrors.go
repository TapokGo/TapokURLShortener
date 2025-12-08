package httperror

import (
	"net/http"
)

type HTTPError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func (e HTTPError) Error() string {
	return e.Message
}

func NotFound(message string) HTTPError {
	return HTTPError{Code: http.StatusNotFound, Message: message}
}

func InvalidRequest(message string) HTTPError {
	return HTTPError{Code: http.StatusBadRequest, Message: message}
}

func InternalServerError(message string) HTTPError {
	return HTTPError{Code: http.StatusInternalServerError, Message: message}
}
