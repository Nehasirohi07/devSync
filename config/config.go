package config

import (
	"log" //message or error printing
	"os"  //use to read environment varaiables from operating system

	"github.com/joho/godotenv" //external package for load .env file
)

// it is a custom data structure that stores all database settings in one place
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

// this is a complete configuration for application
type Config struct {
	Database  DatabaseConfig
	JWTSecret string
	Port      string
}

// this function loads configuration
func LoadConfig() Config {

	err := godotenv.Load() //load .env file

	if err != nil {
		log.Println("Warning: env file not found")
	}

	//since the function's return type Config , we need to return a Config object
	return Config{
		Database: DatabaseConfig{ // here , we are filling the database field of config, which is of type databaseconfig
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
		},

		JWTSecret: os.Getenv("JWT_SECRET"),
		Port:      os.Getenv("PORT"),
	}

}
