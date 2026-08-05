package routes

import (
	"github.com/Nehasirohi07/devSync/handlers"
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {

	router := mux.NewRouter()

	router.HandleFunc("/health", handlers.HealthCheck).Methods("GET")

	return router
}
