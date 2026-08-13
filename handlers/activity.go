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
		`SELECT
		a.id,
		a.user_id,
		u.name,
		a.task_id,
		a.action,
		a.details,
		a.created_at
	FROM activities a
	JOIN users u ON a.user_id = u.id
	WHERE a.task_id = ?
	ORDER BY a.created_at ASC`,
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
			&activity.UserName,
			&activity.TaskID,
			&activity.Action,
			&activity.Details,
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
