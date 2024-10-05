package api

import (
	"net/http"

	"github.com/iwanlaudin/go-microservice/pkg/common/helpers"
)

type ApiResponse struct {
	Message string      `json:"message"`
	Code    int         `json:"code"`
	Errors  interface{} `json:"errors"`
}

func (r *ApiResponse) Err(write http.ResponseWriter) {
	response := struct {
		Status  string `json:"status"`
		Code    int    `json:"code"`
		Message string `json:"message"`
	}{
		Status:  http.StatusText(r.Code),
		Code:    r.Code,
		Message: r.Message,
	}
	helpers.WriteToResponseBody(write, response.Code, response)
}

func (r *ApiResponse) Ok(write http.ResponseWriter, items interface{}) {
	response := struct {
		Status  string      `json:"status"`
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Items   interface{} `json:"items,omitempty"`
	}{
		Status:  http.StatusText(r.Code),
		Code:    r.Code,
		Message: r.Message,
		Items:   items,
	}
	helpers.WriteToResponseBody(write, response.Code, response)
}

func NewAppResponse(message string, code int) *ApiResponse {
	return &ApiResponse{
		Message: message,
		Code:    code,
	}
}

func NewValidationError(errors interface{}) *ApiResponse {
	return &ApiResponse{
		Message: "Invalid parameter",
		Code:    http.StatusBadRequest,
		Errors:  errors,
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
