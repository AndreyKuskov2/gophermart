package middlewares

import (
	"context"
	"net/http"

	"github.com/AndreyKuskov2/gophermart/internal/config"
	"github.com/AndreyKuskov2/gophermart/pkg/jwt"
	"github.com/AndreyKuskov2/gophermart/pkg/logger"
	"github.com/go-chi/render"
)

type contextKey string

const (
	ContextClaims contextKey = "claims"
)

func JwtAuthValidator(cfg *config.Config, log *logger.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := r.Header.Get("Authorization")
			if tokenString == "" {
				log.Log.Error("no authorization token")
				render.Status(r, http.StatusBadRequest)
				render.PlainText(w, r, "no authorization token")
				return
			}
			claims, err := jwt.VerifyToken(tokenString, cfg.JWTSecretToken)
			if err != nil {
				log.Log.Error(err.Error())
				render.Status(r, http.StatusUnauthorized)
				render.PlainText(w, r, "invalid token")
				return
			}

			r = r.Clone(context.WithValue(r.Context(), ContextClaims, claims))
			next.ServeHTTP(w, r)
		})
	}
}
