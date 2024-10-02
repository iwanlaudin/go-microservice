package api

import (
	"encoding/json"
	"net/http"
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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.Code)
	response := ErrorResponse{
		Status:  http.StatusText(e.Code),
		Code:    e.Code,
		Message: e.Message,
	}
	json.NewEncoder(w).Encode(response)
}

func SendResponse(w http.ResponseWriter, code int, data interface{}, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
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
	json.NewEncoder(w).Encode(response)
}

func NewAppError(err error, message string, code int) *AppError {
	return &AppError{
		Error:   err,
		Message: message,
		Code:    code,
	}
}
