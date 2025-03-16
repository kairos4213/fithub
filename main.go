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
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("./app/"))))

	mux.HandleFunc("POST /api/v1/register", apiConfig.createUsersHandler)
	mux.HandleFunc("POST /api/v1/login", apiConfig.loginUsersHandler)
	mux.Handle("PUT /api/v1/users", apiConfig.authMiddleware(http.HandlerFunc(apiConfig.updateUsersHandler)))
	mux.Handle("DELETE /api/v1/users", apiConfig.authMiddleware(http.HandlerFunc(apiConfig.deleteUsersHandler)))

	mux.HandleFunc("POST /api/v1/refresh", apiConfig.refreshHandler)
	mux.HandleFunc("POST /api/v1/revoke", apiConfig.revokeHandler)

	mux.Handle("POST /api/v1/goals", apiConfig.authMiddleware(http.HandlerFunc(apiConfig.createGoalsHandler)))
	mux.Handle("GET /api/v1/goals", apiConfig.authMiddleware(http.HandlerFunc(apiConfig.getAllGoalsHandler)))
	mux.Handle("PUT /api/v1/goals/{id}", apiConfig.authMiddleware(http.HandlerFunc(apiConfig.updateGoalsHandler)))
	mux.Handle("DELETE /api/v1/goals/{id}", apiConfig.authMiddleware(http.HandlerFunc(apiConfig.deleteGoalsHandler)))
	mux.Handle("DELETE /api/v1/goals", apiConfig.authMiddleware(http.HandlerFunc(apiConfig.deleteAllGoalsHandler)))

	mux.Handle("POST /api/v1/metrics/{type}", apiConfig.authMiddleware(http.HandlerFunc(apiConfig.addMetricsHandler)))
	mux.Handle("GET /api/v1/metrics", apiConfig.authMiddleware(http.HandlerFunc(apiConfig.getAllUserMetrics)))
	mux.Handle("PUT /api/v1/metrics/{type}/{id}", apiConfig.authMiddleware(http.HandlerFunc(apiConfig.updateMetricsHandler)))
	mux.Handle("DELETE /api/v1/metrics/{type}/{id}", apiConfig.authMiddleware(http.HandlerFunc(apiConfig.deleteMetricsHandler)))
	mux.Handle("DELETE /api/v1/metrics/{type}", apiConfig.authMiddleware(http.HandlerFunc(apiConfig.deleteAllMetricsHandler)))

	mux.Handle("POST /api/v1/workouts", apiConfig.authMiddleware(http.HandlerFunc(apiConfig.createWorkoutsHandler)))
	mux.Handle("GET /api/v1/workouts", apiConfig.authMiddleware(http.HandlerFunc(apiConfig.getAllUserWorkoutsHandler)))
	mux.Handle("PUT /api/v1/workouts/{id}", apiConfig.authMiddleware(http.HandlerFunc(apiConfig.updateWorkoutsHandler)))
	mux.Handle("DELETE /api/v1/workouts/{id}", apiConfig.authMiddleware(http.HandlerFunc(apiConfig.deleteWorkoutsHandler)))
	mux.Handle("DELETE /api/v1/workouts", apiConfig.authMiddleware(http.HandlerFunc(apiConfig.deleteAllUserWorkoutsHandler)))

	mux.HandleFunc("GET /api/v1/healthz", readinessHandler)

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}
