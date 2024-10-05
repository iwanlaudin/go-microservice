package api

import (
	"net/http"

	"github.com/iwanlaudin/go-microservice/pkg/common/helpers"
)

type AppError struct {
	Status  string      `json:"-"`
	Message string      `json:"message"`
	Code    int         `json:"-"`
	Errors  interface{} `json:"errors"`
}

func (e *AppError) SendResponse(write http.ResponseWriter) {
	response := struct {
		Status  string      `json:"status"`
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Errors  interface{} `json:"data,omitempty"`
	}{
		Status:  http.StatusText(e.Code),
		Code:    e.Code,
		Message: e.Message,
		Errors:  e.Errors,
	}

	helpers.WriteToResponseBody(write, e.Code, response)
}

func SendResponse(write http.ResponseWriter, code int, items interface{}, message string) {
	response := struct {
		Status  string      `json:"status"`
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Items   interface{} `json:"data,omitempty"`
	}{
		Status:  http.StatusText(code),
		Code:    code,
		Message: message,
		Items:   items,
	}

	helpers.WriteToResponseBody(write, code, response)
}

func NewAppError(message string, code int) *AppError {
	return &AppError{
		Status:  http.StatusText(code),
		Message: message,
		Code:    code,
	}
}

// func NewAppErrorWithValidation(err error, message string, code int) *AppError {
// 	var validationErrors []map[string]string

// 	if errs, ok := err.(validator.ValidationErrors); ok {
// 		for _, fieldError := range errs {
// 			validationErrors = append(validationErrors, map[string]string{
// 				"field":   fieldError.Field(),
// 				"tag":     fieldError.Tag(),
// 				"message": fieldError.Error(),
// 			})
// 		}
// 	}

// 	return &AppError{
// 		Status:   http.StatusText(code),
// 		Message: message,
// 		Code:    code,
// 		Errors:  validationErrors,
// 	}
// }
