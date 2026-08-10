package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/Nehasirohi07/devSync/config"
	"github.com/Nehasirohi07/devSync/database"
	"github.com/Nehasirohi07/devSync/models"
	"github.com/Nehasirohi07/devSync/utils"
)

func Register(w http.ResponseWriter, r *http.Request) {

	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user = utils.SanitizeUser(user)

	err = utils.ValidateUser(user)

	if err != nil {
		utils.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	user.Role = "employee"

	hashedPassword, err := utils.HashPassword(user.Password)

	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	user.Password = hashedPassword

	result, err := database.DB.Exec(
		"INSERT INTO users(name, email, password, role) VALUES(? , ? , ?, ?)",
		user.Name,
		user.Email,
		user.Password,
		user.Role,
	)

	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	userID, err := result.LastInsertId()

	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to get user ID ")
		return
	}

	userIDInt := int(userID)

	cfg := config.LoadConfig()

	token, err := utils.GenerateJWT(
		userIDInt,
		user.Role,
		cfg.JWTSecret,
	)

	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	utils.SendSuccess(
		w,
		http.StatusCreated,
		"User registered successfully",
		map[string]interface{}{
			"token": token,
		},
	)

}

func Login(w http.ResponseWriter, r *http.Request) {

	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user = utils.SanitizeUser(user)

	err = utils.ValidateLoginUser(user)

	if err != nil {
		utils.SendError(
			w,
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}

	var StoredUser models.User

	err = database.DB.QueryRow(
		"SELECT id, name , email , password , role FROM users WHERE email = ?",
		user.Email,
	).Scan(
		&StoredUser.ID,
		&StoredUser.Name,
		&StoredUser.Email,
		&StoredUser.Password,
		&StoredUser.Role,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			utils.SendError(w, http.StatusInternalServerError, "Invalid email or password")
			return
		}

		utils.SendError(w, http.StatusInternalServerError, "Database error")
		return
	}

	passwordMatch := utils.ComparePassword(
		StoredUser.Password,
		user.Password,
	)

	if !passwordMatch {
		utils.SendError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	cfg := config.LoadConfig()

	token, err := utils.GenerateJWT(
		StoredUser.ID,
		StoredUser.Role,
		cfg.JWTSecret,
	)

	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	utils.SendSuccess(
		w,
		http.StatusOK,
		"Login successfully",
		map[string]interface{}{
			"token": token,
		},
	)
}
