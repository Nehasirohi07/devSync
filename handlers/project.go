package handlers

import (
	"encoding/json"
	"net/http"

	"time"

	"github.com/Nehasirohi07/devSync/database"
	"github.com/Nehasirohi07/devSync/models"
	"github.com/Nehasirohi07/devSync/utils"
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
