package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"standmaster/api"
	"standmaster/third_party/database"
)

func main() {
	// load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// connect to the database
	db, err := database.NewPostgres(database.PostgresConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
	})
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	} else {
		log.Println("Successfully connected to the database.")
	}
	defer db.Close()

	// create & run the API server
	address := fmt.Sprintf(":%s", os.Getenv("PORT"))
	server := api.NewAPIServer(address, db)
	if err := server.Start(); err != nil {
		log.Fatalf("Error starting the server: %v", err)
	}
}
