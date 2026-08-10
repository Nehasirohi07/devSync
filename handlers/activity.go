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

func CreateActivity(w http.ResponseWriter, r *http.Request) {

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

	var activity models.Activity

	err = json.NewDecoder(r.Body).Decode(&activity)

	if err != nil {
		utils.SendError(
			w,
			http.StatusBadRequest,
			"Invalid request body",
		)
		return
	}

	activity = utils.SanitizeActivity(activity)

	activity.TaskID = activity.TaskID
	activity.UserID = userID

	err = utils.ValidateActivity(activity)

	if err != nil {
		utils.SendError(
			w,
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}

	var assignedTo int
	var managerID int

	err = database.DB.QueryRow(
		`SELECT t.assigned_to, p.manager_id
		FROM tasks t
		JOIN projects p ON t.project_id = p.id
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
			"You are not allowed to create activity for this task",
		)
		return
	}

	result, err := database.DB.Exec(
		`INSERT INTO activities (user_id, task_id, action)
		VALUES (? , ? , ?)`,
		activity.UserID,
		activity.TaskID,
		activity.Action,
	)

	if err != nil {
		utils.SendError(
			w,
			http.StatusInternalServerError,
			"Failed to create activity",
		)
		return
	}

	activityID, err := result.LastInsertId()

	if err != nil {
		utils.SendError(
			w,
			http.StatusInternalServerError,
			"Failed to get activity ID",
		)
		return
	}

	activity.ID = int(activityID)

	utils.SendSuccess(
		w,
		http.StatusCreated,
		"Activity created successfully",
		activity,
	)

}
