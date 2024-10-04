package helpers

import (
	"encoding/json"
	"net/http"
)

// ReadJSONRequest parses a JSON request body into the given struct.
func ReadFromRequestBody(request *http.Request, result interface{}) {
	decoder := json.NewDecoder(request.Body)
	defer request.Body.Close()
	err := decoder.Decode(result)
	PanicIfError(err)
}

// WriteJSONResponse writes a JSON response to the client.
func WriteToResponseBody(writer http.ResponseWriter, statusCode int, response interface{}) {
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(statusCode)
	err := json.NewEncoder(writer).Encode(response)
	PanicIfError(err)
}
