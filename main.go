package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/kairos4213/fithub/internal/database"
	"github.com/kairos4213/fithub/internal/handlers"
	"github.com/kairos4213/fithub/internal/middleware"
	_ "github.com/lib/pq"
)

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

	filePathRoot := os.Getenv("FILEPATH_ROOT")
	if filePathRoot == "" {
		log.Fatal("FILEPATH_ROOT environment variable is not set")
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

	mw := middleware.Middleware{PublicKey: pubKey}

	handler := handlers.Handler{DB: dbQueries, PrivateKey: privKey}

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir(filePathRoot))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fileServer))

	mux.HandleFunc("GET /", handler.Index)

	mux.HandleFunc("GET /login", handler.Login)
	mux.HandleFunc("POST /login", handler.Login)

	mux.HandleFunc("POST /logout", handler.Logout)

	mux.HandleFunc("GET /register", handler.Register)
	mux.HandleFunc("POST /register", handler.Register)
	mux.HandleFunc("POST /users/email", handler.CheckUserEmail)

	mux.Handle("GET /workouts", mw.Auth(http.HandlerFunc(handler.GetUserWorkouts)))
	mux.Handle("POST /workouts", mw.Auth(http.HandlerFunc(handler.CreateUserWorkout)))
	mux.Handle("PUT /workouts/{id}", mw.Auth(http.HandlerFunc(handler.EditUserWorkout)))
	mux.Handle("DELETE /workouts/{id}", mw.Auth(http.HandlerFunc(handler.DeleteUserWorkout)))

	mux.Handle("GET /workouts/{id}", mw.Auth(http.HandlerFunc(handler.GetUserWorkoutExercises)))
	mux.Handle("POST /workouts/{id}", mw.Auth(http.HandlerFunc(handler.AddExerciseToWorkout)))
	mux.Handle("PUT /workouts/{workoutID}/{workoutExerciseID}", mw.Auth(http.HandlerFunc(handler.UpdateWorkoutExercise)))
	mux.Handle("DELETE /workouts/{workoutID}/{workoutExerciseID}", mw.Auth(http.HandlerFunc(handler.DeleteExerciseFromWorkout)))

	mux.Handle("PUT /workouts/{id}/sort", mw.Auth(http.HandlerFunc(handler.UpdateWorkoutExercisesSortOrder)))

	mux.Handle("GET /exercises", mw.Auth(http.HandlerFunc(handler.GetExercisesPage)))
	mux.Handle("POST /exercises", mw.Auth(http.HandlerFunc(handler.GetExerciseByKeyword)))

	mux.Handle("GET /metrics", mw.Auth(http.HandlerFunc(handler.GetAllMetrics)))
	mux.Handle("POST /metrics/{type}", mw.Auth(http.HandlerFunc(handler.LogMetrics)))
	mux.Handle("PUT /metrics/{type}/{id}", mw.Auth(http.HandlerFunc(handler.EditMetrics)))
	mux.Handle("DELETE /metrics/{type}/{id}", mw.Auth(http.HandlerFunc(handler.DeleteMetrics)))

	mux.Handle("GET /goals", mw.Auth(http.HandlerFunc(handler.GetAllGoals)))
	mux.Handle("POST /goals", mw.Auth(http.HandlerFunc(handler.AddNewGoal)))
	mux.Handle("PUT /goals/{id}", mw.Auth(http.HandlerFunc(handler.EditGoal)))
	mux.Handle("DELETE /goals/{id}", mw.Auth(http.HandlerFunc(handler.DeleteGoal)))

	mux.Handle("GET /unauthorized", http.HandlerFunc(handlers.GetUnauthorizedPage))
	mux.Handle("GET /forbidden", http.HandlerFunc(handlers.GetForbiddenPage))

	mux.Handle("GET /admin", mw.AdminAuth(http.HandlerFunc(handler.GetAdminHome)))

	mux.Handle("GET /admin/exercises", mw.AdminAuth(http.HandlerFunc(handler.GetAdminExercisesPage)))
	mux.Handle("POST /admin/exercises", mw.AdminAuth(http.HandlerFunc(handler.AddDBExercise)))
	mux.Handle("PUT /admin/exercises/{id}", mw.AdminAuth(http.HandlerFunc(handler.EditDBExercise)))
	mux.Handle("DELETE /admin/exercises/{id}", mw.AdminAuth(http.HandlerFunc(handler.DeleteDBExercise)))

	mux.HandleFunc("POST /api/v1/register", handler.CreateUser)
	mux.HandleFunc("POST /api/v1/login", handler.LoginUser)
	mux.Handle("PUT /api/v1/users", mw.Auth(http.HandlerFunc(handler.UpdateUser)))
	mux.Handle("DELETE /api/v1/users", mw.Auth(http.HandlerFunc(handler.DeleteUser)))

	mux.HandleFunc("POST /api/v1/refresh", handler.RefreshToken)
	mux.HandleFunc("POST /api/v1/revoke", handler.RevokeToken)

	mux.Handle("POST /api/v1/goals", mw.Auth(http.HandlerFunc(handler.CreateGoal)))
	mux.Handle("GET /api/v1/goals", mw.Auth(http.HandlerFunc(handler.GetAllUserGoals)))
	mux.Handle("PUT /api/v1/goals/{id}", mw.Auth(http.HandlerFunc(handler.UpdateGoal)))
	mux.Handle("DELETE /api/v1/goals/{id}", mw.Auth(http.HandlerFunc(handler.DeleteGoalJSON)))
	mux.Handle("DELETE /api/v1/goals", mw.Auth(http.HandlerFunc(handler.DeleteAllUserGoals)))

	mux.Handle("POST /api/v1/metrics/{type}", mw.Auth(http.HandlerFunc(handler.AddMetric)))
	mux.Handle("GET /api/v1/metrics", mw.Auth(http.HandlerFunc(handler.GetAllUserMetrics)))
	mux.Handle("PUT /api/v1/metrics/{type}/{id}", mw.Auth(http.HandlerFunc(handler.UpdateMetric)))
	mux.Handle("DELETE /api/v1/metrics/{type}/{id}", mw.Auth(http.HandlerFunc(handler.DeleteMetric)))
	mux.Handle("DELETE /api/v1/metrics/{type}", mw.Auth(http.HandlerFunc(handler.DeleteAllUserMetrics)))

	mux.Handle("POST /api/v1/workouts", mw.Auth(http.HandlerFunc(handler.CreateWorkout)))
	mux.Handle("GET /api/v1/workouts", mw.Auth(http.HandlerFunc(handler.GetAllUserWorkouts)))
	mux.Handle("PUT /api/v1/workouts/{id}", mw.Auth(http.HandlerFunc(handler.UpdateWorkout)))
	mux.Handle("DELETE /api/v1/workouts/{id}", mw.Auth(http.HandlerFunc(handler.DeleteWorkout)))
	mux.Handle("DELETE /api/v1/workouts", mw.Auth(http.HandlerFunc(handler.DeleteAllUserWorkouts)))

	mux.HandleFunc("GET /api/v1/healthz", handlers.Readiness)

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mw.Log(mux),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}
