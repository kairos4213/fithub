package server

import (
	"net/http"
	"time"

	"github.com/kairos4213/fithub/internal/handlers"
)

func (s *Server) registerRoutes(mux *http.ServeMux) {
	// Static files
	fileServer := http.FileServer(http.Dir(s.fileDir))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fileServer))

	// Rate limit auth endpoints: 10 requests per minute per IP
	authLimit := s.mw.RateLimit(10, time.Minute)

	// Public pages
	mux.HandleFunc("GET /", s.handler.Index)
	// FIX: /login & /register not showing up in browser network dev tool
	mux.HandleFunc("GET /login", s.handler.Login)
	mux.Handle("POST /login", authLimit(http.HandlerFunc(s.handler.Login)))
	mux.HandleFunc("POST /logout", s.handler.Logout)
	mux.HandleFunc("GET /register", s.handler.Register)
	mux.Handle("POST /register", authLimit(http.HandlerFunc(s.handler.Register)))
	mux.HandleFunc("POST /users/email", s.handler.CheckUserEmail)


	// Error pages
	mux.Handle("GET /unauthorized", http.HandlerFunc(handlers.GetUnauthorizedPage))
	mux.Handle("GET /forbidden", http.HandlerFunc(handlers.GetForbiddenPage))

	s.registerWorkoutRoutes(mux)
	s.registerExerciseRoutes(mux)
	s.registerTemplateRoutes(mux)
	s.registerMetricRoutes(mux)
	s.registerGoalRoutes(mux)
	s.registerAPIRoutes(mux)
}

func (s *Server) registerWorkoutRoutes(mux *http.ServeMux) {
	mux.Handle("GET /workouts", s.mw.Auth(http.HandlerFunc(s.handler.GetUserWorkouts)))
	mux.Handle("POST /workouts", s.mw.Auth(http.HandlerFunc(s.handler.CreateUserWorkout)))
	mux.Handle("PUT /workouts/{id}", s.mw.Auth(http.HandlerFunc(s.handler.EditUserWorkout)))
	mux.Handle("DELETE /workouts/{id}", s.mw.Auth(http.HandlerFunc(s.handler.DeleteUserWorkout)))

	mux.Handle("GET /workouts/{id}", s.mw.Auth(http.HandlerFunc(s.handler.GetUserWorkoutExercises)))
	mux.Handle("POST /workouts/{id}", s.mw.Auth(http.HandlerFunc(s.handler.AddExerciseToWorkout)))
	mux.Handle("PUT /workouts/{workoutID}/{workoutExerciseID}", s.mw.Auth(http.HandlerFunc(s.handler.UpdateWorkoutExercise)))
	mux.Handle("DELETE /workouts/{workoutID}/{workoutExerciseID}", s.mw.Auth(http.HandlerFunc(s.handler.DeleteExerciseFromWorkout)))

	mux.Handle("PUT /workouts/{id}/sort", s.mw.Auth(http.HandlerFunc(s.handler.UpdateWorkoutExercisesSortOrder)))
}

func (s *Server) registerExerciseRoutes(mux *http.ServeMux) {
	mux.Handle("GET /exercises/groups", s.mw.Auth(http.HandlerFunc(s.handler.GetExercisesPage)))
	mux.Handle("GET /exercises/groups/{group}", s.mw.Auth(http.HandlerFunc(s.handler.GetMGExercisesPage)))
	mux.Handle("GET /exercises/{id}", s.mw.Auth(http.HandlerFunc(s.handler.GetSpecificExercisePage)))
	mux.Handle("POST /exercises", s.mw.Auth(http.HandlerFunc(s.handler.GetExerciseByKeyword)))
}

func (s *Server) registerTemplateRoutes(mux *http.ServeMux) {
	mux.Handle("GET /templates", s.mw.Auth(http.HandlerFunc(s.handler.GetAllWorkoutTemplates)))
	mux.Handle("GET /templates/reroll", s.mw.Auth(http.HandlerFunc(s.handler.RerollExercise)))
	mux.Handle("GET /templates/{id}/preview", s.mw.Auth(http.HandlerFunc(s.handler.GetTemplatePreview)))
	mux.Handle("POST /templates/{id}/apply", s.mw.Auth(http.HandlerFunc(s.handler.ApplyTemplate)))
}

