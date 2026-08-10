package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/Nehasirohi07/devSync/config"
	"github.com/Nehasirohi07/devSync/utils"
)

func Auth(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cfg := config.LoadConfig()

		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			utils.SendError(w, http.StatusUnauthorized, "Authorization header required")
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			utils.SendError(w, http.StatusUnauthorized, "Invalid authorization format")
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := utils.ValidateJWT(tokenString, cfg.JWTSecret)

		if err != nil {
			utils.SendError(w, http.StatusUnauthorized, "Invalid or expired token")
			return
		}

		ctx := context.WithValue(r.Context(), "userID", claims.UserID)
		ctx = context.WithValue(ctx, "role", claims.Role)

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
