package main

// @title DevSync API
// @version 1.0
// @description Team collaboration and task management API.
// @host localhost:8081
// @BasePath /
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

import (
	"net/http"

	"log"

	"github.com/Nehasirohi07/devSync/config"
	"github.com/Nehasirohi07/devSync/database"
	"github.com/Nehasirohi07/devSync/routes"

	_ "github.com/Nehasirohi07/devSync/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {

	database.ConnectDB()

	cfg := config.LoadConfig()

	log.Printf("Server running on port %s", cfg.Port)

	router := routes.NewRouter()
	router.HandleFunc("/test-route", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ROUTE WORKING"))
	}).Methods("GET")

	router.PathPrefix("/swagger/").Handler(
		httpSwagger.WrapHandler,
	)

	err := http.ListenAndServe(":"+cfg.Port, router)

	if err != nil {
		log.Fatal("Server failed to start:", err)
	}

}
