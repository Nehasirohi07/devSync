package main

import (
	"net/http"

	"log"

	"github.com/Nehasirohi07/devSync/config"
	"github.com/Nehasirohi07/devSync/database"
	"github.com/Nehasirohi07/devSync/routes"
)

func main() {

	database.ConnectDB()

	cfg := config.LoadConfig()

	log.Printf("Server running on port %s", cfg.Port)

	router := routes.NewRouter()

	err := http.ListenAndServe(":"+cfg.Port, router)

	if err != nil {
		log.Fatal("Server failed to start:", err)
	}

}
