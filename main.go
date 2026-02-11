package main

import (
	"database/sql"
	"log"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/kairos4213/fithub/internal/config"
	"github.com/kairos4213/fithub/internal/database"
	"github.com/kairos4213/fithub/internal/server"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("warning: missing or misconfigured .env: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable is not set")
	}

	filePathRoot := os.Getenv("FILEPATH_ROOT")
	if filePathRoot == "" {
		log.Fatal("FILEPATH_ROOT environment variable is not set")
	}

	tokenSecret := os.Getenv("TOKEN_SECRET")
	if tokenSecret == "" {
		log.Fatal("TOKEN_SECRET environment variable is not set")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not configured")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	dbQueries := database.New(db)
	logger := slog.Default()
	cfg := config.New(dbQueries, logger, tokenSecret)

	srv := server.New(port, filePathRoot, cfg)
	srv.Start()
}
