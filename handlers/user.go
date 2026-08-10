package handlers

import (
	"database/sql"
	"net/http"

	"github.com/Nehasirohi07/devSync/database"
	"github.com/Nehasirohi07/devSync/models"
	"github.com/Nehasirohi07/devSync/utils"
)

// GetMyProfile godoc
// @Summary Get my profile
// @Description Get details of the currently authenticated user.
// @Tags Users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/me [get]
func GetMyProfile(w http.ResponseWriter, r *http.Request) {

	userID, ok := r.Context().Value("userID").(int)

	if !ok {
		utils.SendError(
			w,
			http.StatusUnauthorized,
			"Invalid user",
		)
		return
	}

	var user models.User

	err := database.DB.QueryRow(
		`SELECT id, name , email , role
		FROM users
		WHERE id = ?`,
		userID,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Role,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			utils.SendError(
				w,
				http.StatusNotFound,
				"User not found",
			)
			return
		}

		utils.SendError(
			w,
			http.StatusInternalServerError,
			"Failed to fetch user profiles",
		)
		return
	}

	utils.SendSuccess(
		w,
		http.StatusOK,
		"User profile fetched successfully",
		user,
	)
}

// GetUsers godoc
// @Summary Get users
// @Description Get employees for managers and employees plus managers for admins.
// @Tags Users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/users [get]
func GetUsers(w http.ResponseWriter, r *http.Request) {

	role, ok := r.Context().Value("role").(string)

	if !ok {
		utils.SendError(
			w,
			http.StatusUnauthorized,
			"User role not found",
		)
		return
	}

	var rows *sql.Rows
	var err error

	if role == "manager" {

		rows, err = database.DB.Query(
			`SELECT id, name, email, role
			FROM users
			WHERE role = 'employee'
			ORDER BY id ASC`,
		)
	} else if role == "admin" {
		rows, err = database.DB.Query(
			`SELECT id, name, email, role
			FROM users
			WHERE role IN ('employee', 'manager')
			ORDER BY id ASC`,
		)
	} else {
		utils.SendError(
			w,
			http.StatusForbidden,
			"Manager or admin access required",
		)
		return
	}
	if err != nil {
		utils.SendError(
			w,
			http.StatusInternalServerError,
			"Failed to fetch users",
		)
		return
	}

	defer rows.Close()

	var users []models.User

	for rows.Next() {

		var user models.User

		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Role,
		)

		if err != nil {
			utils.SendError(
				w,
				http.StatusInternalServerError,
				"Failed to read user data",
			)
			return
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		utils.SendError(
			w,
			http.StatusInternalServerError,
			"Failed to read users",
		)
		return
	}

	utils.SendSuccess(
		w,
		http.StatusOK,
		"Users fetched successfully",
		users,
	)

}
