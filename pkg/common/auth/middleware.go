package auth

import (
	"net/http"
	"strings"

	"github.com/iwanlaudin/go-microservice/pkg/common/api"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			api.NewAppError("Missing auth token", http.StatusUnauthorized).SendResponse(w)
			return
		}

		bearerToken := strings.Fields(authHeader)
		if len(bearerToken) < 2 || bearerToken[0] != "Bearer" {
			api.NewAppError("Invalid token format", http.StatusUnauthorized).SendResponse(w)
			return
		}

		user, err := ValidateToken(bearerToken[1])
		if err != nil {
			api.NewAppError(err.Error(), http.StatusUnauthorized).SendResponse(w)
			return
		}

		ctx := ContextWithUser(r.Context(), user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
