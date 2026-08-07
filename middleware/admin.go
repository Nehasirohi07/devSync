package middleware

import (
	"database/sql"
	"net/http"

	"github.com/Nehasirohi07/devSync/utils"
)

func Admin(db *sql.DB, next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		userID, ok := r.Context().Value("userID").(int)

		if !ok {
			utils.SendError(w, http.StatusUnauthorized, "User not authenticated")
			return
		}

		var role string

		err := db.QueryRow(
			`SELECT role,
			FROM users,
			WHERE user_id = ?`,
			userID,
		).Scan(&role)

		if err == sql.ErrNoRows {
			utils.SendError(w, http.StatusUnauthorized, "User not found")
			return
		}

		if err != nil {
			utils.SendError(w, http.StatusInternalServerError, "Failed to fetch user role ")
			return
		}

		if role != "manager" {
			utils.SendError(w, http.StatusForbidden, "manager access required")
			return
		}

		next.ServeHTTP(w, r)
	})
}
