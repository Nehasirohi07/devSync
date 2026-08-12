package models

type AdminDashboard struct {
	TotalUsers                     int `json:"total_users"`
	TotalEmployees                 int `json:"total_employees"`
	TotalManagers                  int `json:"total_managers"`
	PendingManagerRequests         int `json:"pending_manager_requests"`
	PendingAccountDeletionRequests int `json:"pending_account_deletion_requests"`
	DeletedAccounts                int `json:"deleted_accounts"`
}

type ManagerDashboard struct {
	TotalProjects     int `json:"total_projects"`
	TotalTasks        int `json:"total_tasks"`
	CompletedTasks    int `json:"completed_tasks"`
	InProgressTasks   int `json:"in_progress_tasks"`
	PendingTasks      int `json:"pending_tasks"`
	EmployeesAssigned int `json:"employees_assigned"`
}

type EmployeeDashboard struct {
	AssignedTasks   int     `json:"assigned_tasks"`
	CompletedTasks  int     `json:"completed_tasks"`
	InProgressTasks int     `json:"in_progress_tasks"`
	PendingTasks    int     `json:"pending_tasks"`
	AverageProgress float64 `json:"average_progress"`
}
