package routes

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/iwanlaudin/go-microservice/pkg/common/api"
)

func SetupRoutes(r chi.Router) {
	r.Use(api.ErrorLogger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(20 * time.Second)
		w.Write([]byte("hello world!"))
	})

	r.Get("/time-out", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(20 * time.Second)
		w.Write([]byte("hello world!"))
	})

	r.Route("/orders", func(r chi.Router) {
		r.Use(api.AuthMiddleware)

		// ... other routes
	})
}
