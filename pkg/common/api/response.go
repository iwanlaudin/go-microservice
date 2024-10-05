package api

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/iwanlaudin/go-microservice/pkg/common/helpers"
)

type ApiResponse struct {
	Message string
	Code    int
}

func NewAppResponse(message string, code int) *ApiResponse {
	return &ApiResponse{
		Message: message,
		Code:    code,
	}
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

func (r *ApiResponse) ValidationErr(write http.ResponseWriter, err error) {
	validationErrors := validationErrors(err)

	response := struct {
		Status  string      `json:"status"`
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Errors  interface{} `json:"errors,omitempty"`
	}{
		Status:  http.StatusText(r.Code),
		Code:    r.Code,
		Message: r.Message,
		Errors:  validationErrors,
	}
	helpers.WriteToResponseBody(write, response.Code, response)
}

func validationErrors(err error) []map[string]string {
	var validationErrors []map[string]string

	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range errs {
			validationErrors = append(validationErrors, map[string]string{
				"field":   fieldError.Field(),
				"tag":     fieldError.Tag(),
				"message": validationErrorMessage(fieldError),
			})
		}
	}
	return validationErrors
}

func validationErrorMessage(vErr validator.FieldError) string {
	switch vErr.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return "Value is too short"
	// Add more validation tag messages as needed
	default:
		return "Invalid value"
	}
}
