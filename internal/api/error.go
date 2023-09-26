package api

import "net/http"

type apiError struct {
	Message string      `json:"message"`
	Data    interface{} `json:"error"`
	Code    int         `json:"code"`
}

func (e apiError) Error() string {
	return e.Message
}

func newApiError(message string, data interface{}, code ...int) *apiError {
	var errorCode int
	if len(code) > 0 {
		errorCode = code[0]
	} else {
		errorCode = http.StatusBadRequest
	}
	return &apiError{Message: message, Data: data, Code: errorCode}
}
