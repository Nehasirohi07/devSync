package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"time"

	"strconv"

	"github.com/Nehasirohi07/devSync/database"
	"github.com/Nehasirohi07/devSync/models"
	"github.com/Nehasirohi07/devSync/utils"
	"github.com/gorilla/mux"
)

func CreateProject(w http.ResponseWriter, r *http.Request) {

	var project models.Project

	err := json.NewDecoder(r.Body).Decode(&project)

	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	userID, ok := r.Context().Value("userID").(int)

	if !ok {
		utils.SendError(w, http.StatusUnauthorized, "Invalid user")
		return
	}

	project.ManagerID = userID

	result, err := database.DB.Exec(
		`INSERT INTO projects (name , description, manager_id) VALUES(? , ? , ? )`,
		project.Name,
		project.Description,
		project.ManagerID,
	)

	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to create project")
		return
	}

	projectID, err := result.LastInsertId()

	project.ID = int(projectID)

	project.CreatedAt = time.Now()

	utils.SendSuccess(
		w,
		http.StatusCreated,
		"Project created successfully",
		project,
	)

}

func GetProjects(w http.ResponseWriter, r *http.Request) {

	userID, ok := r.Context().Value("userID").(int)

	if !ok {
		utils.SendError(w, http.StatusUnauthorized, "Invalid user")
		return
	}

	rows, err := database.DB.Query(
		"SELECT id, name , description, manager_id, created_at FROM projects WHERE manager_id = ?",
		userID,
	)
	if err != nil {
		utils.SendError(
			w,
			http.StatusInternalServerError,
			"Failed to fetch projects",
		)
		return
	}

	defer rows.Close()

	var projects []models.Project

	for rows.Next() {

		var project models.Project

		err := rows.Scan(
			&project.ID,
			&project.Name,
			&project.Description,
			&project.ManagerID,
			&project.CreatedAt,
		)

		if err != nil {
			utils.SendError(
				w,
				http.StatusInternalServerError,
				"Failed to read project data",
			)
			return
		}

		projects = append(projects, project)

		utils.SendSuccess(
			w,
			http.StatusOK,
			"Project fetched successfully",
			projects,
		)

	}

}

func GetProjectByID(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	projectID, err := strconv.Atoi(vars["id"])

	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid project ID")
		return
	}

	userID, ok := r.Context().Value("userID").(int)

	if !ok {
		utils.SendError(w, http.StatusUnauthorized, "Invalid user")
		return
	}

	var project models.Project

	err = database.DB.QueryRow(
		"SELECT id, name , description , manager_id , created_at FROM projects WHERE = ? AND manager_id = ?",
		projectID,
		userID,
	).Scan(
		&project.ID,
		&project.Name,
		&project.Description,
		&project.ManagerID,
		&project.CreatedAt,
	)

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

	utils.SendSuccess(
		w,
		http.StatusOK,
		"Project fetched successfully",
		project,
	)

}

func UpdateProject(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	projectID, err := strconv.Atoi(vars["id"])

	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid project ID")
		return
	}

	userID, ok := r.Context().Value("userID").(int)

	if !ok {
		utils.SendError(w, http.StatusUnauthorized, "Invalid user")
		return
	}

	var project models.Project

	err = json.NewDecoder(r.Body).Decode(&project)

	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	result, err := database.DB.Exec(
		`UPDATE projects
		SET name = ?, description = ?
		WHERE id = ? AND manager_id = ?`,
		project.Name,
		project.Description,
		projectID,
		userID,
	)

	if err != nil {
		utils.SendError(
			w,
			http.StatusInternalServerError,
			"Failed to update project",
		)
		return
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		utils.SendError(
			w,
			http.StatusInternalServerError,
			"Failed to check updated project",
		)
		return
	}

	if rowsAffected == 0 {
		utils.SendError(
			w,
			http.StatusNotFound,
			"Project not found",
		)
		return
	}

	utils.SendSuccess(
		w,
		http.StatusOK,
		"Project updated successfully",
		nil,
	)
}

func DeleteProject(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	projectID, err := strconv.Atoi(vars["id"])

	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid project ID")
		return
	}

	userID, ok := r.Context().Value("userID").(int)

	if !ok {
		utils.SendError(w, http.StatusUnauthorized, "Invalid user")
		return
	}

	result, err := database.DB.Exec(
		`DELETE FROM projects
		WHERE id = ? AND manager_id = ?`,
		projectID,
		userID,
	)

	if err != nil {
		utils.SendError(
			w,
			http.StatusInternalServerError,
			"Failed to delete project",
		)
		return
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		utils.SendError(
			w,
			http.StatusInternalServerError,
			"Failed to check deleted project",
		)
		return
	}

	if rowsAffected == 0 {
		utils.SendError(
			w,
			http.StatusNotFound,
			"Project not found",
		)
		return
	}

	utils.SendSuccess(
		w,
		http.StatusOK,
		"Project deleted successfully",
		nil,
	)

}
