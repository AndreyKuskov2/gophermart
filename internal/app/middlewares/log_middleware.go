package middlewares

import (
	"net/http"
	"time"

	"github.com/AndreyKuskov2/gophermart/pkg/logger"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
)

func LoggerMiddleware(log *logger.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			next.ServeHTTP(ww, r)

			duration := time.Since(start)

			log.Log.Info("",
				zap.String("uri", r.RequestURI),
				zap.String("method", r.Method),
				zap.Duration("duration", duration),
				zap.Int("status", ww.Status()),
				zap.Int("size", ww.BytesWritten()),
			)
		})
	}
}
