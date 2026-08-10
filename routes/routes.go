package routes

import (
	"github.com/Nehasirohi07/devSync/handlers"
	"github.com/Nehasirohi07/devSync/middleware"
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {

	router := mux.NewRouter()

	// Public routes
	router.HandleFunc("/health", handlers.HealthCheck).Methods("GET")
	router.HandleFunc("/register", handlers.Register).Methods("POST")
	router.HandleFunc("/login", handlers.Login).Methods("POST")

	// Protected routes
	protected := router.PathPrefix("/api").Subrouter()
	protected.Use(middleware.Auth)

	// Manager-only routes
	manager := protected.PathPrefix("").Subrouter()
	manager.Use(middleware.Manager)

	// Project routes - Manager only
	manager.HandleFunc("/projects", handlers.CreateProject).Methods("POST")
	manager.HandleFunc("/projects", handlers.GetProjects).Methods("GET")
	manager.HandleFunc("/projects/{id}", handlers.UpdateProject).Methods("PUT")
	manager.HandleFunc("/projects/{id}", handlers.DeleteProject).Methods("DELETE")

	// Task routes - Manager only
	manager.HandleFunc("/tasks", handlers.CreateTask).Methods("POST")
	manager.HandleFunc("/tasks", handlers.GetTask).Methods("GET")
	manager.HandleFunc("/tasks/{id}", handlers.GetTaskByID).Methods("GET")
	manager.HandleFunc("/tasks/{id}", handlers.UpdateTask).Methods("PUT")
	manager.HandleFunc("/tasks/{id}", handlers.DeleteTask).Methods("DELETE")

	protected.HandleFunc(
		"/manager-requests",
		handlers.CreateManagerRequest,
	).Methods("POST")

	// Employee routes
	protected.HandleFunc("/my-tasks", handlers.GetMyTask).Methods("GET")
	protected.HandleFunc(
		"/my-tasks/{id}/progress",
		handlers.UpdateMyTaskProgress,
	).Methods("PUT")

	protected.HandleFunc(
		"/tasks/{id}/comments",
		handlers.CreateComment,
	).Methods("POST")

	protected.HandleFunc(
		"/tasks/{id}/comments",
		handlers.GetTaskComments,
	).Methods("GET")

	protected.HandleFunc(
		"/comments/{id}",
		handlers.DeleteComment,
	).Methods("DELETE")

	protected.HandleFunc(
		"/tasks/{id}/activities",
		handlers.CreateActivity,
	).Methods("POST")

	protected.HandleFunc(
		"/tasks/{id}/activities",
		handlers.GetTaskActivities,
	).Methods("GET")

	admin := protected.PathPrefix("/admin").Subrouter()
	admin.Use(middleware.Admin)

	admin.HandleFunc(
		"/manager-requests",
		handlers.GetManagerRequests,
	).Methods("GET")

	admin.HandleFunc(
		"/manager-requests/{id}/approve",
		handlers.ApproveManagerRequest,
	).Methods("PUT")

	admin.HandleFunc(
		"/manager-requests/{id}/reject",
		handlers.RejectManagerRequest,
	).Methods("PUT")

	return router
}
