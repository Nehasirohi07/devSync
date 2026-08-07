package handlers

import (
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
