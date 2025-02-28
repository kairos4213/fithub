package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/kairos4213/fithub/internal/database"
)

type Workout struct {
	ID            uuid.UUID
	UserID        uuid.UUID
	Title         string
	Description   string
	Duration      int32
	PlannedDate   time.Time
	DateCompleted time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (cfg *apiConfig) createWorkoutsHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Duration    int32  `json:"duration"`
		PlannedDate string `json:"planned_date"`
	}
	userID := r.Context().Value(userIDKey).(uuid.UUID)

	reqParams := request{}
	if err := parseJSON(r, &reqParams); err != nil {
		respondWithError(w, http.StatusBadRequest, "malformed request", err)
		return
	}

	plannedDate, err := time.Parse(time.DateOnly, reqParams.PlannedDate)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "not a valid date", err)
		return
	}

	workoutDescription := sql.NullString{Valid: false}
	if reqParams.Description != "" {
		workoutDescription.Valid = true
		workoutDescription.String = reqParams.Description
	}

	workout, err := cfg.db.CreateWorkout(r.Context(), database.CreateWorkoutParams{
		UserID:          userID,
		Title:           reqParams.Title,
		Description:     workoutDescription,
		DurationMinutes: reqParams.Duration,
		PlannedDate:     plannedDate,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating workout", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, Workout{
		ID:            workout.ID,
		UserID:        workout.UserID,
		Title:         workout.Title,
		Description:   workout.Description.String,
		Duration:      workout.DurationMinutes,
		PlannedDate:   workout.PlannedDate,
		DateCompleted: workout.DateCompleted.Time,
		CreatedAt:     workout.CreatedAt,
		UpdatedAt:     workout.UpdatedAt,
	})
}

func (cfg *apiConfig) getAllUserWorkoutsHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIDKey).(uuid.UUID)

	userWorkouts, err := cfg.db.GetAllUserWorkouts(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error retrieving user workouts", err)
		return
	}

	response := []Workout{}
	for _, workout := range userWorkouts {
		response = append(response, Workout{
			ID:            workout.ID,
			UserID:        workout.ID,
			Title:         workout.Title,
			Description:   workout.Description.String,
			Duration:      workout.DurationMinutes,
			PlannedDate:   workout.PlannedDate,
			DateCompleted: workout.DateCompleted.Time,
			CreatedAt:     workout.CreatedAt,
			UpdatedAt:     workout.UpdatedAt,
		})
	}
	respondWithJSON(w, http.StatusOK, response)
}
