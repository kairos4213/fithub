package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/kairos4213/fithub/internal/cntx"
	"github.com/kairos4213/fithub/internal/database"
	"github.com/kairos4213/fithub/internal/utils"
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

func (h *Handler) CreateWorkout(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Duration    int32  `json:"duration"`
		PlannedDate string `json:"planned_date"`
	}
	userID := r.Context().Value(cntx.UserIDKey).(uuid.UUID)

	reqParams := request{}
	if err := utils.ParseJSON(r, &reqParams); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "malformed request", err)
		return
	}

	plannedDate, err := time.Parse(time.DateOnly, reqParams.PlannedDate)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "not a valid date", err)
		return
	}

	workoutDescription := sql.NullString{Valid: false}
	if reqParams.Description != "" {
		workoutDescription.Valid = true
		workoutDescription.String = reqParams.Description
	}

	workout, err := h.DB.CreateWorkout(r.Context(), database.CreateWorkoutParams{
		UserID:          userID,
		Title:           reqParams.Title,
		Description:     workoutDescription,
		DurationMinutes: reqParams.Duration,
		PlannedDate:     plannedDate,
	})
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "error creating workout", err)
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, Workout{
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

func (h *Handler) GetAllUserWorkouts(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(cntx.UserIDKey).(uuid.UUID)

	userWorkouts, err := h.DB.GetAllUserWorkouts(r.Context(), userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "error retrieving user workouts", err)
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
	utils.RespondWithJSON(w, http.StatusOK, response)
}

func (h *Handler) UpdateWorkout(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Title         string `json:"title"`
		Description   string `json:"description"`
		Duration      int32  `json:"duration"`
		PlannedDate   string `json:"planned_date"`
		DateCompleted string `json:"date_completed"`
	}

	userID := r.Context().Value(cntx.UserIDKey).(uuid.UUID)
	workoutID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid workout id", err)
		return
	}

	reqParams := request{}
	if err := utils.ParseJSON(r, &reqParams); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "malformed request", err)
		return
	}

	workoutDescription := sql.NullString{Valid: false}
	if reqParams.Description != "" {
		workoutDescription = sql.NullString{Valid: true, String: reqParams.Description}
	}

	plannedDate, err := time.Parse(time.DateOnly, reqParams.PlannedDate)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "incorrect date format", err)
		return
	}

	completionDate := sql.NullTime{Valid: false}
	if reqParams.DateCompleted != "" {
		dateCompleted, err := time.Parse(time.DateOnly, reqParams.DateCompleted)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "incorrect date format", err)
			return
		}
		completionDate = sql.NullTime{Valid: true, Time: dateCompleted}
	}

	updatedWorkout, err := h.DB.UpdateWorkout(r.Context(), database.UpdateWorkoutParams{
		Title:           reqParams.Title,
		Description:     workoutDescription,
		DurationMinutes: reqParams.Duration,
		PlannedDate:     plannedDate,
		DateCompleted:   completionDate,
		ID:              workoutID,
		UserID:          userID,
	})
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "error updating workout", err)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, Workout{
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

func (h *Handler) DeleteWorkout(w http.ResponseWriter, r *http.Request) {
	workoutID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid workout id", err)
		return
	}
	userID := r.Context().Value(cntx.UserIDKey).(uuid.UUID)

	err = h.DB.DeleteWorkout(r.Context(), database.DeleteWorkoutParams{
		ID: workoutID, UserID: userID,
	})
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "error deleting workout", err)
		return
	}

	utils.RespondWithJSON(w, http.StatusNoContent, Workout{})
}

func (h *Handler) DeleteAllUserWorkouts(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(cntx.UserIDKey).(uuid.UUID)

	if err := h.DB.DeleteAllUserWorkouts(r.Context(), userID); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "error deleting user workouts", err)
		return
	}

	utils.RespondWithJSON(w, http.StatusNoContent, []Workout{})
}
