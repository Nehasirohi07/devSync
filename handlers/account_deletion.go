package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/Nehasirohi07/devSync/database"
	"github.com/Nehasirohi07/devSync/models"
	"github.com/Nehasirohi07/devSync/utils"
	"github.com/gorilla/mux"
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
// @Failure 403 {object} utils.Response
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

	var isDeleted bool

	err := database.DB.QueryRow(
		`SELECT is_deleted
	FROM users
	WHERE id = ?`,
		userID,
	).Scan(&isDeleted)

	if err != nil {
		if err == sql.ErrNoRows {
			utils.SendError(
				w,
				http.StatusUnauthorized,
				"User not found",
			)
			return
		}

		utils.SendError(
			w,
			http.StatusInternalServerError,
			"Failed to verify user account",
		)
		return
	}

	if isDeleted {
		utils.SendError(
			w,
			http.StatusForbidden,
			"Your account has already been deleted",
		)
		return
	}

	var request models.AccountDeletionRequest

	err = json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		utils.SendError(
			w,
			http.StatusBadRequest,
			"Invalid request body",
		)
		return
	}

	if strings.TrimSpace(request.Reason) == "" {
		utils.SendError(
			w,
			http.StatusBadRequest,
			"Deletion reason is required",
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

// GetAccountDeletionRequests godoc
// @Summary Get account deletion requests
// @Description Get all account deletion requests. Only admin can access.
// @Tags Account Deletion
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/admin/account-deletion-requests [get]
func GetAccountDeletionRequests(w http.ResponseWriter, r *http.Request) {

	rows, err := database.DB.Query(
		`SELECT
    		adr.id,
    		adr.user_id,
    		u.name,
    		u.email,
    		u.role,
    		adr.reason,
    		adr.status,
    		adr.created_at,
    		adr.reviewed_at,
    		adr.reviewed_by
		FROM account_deletion_requests adr
		JOIN users u ON adr.user_id = u.id
		ORDER BY adr.created_at DESC`,
	)

	if err != nil {
		utils.SendError(
			w,
			http.StatusInternalServerError,
			"Failed to fetch account deletion requests",
		)
		return
	}

	defer rows.Close()

	var requests []models.AccountDeletionRequest

	for rows.Next() {

		var request models.AccountDeletionRequest

		err := rows.Scan(
			&request.ID,
			&request.UserID,
			&request.UserName,
			&request.UserEmail,
			&request.UserRole,
			&request.Reason,
			&request.Status,
			&request.CreatedAt,
			&request.ReviewedAt,
			&request.ReviewedBy,
		)

		if err != nil {
			utils.SendError(
				w,
				http.StatusInternalServerError,
				"Failed to read account deletion request",
			)
			return
		}

		requests = append(requests, request)

		utils.SendSuccess(
			w,
			http.StatusOK,
			"Account deletion requests fetched successfully",
			requests,
		)

	}
}

// ApproveAccountDeletionRequest godoc
// @Summary Approve account deletion request
// @Description Approve a user's account deletion request. Only admin can access.
// @Tags Account Deletion
// @Produce json
// @Security BearerAuth
// @Param id path int true "Deletion Request ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/admin/account-deletion-requests/{id}/approve [put]
func ApproveAccountDeletionRequest(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	requestID, err := strconv.Atoi(vars["id"])

	if err != nil {
		utils.SendError(
			w,
			http.StatusBadRequest,
			"Invalid deletion request ID",
		)
		return
	}

	adminID, ok := r.Context().Value("userID").(int)

	if !ok {
		utils.SendError(
			w,
			http.StatusUnauthorized,
			"Invalid admin",
		)
		return
	}

	var userID int
	var status string

	err = database.DB.QueryRow(
		`SELECT user_id, status
		FROM account_deletion_requests
		WHERE id = ?`,
		requestID,
	).Scan(&userID, &status)

	if err != nil {

		if err == sql.ErrNoRows {
			utils.SendError(
				w,
				http.StatusNotFound,
				"Deletion request not found",
			)
			return
		}

		utils.SendError(
			w,
			http.StatusInternalServerError,
			"Failed to fetch deletion request",
		)
		return
	}

	if status != "pending" {
		utils.SendError(
			w,
			http.StatusBadRequest,
			"Deletion request has already been reviewed",
		)
		return
	}

	tx, err := database.DB.Begin()

	if err != nil {
		utils.SendError(
			w,
			http.StatusInternalServerError,
			"Failed to start transaction",
		)
		return
	}

	_, err = tx.Exec(
		`UPDATE users
		SET is_deleted = TRUE,
			deleted_at = CURRENT_TIMESTAMP
		WHERE id = ?`,
		userID,
	)

	if err != nil {
		tx.Rollback()

		utils.SendError(
			w,
			http.StatusInternalServerError,
			"Failed to delete user account",
		)
		return
	}

	_, err = tx.Exec(
		`UPDATE account_deletion_requests
		SET status = 'approved',
			reviewed_at = CURRENT_TIMESTAMP,
			reviewed_by = ?
		WHERE id = ?`,
		adminID,
		requestID,
	)

	if err != nil {
		tx.Rollback()

		utils.SendError(
			w,
			http.StatusInternalServerError,
			"Failed to update deletion request",
		)
		return
	}

	err = tx.Commit()

	if err != nil {
		utils.SendError(
			w,
			http.StatusInternalServerError,
			"Failed to complete account deletion",
		)
		return
	}

	utils.SendSuccess(
		w,
		http.StatusOK,
		"Account deletion request approved successfully",
		nil,
	)
}

// RejectAccountDeletionRequest godoc
// @Summary Reject account deletion request
// @Description Reject a user's account deletion request. Only admin can access.
// @Tags Account Deletion
// @Produce json
// @Security BearerAuth
// @Param id path int true "Deletion Request ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/admin/account-deletion-requests/{id}/reject [put]
func RejectAccountDeletionRequest(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	requestID, err := strconv.Atoi(vars["id"])

	if err != nil {
		utils.SendError(
			w,
			http.StatusBadRequest,
			"Invalid deletion request ID",
		)
		return
	}

	adminID, ok := r.Context().Value("userID").(int)

	if !ok {
		utils.SendError(
			w,
			http.StatusUnauthorized,
			"Invalid admin",
		)
		return
	}

	var status string

	err = database.DB.QueryRow(
		`SELECT status
		FROM account_deletion_requests
		WHERE id = ?`,
		requestID,
	).Scan(&status)

	if err != nil {

		if err == sql.ErrNoRows {
			utils.SendError(
				w,
				http.StatusNotFound,
				"Deletion request not found",
			)
			return
		}

		utils.SendError(
			w,
			http.StatusInternalServerError,
			"Failed to fetch deletion request",
		)
		return
	}

	if status != "pending" {
		utils.SendError(
			w,
			http.StatusBadRequest,
			"Deletion request has already been reviewed",
		)
		return
	}

	result, err := database.DB.Exec(
		`UPDATE account_deletion_requests
		SET status = 'rejected',
			reviewed_at = CURRENT_TIMESTAMP,
			reviewed_by = ?
		WHERE id = ? AND status = 'pending'`,
		adminID,
		requestID,
	)

	if err != nil {
		utils.SendError(
			w,
			http.StatusInternalServerError,
			"Failed to reject deletion request",
		)
		return
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		utils.SendError(
			w,
			http.StatusInternalServerError,
			"Failed to check deletion request",
		)
		return
	}

	if rowsAffected == 0 {
		utils.SendError(
			w,
			http.StatusNotFound,
			"Deletion request not found or already reviewed",
		)
		return
	}

	utils.SendSuccess(
		w,
		http.StatusOK,
		"Account deletion request rejected successfully",
		nil,
	)
}

// GetMyAccountDeletionRequest godoc
// @Summary Get my account deletion request
// @Description Get the latest account deletion request of the authenticated user.
// @Tags Account Deletion
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/account-deletion-request [get]
func GetMyAccountDeletionRequest(w http.ResponseWriter, r *http.Request) {

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

	err := database.DB.QueryRow(
		`SELECT
			id,
			user_id,
			reason,
			status,
			created_at,
			reviewed_at,
			reviewed_by
		FROM account_deletion_requests
		WHERE user_id = ?
		ORDER BY created_at DESC
		LIMIT 1`,
		userID,
	).Scan(
		&request.ID,
		&request.UserID,
		&request.Reason,
		&request.Status,
		&request.CreatedAt,
		&request.ReviewedAt,
		&request.ReviewedBy,
	)

	if err != nil {

		if err == sql.ErrNoRows {
			utils.SendError(
				w,
				http.StatusNotFound,
				"No account deletion request found",
			)
			return
		}

		utils.SendError(
			w,
			http.StatusInternalServerError,
			"Failed to fetch account deletion request",
		)
		return
	}

	utils.SendSuccess(
		w,
		http.StatusOK,
		"Account deletion request fetched successfully",
		request,
	)
}
