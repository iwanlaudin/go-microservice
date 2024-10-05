package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/iwanlaudin/go-microservice/pkg/common/api"
)

func SetupRoutes(r chi.Router) {
	r.Use(api.ErrorLogger)

	r.Get("/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("panic")
	})

	r.Route("/auth", func(r chi.Router) {
		r.Get("/sign-in", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("This is the sign-in endpoint, no authentication required!"))
		})

		r.Get("/sign-up", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("This is the sign-in endpoint, no authentication required!"))
		})

		r.Get("/me", func(w http.ResponseWriter, r *http.Request) {
			user := map[string]interface{}{
				"id":         1,
				"first_name": "Iwan",
				"last_name":  "La Udin",
				"email":      "iwanlaudin@gmail.com",
			}

			api.NewAppResponse("successfully", http.StatusOK).Ok(w, user)
		})
	})
}
