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

// CreateActivity godoc
// @Summary Create task activity
// @Description Create an activity entry for a task. Only the assigned employee or project manager can create activity.
// @Tags Activities
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Task ID"
// @Param activity body models.Activity true "Activity details"
// @Success 201 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/tasks/{id}/activities [post]
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

	activity.TaskID = taskID
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

// GetTaskActivities godoc
// @Summary Get task activities
// @Description Get activity history for a task. Only the assigned employee or project manager can view it.
// @Tags Activities
// @Produce json
// @Security BearerAuth
// @Param id path int true "Task ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/tasks/{id}/activities [get]
func GetTaskActivities(w http.ResponseWriter, r *http.Request) {

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
			"You are not allowed to view activity for this task",
		)
		return
	}

	rows, err := database.DB.Query(
		`SELECT id, user_id, task_id, action , created_at
		FROM activities
		WHERE task_id  = ?
	 	ORDER BY created_at ASC`,
		taskID,
	)

	if err != nil {
		utils.SendError(
			w,
			http.StatusInternalServerError,
			"Failed to fetch activities",
		)
		return
	}

	defer rows.Close()

	var activities []models.Activity

	for rows.Next() {

		var activity models.Activity

		err := rows.Scan(
			&activity.ID,
			&activity.UserID,
			&activity.TaskID,
			&activity.Action,
			&activity.CreatedAt,
		)

		if err != nil {
			utils.SendError(
				w,
				http.StatusInternalServerError,
				"Failed to read activity data",
			)
			return
		}

		activities = append(activities, activity)

	}

	utils.SendSuccess(
		w,
		http.StatusOK,
		"Activities fetched successfully",
		activities,
	)
}