func (s *Server) registerMetricRoutes(mux *http.ServeMux) {
	mux.Handle("GET /metrics", s.mw.Auth(http.HandlerFunc(s.handler.GetAllMetrics)))
	mux.Handle("POST /metrics/{type}", s.mw.Auth(http.HandlerFunc(s.handler.LogMetrics)))
	mux.Handle("PUT /metrics/{type}/{id}", s.mw.Auth(http.HandlerFunc(s.handler.EditMetrics)))
	mux.Handle("DELETE /metrics/{type}/{id}", s.mw.Auth(http.HandlerFunc(s.handler.DeleteMetrics)))
}

func (s *Server) registerGoalRoutes(mux *http.ServeMux) {
	mux.Handle("GET /goals", s.mw.Auth(http.HandlerFunc(s.handler.GetAllGoals)))
	mux.Handle("POST /goals", s.mw.Auth(http.HandlerFunc(s.handler.AddNewGoal)))
	mux.Handle("PUT /goals/{id}", s.mw.Auth(http.HandlerFunc(s.handler.EditGoal)))
	mux.Handle("DELETE /goals/{id}", s.mw.Auth(http.HandlerFunc(s.handler.DeleteGoal)))
}

func (s *Server) registerAPIRoutes(mux *http.ServeMux) {
	// Auth
	mux.HandleFunc("POST /api/v1/register", s.handler.CreateUser)
	mux.HandleFunc("POST /api/v1/login", s.handler.LoginUser)
	mux.HandleFunc("POST /api/v1/refresh", s.handler.RefreshToken)
	mux.HandleFunc("POST /api/v1/revoke", s.handler.RevokeToken)

	// Users
	mux.Handle("PUT /api/v1/users", s.mw.Auth(http.HandlerFunc(s.handler.UpdateUser)))
	mux.Handle("DELETE /api/v1/users", s.mw.Auth(http.HandlerFunc(s.handler.DeleteUser)))

	// Goals
	mux.Handle("POST /api/v1/goals", s.mw.Auth(http.HandlerFunc(s.handler.CreateGoal)))
	mux.Handle("GET /api/v1/goals", s.mw.Auth(http.HandlerFunc(s.handler.GetAllUserGoals)))
	mux.Handle("PUT /api/v1/goals/{id}", s.mw.Auth(http.HandlerFunc(s.handler.UpdateGoal)))
	mux.Handle("DELETE /api/v1/goals/{id}", s.mw.Auth(http.HandlerFunc(s.handler.DeleteGoalJSON)))
	mux.Handle("DELETE /api/v1/goals", s.mw.Auth(http.HandlerFunc(s.handler.DeleteAllUserGoals)))

	// Metrics
	mux.Handle("POST /api/v1/metrics/{type}", s.mw.Auth(http.HandlerFunc(s.handler.AddMetric)))
	mux.Handle("GET /api/v1/metrics", s.mw.Auth(http.HandlerFunc(s.handler.GetAllUserMetrics)))
	mux.Handle("PUT /api/v1/metrics/{type}/{id}", s.mw.Auth(http.HandlerFunc(s.handler.UpdateMetric)))
	mux.Handle("DELETE /api/v1/metrics/{type}/{id}", s.mw.Auth(http.HandlerFunc(s.handler.DeleteMetric)))
	mux.Handle("DELETE /api/v1/metrics/{type}", s.mw.Auth(http.HandlerFunc(s.handler.DeleteAllUserMetrics)))

	// Workouts
	mux.Handle("POST /api/v1/workouts", s.mw.Auth(http.HandlerFunc(s.handler.CreateWorkout)))
	mux.Handle("GET /api/v1/workouts", s.mw.Auth(http.HandlerFunc(s.handler.GetAllUserWorkouts)))
	mux.Handle("PUT /api/v1/workouts/{id}", s.mw.Auth(http.HandlerFunc(s.handler.UpdateWorkout)))
	mux.Handle("DELETE /api/v1/workouts/{id}", s.mw.Auth(http.HandlerFunc(s.handler.DeleteWorkout)))
	mux.Handle("DELETE /api/v1/workouts", s.mw.Auth(http.HandlerFunc(s.handler.DeleteAllUserWorkouts)))

	// Health
	mux.HandleFunc("GET /api/v1/healthz", handlers.Readiness)
}
