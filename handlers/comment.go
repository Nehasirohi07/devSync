package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Nehasirohi07/devSync/database"
	"github.com/Nehasirohi07/devSync/models"
	"github.com/Nehasirohi07/devSync/utils"
	"github.com/gorilla/mux"
)

func CreateComnnet(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	taskID, err := strconv.Atoi(vars["id"])

	if err != nil {
		utils.SendError(
			w,
			http.StatusBadRequest,
			"Invalid task ID",
		)
		return
	}

	userID, ok := r.Context().Value("userID").(int)

	if !ok {
		utils.SendError(
			w,
			http.StatusUnauthorized,
			"Invalid user",
		)
		return
	}

	var comment models.Comment

	err = json.NewDecoder(r.Body).Decode(&comment)

	if err != nil {
		utils.SendError(
			w,
			http.StatusBadRequest,
			"Invalid request body",
		)
		return
	}

	comment = utils.SanitizeComment(comment)

	comment = utils.SanitizeComment(comment)

	comment.TaskID = taskID
	comment.UserID = userID
	err = utils.ValidateComment(comment)

	if err != nil {
		utils.SendError(
			w,
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}

	var taskExists int

	err = database.DB.QueryRow(
		"SELECT id FROM tasks WHERE id  = ?",
		taskID,
	).Scan(&taskExists)

	if err != nil {
		if err == sql.ErrNoRows {
			utils.SendError(
				w,
				http.StatusNotFound,
				"Task not found",
			)
			return
		}
		result, err := database.DB.Exec(
			`INSERT INTO comments (task_id, user_id, content)
			VALUES( ?, ? , ?)`,
			comment.TaskID,
			comment.UserID,
			comment.Content,
		)

		if err != nil {
			utils.SendError(
				w,
				http.StatusInternalServerError,
				"Failed to create comment",
			)
			return
		}

		commentID, err := result.LastInsertId()

		if err != nil {
			utils.SendError(
				w,
				http.StatusInternalServerError,
				"Failed to get comment ID",
			)
			return
		}

		comment.ID = int(commentID)

		utils.SendSuccess(
			w,
			http.StatusCreated,
			"Comment created successfully",
			comment,
		)

	}
}

func GetTaskComments(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	taskID, err := strconv.Atoi(vars["id"])

	if err != nil {
		utils.SendError(
			w,
			http.StatusBadRequest,
			"Invalid task ID",
		)
		return
	}

	userID, ok := r.Context().Value("userID").(int)

	if !ok {
		utils.SendError(
			w,
			http.StatusUnauthorized,
			"Invalid user",
		)
		return
	}

	var assignedTo int
	var managerID int

	err = database.DB.QueryRow(
		`SELECT t.assigned_to, p.manager_id
		FROM tasks t
		JOIN projects p ON t.project_id  = p.id
		WHERE t.id = ?`,
		taskID,
	).Scan(&assignedTo, &managerID)

	if err != nil {
		if err == sql.ErrNoRows {
			utils.SendError(
				w,
				http.StatusNotFound,
				"Task not found",
			)
			return
		}

		utils.SendError(
			w,
			http.StatusInternalServerError,
			"Failed to check task",
		)
		return
	}

	if userID != assignedTo && userID != managerID {
		utils.SendError(
			w,
			http.StatusForbidden,
			"You are not allowed to view comments for this task",
		)
		return
	}

	rows, err := database.DB.Query(
		`SELECT id, task_id, user_id, content, created_at
		FROM comments
		WHERE task_id = ?
		ORDER BY created_at ASC`,
		taskID,
	)

	if err != nil {
		utils.SendError(
			w,
			http.StatusInternalServerError,
			"Failed to fetch comments",
		)
		return
	}

	defer rows.Close()

	var comments []models.Comment

	for rows.Next() {
		var comment models.Comment

		err := rows.Scan(
			&comment.ID,
			&comment.TaskID,
			&comment.UserID,
			&comment.Content,
			&comment.CreatedAt,
		)

		if err != nil {
			utils.SendError(
				w,
				http.StatusInternalServerError,
				"Failed to read comment data",
			)
			return
		}

		comments = append(comments, comment)

	}

	utils.SendSuccess(
		w,
		http.StatusOK,
		"Comments fetched successfully",
		comments,
	)
}

func DeleteComment(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	commentID, err := strconv.Atoi(vars["id"])

	if err != nil {
		utils.SendError(
			w,
			http.StatusBadRequest,
			"Invalid comment ID",
		)
		return
	}

	userID, ok := r.Context().Value("userID").(int)

	if !ok {
		utils.SendError(
			w,
			http.StatusUnauthorized,
			"Invalid user",
		)
		return
	}

	result, err := database.DB.Exec(
		`DELETE FROM comments
		WHERE id = ? AND user_id = ?`,
		commentID,
		userID,
	)

	if err != nil {
		utils.SendError(
			w,
			http.StatusInternalServerError,
			"Failed to delete comment",
		)
		return
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		utils.SendError(
			w,
			http.StatusInternalServerError,
			"Failed to check deleted comment",
		)
		return
	}

	if rowsAffected == 0 {
		utils.SendError(
			w,
			http.StatusNotFound,
			"Comment not found or you are not allowed to delete it",
		)
		return
	}

	utils.SendSuccess(
		w,
		http.StatusOK,
		"Comment deleted successfully",
		nil,
	)
}
