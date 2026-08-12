package handlers

import (
	"database/sql"
	"net/http"

	"github.com/Nehasirohi07/devSync/database"
	"github.com/Nehasirohi07/devSync/models"
	"github.com/Nehasirohi07/devSync/utils"
)

// GetAdminDashboard godoc
// @Summary Get admin dashboard
// @Description Get overall system statistics for the authenticated admin.
// @Tags Dashboard
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/admin/dashboard [get]
func GetAdminDashboard(w http.ResponseWriter, r *http.Request) {

	role, ok := r.Context().Value("role").(string)

	if !ok {
		utils.SendError(
			w,
			http.StatusUnauthorized,
			"Invalid role",
		)
		return
	}

	if role != "admin" {
		utils.SendError(
			w,
			http.StatusForbidden,
			"Only admin can access admin dashboard",
		)
		return
	}

	var dashboard models.AdminDashboard

	err := database.DB.QueryRow(
		`SELECT
			COUNT(*) AS total_users,
			COALESCE(SUM(CASE WHEN role = 'employee' THEN 1 ELSE 0 END), 0) AS total_employees,
			COALESCE(SUM(CASE WHEN role = 'manager' THEN 1 ELSE 0 END), 0) AS total_manager,
			(
				SELECT COUNT(*)
				FROM manager_requests
				WHERE status = 'pending'
			) AS pending_manager_requests,
			(
				SELECT COUNT(*)
				FROM account_deletion_requests
				WHERE status = 'pending'
			) AS pending_account_deletion_requests,
			COALESCE(SUM(CASE WHEN is_deleted = true THEN 1 ELSE 0 END), 0) AS deleted_accounts
			FROM users`,
	).Scan(
		&dashboard.TotalUsers,
		&dashboard.TotalEmployees,
		&dashboard.TotalManagers,
		&dashboard.PendingManagerRequests,
		&dashboard.PendingAccountDeletionRequests,
		&dashboard.DeletedAccounts,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			utils.SendError(
				w,
				http.StatusInternalServerError,
				"Failed to fetch admin dashboard",
			)
			return
		}
	}

	utils.SendSuccess(
		w,
		http.StatusOK,
		"Admin dashboard fetched successfully",
		dashboard,
	)
}

// GetManagerDashboard godoc
// @Summary Get manager dashboard
// @Description Get dashboard statistics for the authenticated manager.
// @Tags Dashboard
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/dashboard/manager [get]
func GetManagerDashboard(w http.ResponseWriter, r *http.Request) {

	userID, ok := r.Context().Value("userID").(int)

	if !ok {
		utils.SendError(
			w,
			http.StatusUnauthorized,
			"Invalid user",
		)
		return
	}

	role, ok := r.Context().Value("role").(string)

	if !ok {
		utils.SendError(
			w,
			http.StatusUnauthorized,
			"Invalid role",
		)
		return
	}

	if role != "manager" {
		utils.SendError(
			w,
			http.StatusForbidden,
			"Only managers can access manager dashboard",
		)
		return
	}

	var dashboard models.ManagerDashboard

	err := database.DB.QueryRow(
		`SELECT
			(
				SELECT COUNT(*)
				FROM projects
				WHERE manager_id = ?
			) AS total_projects,

			(
				SELECT COUNT(*)
				FROM tasks t
				JOIN projects p ON t.project_id = p.id
				WHERE p.manager_id = ?
			) AS total_tasks,

			(
				SELECT COUNT(*)
				FROM tasks t
				JOIN projects p ON t.project_id = p.id
				WHERE p.manager_id = ?
				AND t.status = 'completed'
			) AS completed_tasks,

			(
				SELECT COUNT(*)
				FROM tasks t
				JOIN projects p ON t.project_id = p.id
				WHERE p.manager_id = ?
				AND t.status = 'in_progress'
			) AS in_progress_tasks,

			(
				SELECT COUNT(*)
				FROM tasks t
				JOIN projects p ON t.project_id = p.id
				WHERE p.manager_id = ?
				AND t.status = 'pending'
			) AS pending_tasks,

			(
				SELECT COUNT(DISTINCT t.assigned_to)
				FROM tasks t
				JOIN projects p ON t.project_id = p.id
				WHERE p.manager_id = ?
				AND t.assigned_to IS NOT NULL
			) AS employees_assigned`,
		userID,
		userID,
		userID,
		userID,
		userID,
		userID,
	).Scan(
		&dashboard.TotalProjects,
		&dashboard.TotalTasks,
		&dashboard.CompletedTasks,
		&dashboard.InProgressTasks,
		&dashboard.PendingTasks,
		&dashboard.EmployeesAssigned,
	)

	if err != nil {
		utils.SendError(
			w,
			http.StatusInternalServerError,
			"Failed to fetch manager dashboard",
		)
		return
	}

	utils.SendSuccess(
		w,
		http.StatusOK,
		"Manager dashboard fetched successfully",
		dashboard,
	)
}

// GetEmployeeDashboard godoc
// @Summary Get employee dashboard
// @Description Get dashboard statistics for the authenticated employee.
// @Tags Dashboard
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/employee/dashboard [get]
func GetEmployeeDashboard(w http.ResponseWriter, r *http.Request) {

	userID, ok := r.Context().Value("userID").(int)

	if !ok {
		utils.SendError(
			w,
			http.StatusUnauthorized,
			"Invalid user",
		)
		return
	}

	role, ok := r.Context().Value("role").(string)

	if !ok {
		utils.SendError(
			w,
			http.StatusUnauthorized,
			"Invalid role",
		)
		return
	}

	if role != "employee" {
		utils.SendError(
			w,
			http.StatusForbidden,
			"Only employees can access employee dashboard",
		)
		return
	}

	var dashboard models.EmployeeDashboard

	err := database.DB.QueryRow(
		`SELECT
			COUNT(*) AS assigned_tasks,

			COALESCE(
				SUM(CASE
					WHEN status = 'completed' THEN 1
					ELSE 0
				END),
				0
			) AS completed_tasks,

			COALESCE(
				SUM(CASE
					WHEN status = 'in_progress' THEN 1
					ELSE 0
				END),
				0
			) AS in_progress_tasks,

			COALESCE(
				SUM(CASE
					WHEN status = 'pending' THEN 1
					ELSE 0
				END),
				0
			) AS pending_tasks,

			COALESCE(AVG(progress), 0) AS average_progress

		FROM tasks
		WHERE assigned_to = ?`,
		userID,
	).Scan(
		&dashboard.AssignedTasks,
		&dashboard.CompletedTasks,
		&dashboard.InProgressTasks,
		&dashboard.PendingTasks,
		&dashboard.AverageProgress,
	)

	if err != nil {
		utils.SendError(
			w,
			http.StatusInternalServerError,
			"Failed to fetch employee dashboard",
		)
		return
	}

	utils.SendSuccess(
		w,
		http.StatusOK,
		"Employee dashboard fetched successfully",
		dashboard,
	)
}
