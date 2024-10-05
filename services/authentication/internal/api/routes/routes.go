package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/iwanlaudin/go-microservice/pkg/common/api"
	"github.com/iwanlaudin/go-microservice/pkg/common/auth"
)

func SetupRoutes(r chi.Router) {
	r.Use(api.ErrorLogger)

	r.Get("/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("panic")
	})

	r.Get("/auth/sign-in", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("This is the sign-in endpoint, no authentication required!"))
	})

	r.Route("/auth", func(r chi.Router) {
		r.Use(auth.AuthMiddleware)

		r.Get("/me", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("hello world!"))
		})
	})
}
