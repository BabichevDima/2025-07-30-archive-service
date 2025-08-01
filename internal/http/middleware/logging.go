package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

func RequestLogger(logger *zap.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		logger.Info("request started",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("user_agent", r.UserAgent()),
		)

		next.ServeHTTP(w, r)

		logger.Info("request completed",
			zap.Duration("duration", time.Since(start)),
		)
	})
}
