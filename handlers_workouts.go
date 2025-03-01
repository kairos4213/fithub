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
	respondWithJSON(w, http.StatusOK, response)
}

func (cfg *apiConfig) updateWorkoutsHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Title         string `json:"title"`
		Description   string `json:"description"`
		Duration      int32  `json:"duration"`
		PlannedDate   string `json:"planned_date"`
		DateCompleted string `json:"date_completed"`
	}

	userID := r.Context().Value(userIDKey).(uuid.UUID)
	workoutID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid workout id", err)
		return
	}

	reqParams := request{}
	if err := parseJSON(r, &reqParams); err != nil {
		respondWithError(w, http.StatusBadRequest, "malformed request", err)
		return
	}

	workoutDescription := sql.NullString{Valid: false}
	if reqParams.Description != "" {
		workoutDescription = sql.NullString{Valid: true, String: reqParams.Description}
	}

	plannedDate, err := time.Parse(time.DateOnly, reqParams.PlannedDate)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "incorrect date format", err)
		return
	}

	completionDate := sql.NullTime{Valid: false}
	if reqParams.DateCompleted != "" {
		dateCompleted, err := time.Parse(time.DateOnly, reqParams.DateCompleted)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "incorrect date format", err)
			return
		}
		completionDate = sql.NullTime{Valid: true, Time: dateCompleted}
	}

	updatedWorkout, err := cfg.db.UpdateWorkout(r.Context(), database.UpdateWorkoutParams{
		Title:           reqParams.Title,
		Description:     workoutDescription,
		DurationMinutes: reqParams.Duration,
		PlannedDate:     plannedDate,
		DateCompleted:   completionDate,
		ID:              workoutID,
		UserID:          userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error updating workout", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Workout{
		ID:            updatedWorkout.ID,
		UserID:        updatedWorkout.UserID,
		Title:         updatedWorkout.Title,
		Description:   updatedWorkout.Description.String,
		Duration:      updatedWorkout.DurationMinutes,
		PlannedDate:   updatedWorkout.PlannedDate,
		DateCompleted: updatedWorkout.DateCompleted.Time,
		CreatedAt:     updatedWorkout.CreatedAt,
		UpdatedAt:     updatedWorkout.UpdatedAt,
	})
}
