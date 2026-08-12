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

// CreateTask godoc
// @Summary Create a task
// @Description Create a new task inside a project and assign it to an employee.
// @Tags Tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param task body models.Task true "Task details"
// @Success 201 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/tasks [post]
func CreateTask(w http.ResponseWriter, r *http.Request) {

	var task models.Task

	err := json.NewDecoder(r.Body).Decode(&task)

	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	task = utils.SanitizeTask(task)

	err = utils.ValidateTask(task)

	if err != nil {
		utils.SendError(
			w,
			http.StatusBadRequest,
			err.Error(),
		)
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
			err.Error(),
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

	_, err = database.DB.Exec(
		`INSERT INTO activities (user_id, task_id, action, details)
	VALUES (?, ?, ?, ?)`,
		userID,
		task.ID,
		"task_created",
		"Task created and assigned to employee",
	)

	if err != nil {
		utils.SendError(
			w,
			http.StatusInternalServerError,
			"Task created but failed to create activity",
		)
		return
	}

	utils.SendSuccess(
		w,
		http.StatusCreated,
		"Task created successfully",
		task,
	)

}

// GetTask godoc
// @Summary Get manager tasks
// @Description Get all tasks belonging to projects managed by the authenticated manager.
// @Tags Tasks
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/tasks [get]
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
		JOIN projects p ON t.project_id = p.id
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

// GetTaskByID godoc
// @Summary Get task by ID
// @Description Get a specific task belonging to a project managed by the authenticated manager.
// @Tags Tasks
// @Produce json
// @Security BearerAuth
// @Param id path int true "Task ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/tasks/{id} [get]
func GetTaskByID(w http.ResponseWriter, r *http.Request) {

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

	var task models.Task

	err = database.DB.QueryRow(
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
		JOIN projects p ON t.Project_id = p.id
		WHERE t.id = ? AND p.manager_id = ?`,
		taskID,
		userID,
	).Scan(
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
			"Failed to fetch task",
		)
		return
	}

	utils.SendSuccess(
		w,
		http.StatusOK,
		"Task fetched successfully",
		task,
	)
}

// UpdateTask godoc
// @Summary Update a task
// @Description Update a task belonging to a project managed by the authenticated manager.
// @Tags Tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Task ID"
// @Param task body models.Task true "Updated task details"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/tasks/{id} [put]
func UpdateTask(w http.ResponseWriter, r *http.Request) {

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

	var task models.Task

	err = json.NewDecoder(r.Body).Decode(&task)

	if err != nil {
		utils.SendError(
			w,
			http.StatusBadRequest,
			"Invalid request body",
		)
		return
	}

	task = utils.SanitizeTask(task)

	err = utils.ValidateTask(task)

	if err != nil {
		utils.SendError(
			w,
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}

	result, err := database.DB.Exec(
		`UPDATE tasks t
		JOIN projects p ON t.project_id = p.id
		SET t.project_id = ?,
			t.assigned_to= ?,
			t.title = ?,
			t.description = ?,
			t.status = ?,
			t.progress = ?
		WHERE t.id = ? AND p.manager_id = ?`,
		task.ProjectID,
		task.AssignedTo,
		task.Title,
		task.Description,
		task.Status,
		task.Progress,
		taskID,
		userID,
	)

	if err != nil {
		utils.SendError(
			w,
			http.StatusInternalServerError,
			"Failed to update task",
		)
		return
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		utils.SendError(
			w,
			http.StatusInternalServerError,
			"Failed to check update task",
		)
		return
	}

	if rowsAffected == 0 {
		utils.SendError(
			w,
			http.StatusNotFound,
			"Task not found",
		)
		return
	}

	task.ID = taskID

	utils.SendSuccess(
		w,
		http.StatusOK,
		"Task updated successfully",
		task,
	)

}

// DeleteTask godoc
// @Summary Delete a task
// @Description Delete a task belonging to a project managed by the authenticated manager.
// @Tags Tasks
// @Produce json
// @Security BearerAuth
// @Param id path int true "Task ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/tasks/{id} [delete]
func DeleteTask(w http.ResponseWriter, r *http.Request) {

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

	result, err := database.DB.Exec(
		`DELETE T
		FROM tasks t
		JOIN projects p ON t.project_id = p.id
		WHERE t.id = ? AND p.manager_id = ?`,
		taskID,
		userID,
	)

	if err != nil {
		utils.SendError(
			w,
			http.StatusInternalServerError,
			"Failed to delete task",
		)
		return
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		utils.SendError(
			w,
			http.StatusInternalServerError,
			"Failed to check deleted task",
		)
		return
	}

	if rowsAffected == 0 {
		utils.SendError(
			w,
			http.StatusNotFound,
			"Task not found",
		)
		return
	}

	utils.SendSuccess(
		w,
		http.StatusOK,
		"Task deleted successfully",
		nil,
	)

}

// GetMyTask godoc
// @Summary Get my assigned tasks
// @Description Get all tasks assigned to the authenticated employee.
// @Tags Employee Tasks
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/my-tasks [get]
func GetMyTask(w http.ResponseWriter, r *http.Request) {

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
			id,
			project_id,
			assigned_to,
			title,
			description,
			status,
			progress,
			created_at
		FROM tasks
		WHERE assigned_to = ?`,
		userID,
	)

	if err != nil {
		utils.SendError(
			w,
			http.StatusInternalServerError,
			"Failed to fetch assigned tasks",
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
		"Assigned tasks fetched successfully",
		tasks,
	)
}

// UpdateMyTaskProgress godoc
// @Summary Update task progress
// @Description Update the status and progress of a task assigned to the authenticated employee.
// @Tags Employee Tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Task ID"
// @Param task body models.Task true "Task progress"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/my-tasks/{id}/progress [put]
func UpdateMyTaskProgress(w http.ResponseWriter, r *http.Request) {

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

	var task models.Task

	err = json.NewDecoder(r.Body).Decode(&task)

	if err != nil {
		utils.SendError(
			w,
			http.StatusBadRequest,
			"Invalid request body",
		)
		return
	}

	task = utils.SanitizeTask(task)

	err = utils.ValidateTaskProgress(task)

	if err != nil {
		utils.SendError(
			w,
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}

	result, err := database.DB.Exec(
		`UPDATE tasks
		SET status = ?, progress = ?
		WHERE id = ? AND assigned_to = ?`,
		task.Status,
		task.Progress,
		taskID,
		userID,
	)

	if err != nil {

		utils.SendError(
			w,
			http.StatusInternalServerError,
			"Failed to update task progress",
		)
		return
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		utils.SendError(
			w,
			http.StatusInternalServerError,
			"Failed to check update task",
		)
		return
	}

	if rowsAffected == 0 {
		utils.SendError(
			w,
			http.StatusNotFound,
			"Task not found or not assigned to you",
		)
		return
	}

	_, err = database.DB.Exec(
		`INSERT INTO activities (user_id, task_id, action, details)
	VALUES (?, ?, ?, ?)`,
		userID,
		taskID,
		"progress_updated",
		"Task progress updated to "+strconv.Itoa(task.Progress)+"%",
	)

	if err != nil {
		utils.SendError(
			w,
			http.StatusInternalServerError,
			"Task updated but failed to create activity",
		)
		return
	}

	task.ID = taskID

	utils.SendSuccess(
		w,
		http.StatusOK,
		"Task progress updated successfully",
		task,
	)
}
