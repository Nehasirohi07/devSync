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
