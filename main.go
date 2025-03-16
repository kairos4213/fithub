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

type api struct {
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

	api := api{
		db:         dbQueries,
		privateKey: privKey,
		publicKey:  pubKey,
	}

	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("./app/"))))

	mux.HandleFunc("POST /api/v1/register", api.createUsersHandler)
	mux.HandleFunc("POST /api/v1/login", api.loginUsersHandler)
	mux.Handle("PUT /api/v1/users", api.authMiddleware(http.HandlerFunc(api.updateUsersHandler)))
	mux.Handle("DELETE /api/v1/users", api.authMiddleware(http.HandlerFunc(api.deleteUsersHandler)))

	mux.HandleFunc("POST /api/v1/refresh", api.refreshHandler)
	mux.HandleFunc("POST /api/v1/revoke", api.revokeHandler)

	mux.Handle("POST /api/v1/goals", api.authMiddleware(http.HandlerFunc(api.createGoalsHandler)))
	mux.Handle("GET /api/v1/goals", api.authMiddleware(http.HandlerFunc(api.getAllGoalsHandler)))
	mux.Handle("PUT /api/v1/goals/{id}", api.authMiddleware(http.HandlerFunc(api.updateGoalsHandler)))
	mux.Handle("DELETE /api/v1/goals/{id}", api.authMiddleware(http.HandlerFunc(api.deleteGoalsHandler)))
	mux.Handle("DELETE /api/v1/goals", api.authMiddleware(http.HandlerFunc(api.deleteAllGoalsHandler)))

	mux.Handle("POST /api/v1/metrics/{type}", api.authMiddleware(http.HandlerFunc(api.addMetricsHandler)))
	mux.Handle("GET /api/v1/metrics", api.authMiddleware(http.HandlerFunc(api.getAllUserMetrics)))
	mux.Handle("PUT /api/v1/metrics/{type}/{id}", api.authMiddleware(http.HandlerFunc(api.updateMetricsHandler)))
	mux.Handle("DELETE /api/v1/metrics/{type}/{id}", api.authMiddleware(http.HandlerFunc(api.deleteMetricsHandler)))
	mux.Handle("DELETE /api/v1/metrics/{type}", api.authMiddleware(http.HandlerFunc(api.deleteAllMetricsHandler)))

	mux.Handle("POST /api/v1/workouts", api.authMiddleware(http.HandlerFunc(api.createWorkoutsHandler)))
	mux.Handle("GET /api/v1/workouts", api.authMiddleware(http.HandlerFunc(api.getAllUserWorkoutsHandler)))
	mux.Handle("PUT /api/v1/workouts/{id}", api.authMiddleware(http.HandlerFunc(api.updateWorkoutsHandler)))
	mux.Handle("DELETE /api/v1/workouts/{id}", api.authMiddleware(http.HandlerFunc(api.deleteWorkoutsHandler)))
	mux.Handle("DELETE /api/v1/workouts", api.authMiddleware(http.HandlerFunc(api.deleteAllUserWorkoutsHandler)))

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
