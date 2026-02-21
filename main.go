package main

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/kairos4213/fithub/internal/config"
	"github.com/kairos4213/fithub/internal/database"
	"github.com/kairos4213/fithub/internal/server"
	_ "github.com/lib/pq"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found, reading env vars from environment")
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
	if len(tokenSecret) < 32 {
		log.Fatal("TOKEN_SECRET must be at least 32 characters")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not configured")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	dbQueries := database.New(db)
	logger := slog.Default()

	env := os.Getenv("GO_ENV")
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		if env == "production" {
			log.Fatal("BASE_URL must be set in production")
		}
		baseURL = fmt.Sprintf("http://localhost:%s", port)
	}

	oauthProviders := make(map[string]config.OAuthProvider)
	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	if googleClientID != "" && googleClientSecret != "" {
		oauthProviders["google"] = config.OAuthProvider{
			ClientID:     googleClientID,
			ClientSecret: googleClientSecret,
			RedirectURL:  baseURL + "/auth/google/callback",
		}
	} else {
		log.Println("WARNING: GOOGLE_CLIENT_ID or GOOGLE_CLIENT_SECRET not set; Google OAuth disabled")
	}

	cfg := config.New(dbQueries, db, logger, tokenSecret, oauthProviders)

	srv := server.New(port, filePathRoot, cfg, db)
	srv.Start()
}
