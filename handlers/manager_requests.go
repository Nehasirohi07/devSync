package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/Nehasirohi07/devSync/database"
	"github.com/Nehasirohi07/devSync/models"
	"github.com/Nehasirohi07/devSync/utils"
	"github.com/gorilla/mux"
)

func CreateManagerRequest(w http.ResponseWriter, r http.Request) {

	userID, ok := r.Context().Value("userID").(int)

	if !ok {
		utils.SendError(
			w,
			http.StatusUnauthorized,
			"Invalid user",
		)
		return
	}

	var existingRequest models.ManagerRequest

	err := database.DB.QueryRow(
		`SELECT id, user_id, status, created_at
		FROM manager_requests
		WHERE user_id = ? AND status = 'pending'`,
		userID,
	).Scan(
		&existingRequest.ID,
		&existingRequest.UserID,
		&existingRequest.Status,
		&existingRequest.CreatedAt,
	)

	if err == nil {
		utils.SendError(
			w,
			http.StatusConflict,
			"Manager request already pending",
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
		`INSERT INTO manager_requests (user_id, status)
		VALUES(?, 'pending')`,
		userID,
	)

	if err != nil {
		utils.SendError(
			w,
			http.StatusInternalServerError,
			"Failed to create manager request",
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

	request := models.ManagerRequest{
		ID:     int(requestID),
		UserID: userID,
		Status: "pending",
	}

	utils.SendSuccess(
		w,
		http.StatusCreated,
		"Manager request submitted successfully",
		request,
	)
}

func GetManagerRequests(w http.ResponseWriter, r *http.Request) {

	rows, err := database.DB.Query(
		`SELECT id, user_id, status, created_at
		FROM manager_requests
		WHERE status = 'pending'
		ORDER BY created_at ASC`,
	)

	if err != nil {
		utils.SendError(
			w,
			http.StatusInternalServerError,
			"Failed to fetch manager requests",
		)
		return
	}

	defer rows.Close()

	var requests []models.ManagerRequest

	for rows.Next() {

		var request models.ManagerRequest

		err := rows.Scan(
			&request.ID,
			&request.UserID,
			&request.Status,
			&request.CreatedAt,
		)

		if err != nil {
			utils.SendError(
				w,
				http.StatusInternalServerError,
				"Failed to read manager request",
			)
			return
		}

		requests = append(requests, request)
	}

	if err := rows.Err(); err != nil {
		utils.SendError(
			w,
			http.StatusInternalServerError,
			"Failed to read manager requests",
		)
		return
	}

	utils.SendSuccess(
		w,
		http.StatusOK,
		"Manager requests fetched successfully",
		requests,
	)

}

func ApproveManagerRequest(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	requestID, err := strconv.Atoi(vars["id"])

	if err != nil {
		utils.SendError(
			w,
			http.StatusBadRequest,
			"Invalid manager request ID",
		)
		return
	}

	var userID int

	err = database.DB.QueryRow(
		`SELECT user_id
		FROM manager_requests
		WHERE id = ? AND status = 'pending'`,
		requestID,
	).Scan(&userID)

	if err != nil {
		if err == sql.ErrNoRows {
			utils.SendError(
				w,
				http.StatusNotFound,
				"Pending manager request not found",
			)
			return
		}

		utils.SendError(
			w,
			http.StatusInternalServerError,
			"Failed to fetch manager request",
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
		SET role = 'manager'
		WHERE id = ?`,
		userID,
	)

	if err != nil {
		tx.Rollback()

		utils.SendError(
			w,
			http.StatusInternalServerError,
			"Failed to update user role",
		)
		return
	}

	_, err = tx.Exec(
		`UPDATE manager_requests
				SET status = 'approved'
				WHERE id = ?`,
		requestID,
	)

	if err != nil {
		tx.Rollback()

		utils.SendError(
			w,
			http.StatusInternalServerError,
			"Failed to approve manager request",
		)
		return
	}

	err = tx.Commit()

	if err != nil {
		utils.SendError(
			w,
			http.StatusInternalServerError,
			"Failed to complete manager approval",
		)
		return
	}

	utils.SendSuccess(
		w,
		http.StatusOK,
		"Manager request approved successfully",
		map[string]interface{}{
			"request_id": requestID,
			"user_id":    userID,
			"role":       "manager",
		},
	)
}

func RejectManagerRequest(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	requestID, err := strconv.Atoi(vars["id"])

	if err != nil {
		utils.SendError(
			w,
			http.StatusBadRequest,
			"Invalid manager request ID",
		)
		return
	}

	result, err := database.DB.Exec(
		`UPDATE manager_requests
		SET status = 'rejected'
		WHERE id = ? AND status = 'pending'`,
		requestID,
	)

	if err != nil {
		utils.SendError(
			w,
			http.StatusInternalServerError,
			"Failed to reject manager request",
		)
		return
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		utils.SendError(
			w,
			http.StatusInternalServerError,
			"Failed to check rejected manager request",
		)
		return
	}

	if rowsAffected == 0 {
		utils.SendError(
			w,
			http.StatusNotFound,
			"Pending manager request not found",
		)
		return
	}

	utils.SendSuccess(
		w,
		http.StatusOK,
		"Manager request rejected successfully",
		map[string]interface{}{
			"request_id": requestID,
			"status":     "rejected",
		},
	)
}
