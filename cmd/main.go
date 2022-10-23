package main

import (
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/JohnDennehy101/recipe-app-backend-golang/internal/driver"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

const projectDirName = "back-end"

func loadEnv() {
	projectName := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	currentWorkDirectory, _ := os.Getwd()
	rootPath := projectName.Find([]byte(currentWorkDirectory))

	err := godotenv.Load(string(rootPath) + `/.env`)

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func main() {

	loadEnv()
	databaseHost := os.Getenv("DATABASE_HOST")
	databasePort := os.Getenv("DATABASE_PORT")
	databasePassword := os.Getenv("DATABASE_PASSWORD")
	databaseName := os.Getenv("DATABASE_NAME")
	databaseUser := os.Getenv("DATABASE_USER")

	log.Println("Connecting to database...")
	db, err := driver.ConnectSQL(fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s", databaseHost, databasePort, databaseName, databaseUser, databasePassword))

	if err != nil {
		log.Fatal("Cannot connect to database")
	}

	defer db.SQL.Close()
}
