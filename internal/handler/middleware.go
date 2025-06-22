package handler

import (
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

func LoggingMiddleware(logger zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			logger.Info().
				Str("method", r.Method).
				Str("url", r.URL.String()).
				Msg("started request")

			next.ServeHTTP(w, r)

			duration := time.Since(start)
			logger.Info().
				Str("method", r.Method).
				Str("url", r.URL.String()).
				Dur("duration", duration).
				Msg("completed request")
		})
	}
}
