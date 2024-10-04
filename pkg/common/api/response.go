package api

import (
	"net/http"

	"github.com/iwanlaudin/go-microservice/pkg/common/helpers"
)

type AppError struct {
	Error   error  `json:"-"`
	Message string `json:"message"`
	Code    int    `json:"-"`
}

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func (e *AppError) SendResponse(w http.ResponseWriter) {
	response := ErrorResponse{
		Status:  http.StatusText(e.Code),
		Code:    e.Code,
		Message: e.Message,
	}

	helpers.WriteToResponseBody(w, e.Code, response)
}

func SendResponse(w http.ResponseWriter, code int, data interface{}, message string) {
	response := struct {
		Status  string      `json:"status"`
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data,omitempty"`
	}{
		Status:  http.StatusText(code),
		Code:    code,
		Message: message,
		Data:    data,
	}

	helpers.WriteToResponseBody(w, code, response)
}

func NewAppError(err error, message string, code int) *AppError {
	return &AppError{
		Error:   err,
		Message: message,
		Code:    code,
	}
}
