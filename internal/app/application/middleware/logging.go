package middleware

import (
	service "github.com/renderview-inc/backend/internal/app/application/services/logger"
	"github.com/renderview-inc/backend/internal/app/application/services/logger/option"
	"net/http"
)

func LoggingMiddleware(next http.Handler, logService *service.LogService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		correlationID := r.Context().Value(service.CorrelationID)
		logService.Info(r.Context(), "incoming request",
			option.Any("method", r.Method),
			option.Any("url", r.URL.Path),
			option.Any("cid", correlationID),
		)
		next.ServeHTTP(w, r)
	})
}
