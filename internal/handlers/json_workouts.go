package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/kairos4213/fithub/internal/cntx"
	"github.com/kairos4213/fithub/internal/database"
	"github.com/kairos4213/fithub/internal/utils"
	"github.com/kairos4213/fithub/internal/validate"
)

type Workout struct {
	ID            string `json:"id,omitempty"`
	UserID        string `json:"user_id,omitempty"`
	Title         string `json:"title,omitempty"`
	Description   string `json:"description,omitempty"`
	Duration      string `json:"duration,omitempty"`
	PlannedDate   string `json:"planned_date,omitempty"`
	DateCompleted string `json:"date_completed,omitempty"`
	CreatedAt     string `json:"created_at,omitempty"`
	UpdatedAt     string `json:"updated_at,omitempty"`
}

func (h *Handler) CreateWorkout(w http.ResponseWriter, r *http.Request) {
	userID, ok := cntx.UserID(r.Context())
	if !ok {
		utils.RespondWithError(w, http.StatusInternalServerError, "missing user id in context", nil)
		return
	}

	reqParams := Workout{}
	if err := utils.ParseJSON(r, &reqParams); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "malformed request", err)
		return
	}

	if errs := validate.Fields(
		validate.Required(reqParams.Title, "title"),
		validate.Required(reqParams.Duration, "duration"),
		validate.Required(reqParams.PlannedDate, "planned date"),
		validate.MaxLen(reqParams.Title, 100, "title"),
		validate.MaxLen(reqParams.Description, 500, "description"),
	); errs != nil {
		utils.RespondWithError(w, http.StatusBadRequest, errs[0].Error(), nil)
		return
	}

	plannedDate, err := time.Parse(time.DateOnly, reqParams.PlannedDate)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "planned date must be in YYYY-MM-DD format", err)
		return
	}

	workoutDescription := sql.NullString{Valid: false}
	if reqParams.Description != "" {
		workoutDescription.Valid = true
		workoutDescription.String = reqParams.Description
	}

	workoutDuration, err := strconv.Atoi(reqParams.Duration)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "duration must be a number", err)
		return
	}

	workout, err := h.cfg.DB.CreateWorkout(r.Context(), database.CreateWorkoutParams{
		UserID:          userID,
		Title:           reqParams.Title,
		Description:     workoutDescription,
		DurationMinutes: int32(workoutDuration),
		PlannedDate:     plannedDate,
	})
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "error creating workout", err)
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, Workout{
		ID:            workout.ID.String(),
		UserID:        workout.UserID.String(),
		Title:         workout.Title,
		Description:   workout.Description.String,
		Duration:      strconv.FormatInt(int64(workout.DurationMinutes), 10),
		PlannedDate:   workout.PlannedDate.Format(time.DateOnly),
		DateCompleted: workout.DateCompleted.Time.Format(time.DateOnly),
		CreatedAt:     workout.CreatedAt.Format(time.RFC822),
		UpdatedAt:     workout.UpdatedAt.Format(time.RFC822),
	})
}

func (h *Handler) GetAllUserWorkouts(w http.ResponseWriter, r *http.Request) {
	userID, ok := cntx.UserID(r.Context())
	if !ok {
		utils.RespondWithError(w, http.StatusInternalServerError, "missing user id in context", nil)
		return
	}

	userWorkouts, err := h.cfg.DB.GetAllUserWorkouts(r.Context(), userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "error retrieving user workouts", err)
		return
	}

	response := []Workout{}
	for _, workout := range userWorkouts {
		log.Printf("workout duration: %v", workout.DurationMinutes)
		response = append(response, Workout{
			ID:            workout.ID.String(),
			UserID:        workout.UserID.String(),
			Title:         workout.Title,
			Description:   workout.Description.String,
			Duration:      strconv.FormatInt(int64(workout.DurationMinutes), 10),
			PlannedDate:   workout.PlannedDate.Format(time.DateOnly),
			DateCompleted: workout.DateCompleted.Time.Format(time.DateOnly),
			CreatedAt:     workout.CreatedAt.Format(time.RFC822),
			UpdatedAt:     workout.UpdatedAt.Format(time.RFC822),
		})
	}
	utils.RespondWithJSON(w, http.StatusOK, response)
}

func (h *Handler) UpdateWorkout(w http.ResponseWriter, r *http.Request) {
	userID, ok := cntx.UserID(r.Context())
	if !ok {
		utils.RespondWithError(w, http.StatusInternalServerError, "missing user id in context", nil)
		return
	}
	workoutID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid workout id", err)
		return
	}

	reqParams := Workout{}
	if err := utils.ParseJSON(r, &reqParams); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "malformed request", err)
		return
	}

	if errs := validate.Fields(
		validate.Required(reqParams.Title, "title"),
		validate.Required(reqParams.Duration, "duration"),
		validate.Required(reqParams.PlannedDate, "planned date"),
		validate.MaxLen(reqParams.Title, 100, "title"),
		validate.MaxLen(reqParams.Description, 500, "description"),
	); errs != nil {
		utils.RespondWithError(w, http.StatusBadRequest, errs[0].Error(), nil)
		return
	}

	workoutDescription := sql.NullString{Valid: false}
	if reqParams.Description != "" {
		workoutDescription = sql.NullString{Valid: true, String: reqParams.Description}
	}

	workoutDuration, err := strconv.Atoi(reqParams.Duration)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "duration must be a number", err)
		return
	}

	plannedDate, err := time.Parse(time.DateOnly, reqParams.PlannedDate)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "planned date must be in YYYY-MM-DD format", err)
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

	updatedWorkout, err := h.cfg.DB.UpdateWorkout(r.Context(), database.UpdateWorkoutParams{
		Title:           reqParams.Title,
		Description:     workoutDescription,
		DurationMinutes: int32(workoutDuration),
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
		ID:            updatedWorkout.ID.String(),
		UserID:        updatedWorkout.UserID.String(),
		Title:         updatedWorkout.Title,
		Description:   updatedWorkout.Description.String,
		Duration:      strconv.FormatInt(int64(updatedWorkout.DurationMinutes), 10),
		PlannedDate:   updatedWorkout.PlannedDate.Format(time.DateOnly),
		DateCompleted: updatedWorkout.DateCompleted.Time.Format(time.DateOnly),
		CreatedAt:     updatedWorkout.CreatedAt.Format(time.RFC822),
		UpdatedAt:     updatedWorkout.UpdatedAt.Format(time.RFC822),
	})
}

func (h *Handler) DeleteWorkout(w http.ResponseWriter, r *http.Request) {
	workoutID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid workout id", err)
		return
	}
	userID, ok := cntx.UserID(r.Context())
	if !ok {
		utils.RespondWithError(w, http.StatusInternalServerError, "missing user id in context", nil)
		return
	}

	err = h.cfg.DB.DeleteWorkout(r.Context(), database.DeleteWorkoutParams{
		ID: workoutID, UserID: userID,
	})
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "error deleting workout", err)
		return
	}

	utils.RespondWithJSON(w, http.StatusNoContent, Workout{})
}

func (h *Handler) DeleteAllUserWorkouts(w http.ResponseWriter, r *http.Request) {
	userID, ok := cntx.UserID(r.Context())
	if !ok {
		utils.RespondWithError(w, http.StatusInternalServerError, "missing user id in context", nil)
		return
	}

	if err := h.cfg.DB.DeleteAllUserWorkouts(r.Context(), userID); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "error deleting user workouts", err)
		return
	}

	utils.RespondWithJSON(w, http.StatusNoContent, []Workout{})
}
