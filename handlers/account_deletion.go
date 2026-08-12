package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/Nehasirohi07/devSync/database"
	"github.com/Nehasirohi07/devSync/models"
	"github.com/Nehasirohi07/devSync/utils"
)

// CreateAccountDeletionRequest godoc
// @Summary Create account deletion request
// @Description Send a request to admin for account deletion.
// @Tags Account Deletion
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.AccountDeletionRequest true "Account deletion request"
// @Success 201 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 409 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/account-deletion-request [post]
func CreateAccountDeletionRequest(w http.ResponseWriter, r *http.Request) {

	userID, ok := r.Context().Value("userID").(int)

	if !ok {
		utils.SendError(
			w,
			http.StatusUnauthorized,
			"Invalid user",
		)
		return
	}

	var request models.AccountDeletionRequest

	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		utils.SendError(
			w,
			http.StatusBadRequest,
			"Invalid request body",
		)
		return
	}

	var existingID int

	err = database.DB.QueryRow(
		`SELECT id
		FROM account_deletion_requests
		WHERE user_id = ? AND status = 'pending'
		LIMIT 1`,
		userID,
	).Scan(&existingID)

	if err == nil {
		utils.SendError(
			w,
			http.StatusConflict,
			"Account deletion request already pending",
		)
		return
	}

	if err != sql.ErrNoRows {
		utils.SendError(
			w,
			http.StatusInternalServerError,
			"Failed to check existing request",
		)
		return
	}

	result, err := database.DB.Exec(
		`INSERT INTO account_deletion_requests
		(user_id, reason)
		VALUES (? , ?)`,
		userID,
		request.Reason,
	)

	if err != nil {
		utils.SendError(
			w,
			http.StatusInternalServerError,
			"Failed to create account deletion request",
		)
		return
	}

	requestID, err := result.LastInsertId()

	if err != nil {
		utils.SendError(
			w,
			http.StatusInternalServerError,
			"Failed to get request ID",
		)
		return
	}

	request.ID = int(requestID)
	request.UserID = userID
	request.Status = "pending"

	utils.SendSuccess(
		w,
		http.StatusCreated,
		"Account deletion request created successfully",
		request,
	)

}
