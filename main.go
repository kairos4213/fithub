package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/kairos4213/fithub/internal/database"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	db         *database.Queries
	privateKey []byte
	publicKey  []byte
}

func main() {
	privKey, err := os.ReadFile("private_key.pem")
	if err != nil {
		log.Fatalf("missing private key: %v", err)
	}

	pubKey, err := os.ReadFile("public_key.pem")
	if err != nil {
		log.Fatalf("missing public key: %v", err)
	}

	err = godotenv.Load(".env")
	if err != nil {
		log.Printf("warning: missing or misconfigured .env: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable is not set")
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

	apiConfig := apiConfig{
		db:         dbQueries,
		privateKey: privKey,
		publicKey:  pubKey,
	}

	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))

	mux.HandleFunc("POST /api/users", apiConfig.handlerUsersCreate)
	mux.HandleFunc("GET /api/users", apiConfig.handlerUsersLogin)

	mux.HandleFunc("GET /api/healthz", handlerReadiness)

	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}
