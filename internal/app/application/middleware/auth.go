package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/renderview-inc/backend/internal/app/application/services"
)

func AuthMiddleware(next http.Handler, authService *services.AuthService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bearerToken := r.Header.Get("Authorization")
		if bearerToken == "" || !strings.HasPrefix(bearerToken, "Bearer ") {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(bearerToken, "Bearer ")

		ctx, cancel := context.WithTimeout(r.Context(), 300*time.Millisecond)
		defer cancel()

		err := authService.Authorize(ctx, token)
		if err != nil {
			if err == services.ErrAccessTokenInvalid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		next.ServeHTTP(w, r)
	})
}
