package database

import (
	"database/sql"
	"fmt"
	"log"

	"time"

	"github.com/Nehasirohi07/devSync/config"
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func ConnectDB() {

	cfg := config.LoadConfig()

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",

		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		log.Fatalf("Failed to open database:%v", err)
	}

	maxRetries := 5

	for i := 1; i <= maxRetries; i++ {

		err = db.Ping()

		if err == nil {
			break
		}

		if i == maxRetries {
			log.Fatal("Failed to connect to database after retries:", err)
		}

		log.Printf("Database not ready... Retry %d/%d", i, maxRetries)

		time.Sleep(2 * time.Second)

	}

	DB = db

}
