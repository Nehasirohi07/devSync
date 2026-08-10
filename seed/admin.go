package main

import (
	"log"
	"os"

	"github.com/Nehasirohi07/devSync/database"
	"github.com/Nehasirohi07/devSync/utils"
)

func main() {

	database.ConnectDB()

	defer database.DB.Close()

	adminName := os.Getenv("ADMIN_NAME")
	adminEmail := os.Getenv("ADMIN_EMAIL")
	adminPassword := os.Getenv("ADMIN_PASSWORD")

	if adminName == "" || adminEmail == "" || adminPassword == "" {
		log.Fatal("Admin environment varaiables are missing")
	}

	var existingID int

	err := database.DB.QueryRow(
		`SELECT id 
		FROM users
		WHERE email = ?`,
		adminEmail,
	).Scan(&existingID)

	if err == nil {
		log.Printf(
			"Admin already exists with ID: %d",
			existingID,
		)
		return
	}

	hashedpassword, err := utils.HashPassword(adminPassword)

	if err != nil {
		log.Fatal("Failed to hash admin password:", err)
	}

	_, err = database.DB.Exec(
		`INSERT INTO users (name, email, password, role)
		VALUES (? , ? ,? ,'admin')`,
		adminName,
		adminEmail,
		hashedpassword,
	)

	if err != nil {
		log.Fatal("Failed to craete admin:", err)
	}

	log.Println("Admin created successfully")

}
