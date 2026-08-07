package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/Nehasirohi07/devSync/database"
	"github.com/Nehasirohi07/devSync/models"
	"github.com/Nehasirohi07/devSync/utils"
)

func CreateTask(w http.ResponseWriter, r *http.Request) {

	var task models.Task

	err := json.NewDecoder(r.Body).Decode(&task)

	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	userID, ok := r.Context().Value("userID").(int)

	if !ok {
		utils.SendError(w, http.StatusUnauthorized, "Invalid user")
		return
	}

	var managerID int

	err = database.DB.QueryRow(
		"SELECT manager_id FROM projects WHERE id = ?",
		task.ProjectID,
	).Scan(&managerID)

	if err != nil {
		if err == sql.ErrNoRows {
			utils.SendError(
				w,
				http.StatusNotFound,
				"Project not found",
			)
			return
		}

		utils.SendError(
			w,
			http.StatusInternalServerError,
			"Failed to fetch project",
		)
		return
	}

	if managerID != userID {
		utils.SendError(
			w,
			http.StatusForbidden,
			"You are not allowed to create a task in this project",
		)
		return
	}

	result, err := database.DB.Exec(
		`INSERT INTO tasks
		(project_id, assigned_to , title , description)
		VALUES(? , ? ,?, ?)`,
		task.ProjectID,
		task.AssignedTo,
		task.Title,
		task.Description,
	)

	if err != nil {
		utils.SendError(
			w,
			http.StatusInternalServerError,
			"Failed to create task",
		)
		return
	}

	taskID, err := result.LastInsertId()

	if err != nil {
		utils.SendError(
			w,
			http.StatusInternalServerError,
			"Failed to get task ID",
		)
		return
	}

	task.ID = int(taskID)

	utils.SendSuccess(
		w,
		http.StatusCreated,
		"Task created successfully",
		task,
	)

}

func GetTask(w http.ResponseWriter, r *http.Request) {

	userID, ok := r.Context().Value("userID").(int)

	if !ok {
		utils.SendError(
			w,
			http.StatusUnauthorized,
			"Invalid user",
		)
		return
	}

	rows, err := database.DB.Query(
		`SELECT
			t.id,
			t.project_id,
			t.assigned_to,
			t.title,
			t.description,
			t.status,
			t.progress,
			t.created_at
		FROM tasks t
		JOIN project p ON t.project_id = p.id
		WHERE p.Manager_id = ?`,
		userID,
	)

	if err != nil {
		utils.SendError(
			w,
			http.StatusInternalServerError,
			"Failed to fetch tasks",
		)
		return
	}

	defer rows.Close()

	var tasks []models.Task

	for rows.Next() {
		var task models.Task

		err := rows.Scan(
			&task.ID,
			&task.ProjectID,
			&task.AssignedTo,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.Progress,
			&task.CreatedAt,
		)

		if err != nil {
			utils.SendError(
				w,
				http.StatusInternalServerError,
				"Failed to read task data",
			)
			return
		}

		tasks = append(tasks, task)
	}

	utils.SendSuccess(
		w,
		http.StatusOK,
		"Task fetched successfully",
		tasks,
	)
}
