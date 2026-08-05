package main

import (
	"net/http"

	"log"

	"github.com/Nehasirohi07/devSync/config"
	"github.com/Nehasirohi07/devSync/database"
)

func main() {

	database.ConnectDB()

	cfg := config.LoadConfig()

	log.Printf("Server running on port %s", cfg.Port)

	err := http.ListenAndServe(":"+cfg.Port, nil)

	if err != nil {
		log.Fatal("Server failed to start:", err)
	}

}
