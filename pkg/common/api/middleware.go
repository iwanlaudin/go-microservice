package api

import (
	"log"
	"net/http"
	"runtime/debug"
)

func ErrorLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Log the error and stack trace
				log.Printf("Panic: %v\nStack trace: %s", err, debug.Stack())

				appError := NewAppError(nil, "Internal Server Error", http.StatusInternalServerError)
				appError.SendResponse(w)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
