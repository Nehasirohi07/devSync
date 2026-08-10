package middleware

import (
	"net/http"

	"github.com/Nehasirohi07/devSync/utils"
)

func Manager(next http.Handler) http.Handler {

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

		if role != "manager" {
			utils.SendError(
				w,
				http.StatusForbidden,
				"Manager access required",
			)
			return
		}

		next.ServeHTTP(w, r)
	})
}
