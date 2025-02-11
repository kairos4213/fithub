package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

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

	mux.HandleFunc("POST /api/users", apiConfig.createUsersHandler)
	mux.HandleFunc("GET /api/users", apiConfig.loginUsersHandler)
	mux.Handle("PUT /api/users", apiConfig.authMiddleware(http.HandlerFunc(apiConfig.updateUsersHandler)))
	// TODO: getUsersHandler
	// TODO: deleteUsersHandler

	mux.HandleFunc("POST /api/refresh", apiConfig.refreshHandler)
	mux.HandleFunc("POST /api/revoke", apiConfig.revokeHandler)

	mux.Handle("POST /api/goals", apiConfig.authMiddleware(http.HandlerFunc(apiConfig.createGoalsHandler)))
	mux.Handle("GET /api/goals", apiConfig.authMiddleware(http.HandlerFunc(apiConfig.getAllGoalsHandler)))
	mux.Handle("PUT /api/goals", apiConfig.authMiddleware(http.HandlerFunc(apiConfig.updateGoalsHandler)))
	mux.Handle("DELETE /api/goals/{id}", apiConfig.authMiddleware(http.HandlerFunc(apiConfig.deleteGoalsHandler)))
	// TODO: deleteAllGoalsHandler
	// TODO: getGoalsHandler

	mux.HandleFunc("GET /api/healthz", readinessHandler)

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}
