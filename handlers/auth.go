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

// Register godoc
// @Summary Register a new employee
// @Description Register a new employee account. The role is automatically set to employee.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param user body models.User true "User registration details"
// @Success 201 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /register [post]
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

// Login godoc
// @Summary Login user
// @Description Login with email and password and receive a JWT token.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param user body models.User true "Login credentials"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /login [post]
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
