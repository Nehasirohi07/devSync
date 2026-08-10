package middleware

import (
	"net/http"

	"github.com/Nehasirohi07/devSync/utils"
)

func Admin(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		role, ok := r.Context().Value("role").(string)

		if !ok {
			utils.SendError(
				w,
				http.StatusUnauthorized,
				"User role not found",
			)
			return
		}

		if role != "admin" {
			utils.SendError(
				w,
				http.StatusForbidden,
				"Admin access required",
			)
			return
		}

		next.ServeHTTP(w, r)
	})
}
