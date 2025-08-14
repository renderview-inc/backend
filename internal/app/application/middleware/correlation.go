package middleware

import (
	"context"
	"github.com/google/uuid"
	service "github.com/renderview-inc/backend/internal/app/application/services/logger"
	"net/http"
)

func CorrelationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		correlationID := r.Header.Get("Correlation-ID")
		if correlationID == "" {
			correlationID = uuid.New().String()
		}

		w.Header().Set("Correlation-ID", correlationID)
		ctx := context.WithValue(r.Context(), service.CorrelationID, correlationID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
