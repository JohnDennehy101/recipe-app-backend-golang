package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/JohnDennehy101/recipe-app-backend-golang/internal/driver"
	"github.com/JohnDennehy101/recipe-app-backend-golang/models"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

const projectDirName = "back-end"

const version = "1.0.0"

type config struct {
	port int
	env  string
}

type application struct {
	config config
	logger *log.Logger
	models models.Models
}

type AppStatus struct {
	Status      string
	Environment string
	Version     string
}

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

	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "Server port")
	flag.StringVar(&cfg.env, "env", "development", "application environment (development|production)")
	flag.Parse()

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

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	app := &application{
		config: cfg,
		logger: logger,
		models: models.NewModels(db.SQL),
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Println("Starting server on port", cfg.port)

	err = srv.ListenAndServe()

	if err != nil {
		log.Println(err)
	}

	defer db.SQL.Close()
}
