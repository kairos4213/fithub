package handlers

import (
	"database/sql"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kairos4213/fithub/internal/cntx"
	"github.com/kairos4213/fithub/internal/database"
	"github.com/kairos4213/fithub/internal/templates"
)

func (h *Handler) GetUserWorkouts(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(cntx.UserIDKey).(uuid.UUID)
	workouts, err := h.cfg.DB.GetAllUserWorkouts(r.Context(), userID)
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to get user workouts", slog.String("error", err.Error()))
		return
	}

	w.Header().Set("Content-Type", "text/html")
	contents := templates.Workouts(workouts)
	err = templates.Layout(contents, "Fithub | Workouts", true).Render(r.Context(), w)
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to render user workouts", slog.String("error", err.Error()))
		return
	}
}

func (h *Handler) CreateUserWorkout(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(cntx.UserIDKey).(uuid.UUID)

	reqTitle := r.FormValue("title")
	reqDescription := r.FormValue("workout-description")
	reqDuration := r.FormValue("duration")
	reqPlannedDate := r.FormValue("planned-date")

	if reqTitle == "" || reqDuration == "" || reqPlannedDate == "" {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Info("unable to create workout: missing required form info", slog.String("user", userID.String()))
		return
	}

	description := sql.NullString{Valid: false}
	if reqDescription != "" {
		description.String = reqDescription
		description.Valid = true
	}

	duration, err := strconv.ParseInt(reqDuration, 10, 32)
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to parse duration", slog.String("error", err.Error()))
		return
	}

	plannedDate, err := time.Parse(time.DateOnly, reqPlannedDate)
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to parse planned date", slog.String("error", err.Error()))
		return
	}

	_, err = h.cfg.DB.CreateWorkout(r.Context(), database.CreateWorkoutParams{
		UserID: userID, Title: reqTitle, Description: description, DurationMinutes: int32(duration), PlannedDate: plannedDate,
	})
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to create workout", slog.String("error", err.Error()))
		return
	}

	workouts, err := h.cfg.DB.GetAllUserWorkouts(r.Context(), userID)
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to fetch all workouts", slog.String("error", err.Error()))
		return
	}

	err = templates.WorkoutsTableBody(workouts).Render(r.Context(), w)
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to render workouts table body", slog.String("error", err.Error()))
		return
	}
}

func (h *Handler) EditUserWorkout(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(cntx.UserIDKey).(uuid.UUID)
	workoutID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to parse workout id", slog.String("error", err.Error()))
		return
	}

	reqTitle := r.FormValue("title")
	reqDescription := r.FormValue("workout-description")
	reqDuration := r.FormValue("duration")
	reqPlannedDate := r.FormValue("planned-date")
	reqCompletionDate := r.FormValue("date-completed")

	description := sql.NullString{Valid: false}
	if reqDescription != "" {
		description.String = reqDescription
		description.Valid = true
	}

	duration, err := strconv.ParseInt(reqDuration, 10, 32)
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to parse duration", slog.String("error", err.Error()))
		return
	}

	plannedDate, err := time.Parse(time.DateOnly, reqPlannedDate)
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to parse planned date", slog.String("error", err.Error()))
		return
	}

	dateCompleted := sql.NullTime{Valid: false}
	if reqCompletionDate != "" {
		date, err := time.Parse(time.DateOnly, reqCompletionDate)
		if err != nil {
			HandleInternalServerError(w, r)
			h.cfg.Logger.Error("failed to parse date completed", slog.String("error", err.Error()))
			return
		}
		dateCompleted.Time = date
		dateCompleted.Valid = true
	}

	updatedWorkout, err := h.cfg.DB.UpdateWorkout(r.Context(), database.UpdateWorkoutParams{
		Title:           reqTitle,
		Description:     description,
		DurationMinutes: int32(duration),
		PlannedDate:     plannedDate,
		DateCompleted:   dateCompleted,
		ID:              workoutID,
		UserID:          userID,
	})
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to update workout", slog.String("error", err.Error()))
		return
	}

	currentURL := r.Header.Get("HX-Current-URL")
	if strings.Contains(currentURL, "/workouts/"+workoutID.String()) {
		err = templates.WorkoutInfo(updatedWorkout).Render(r.Context(), w)
		if err != nil {
			HandleInternalServerError(w, r)
			h.cfg.Logger.Error("failed to render workout info", slog.String("error", err.Error()))
			return
		}
		return
	}

	workouts, err := h.cfg.DB.GetAllUserWorkouts(r.Context(), userID)
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to get all user workouts", slog.String("error", err.Error()))
		return
	}

	err = templates.WorkoutsTableBody(workouts).Render(r.Context(), w)
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to render workouts table body", slog.String("error", err.Error()))
		return
	}
}

func (h *Handler) DeleteUserWorkout(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(cntx.UserIDKey).(uuid.UUID)
	workoutID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to parse workout id", slog.String("error", err.Error()))
		return
	}

	err = h.cfg.DB.DeleteWorkout(r.Context(), database.DeleteWorkoutParams{
		ID:     workoutID,
		UserID: userID,
	})
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to delete workout", slog.String("error", err.Error()))
		return
	}

	currentURL := r.Header.Get("HX-Current-URL")
	if strings.Contains(currentURL, "/workouts/"+workoutID.String()) {
		w.Header().Set("HX-Location", `{ "path": "/workouts" }`)
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetUserWorkoutExercises(w http.ResponseWriter, r *http.Request) {
	workoutID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to parse workout id", slog.String("error", err.Error()))
		return
	}

	workout, err := h.cfg.DB.GetWorkoutByID(r.Context(), workoutID)
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to fetch workout", slog.String("error", err.Error()))
		return
	}

	workoutExercises, err := h.cfg.DB.WorkoutAndExercises(r.Context(), workoutID)
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to get workout exercises", slog.String("error", err.Error()))
		return
	}

	contents := templates.WorkoutPage(workout, workoutExercises, []database.Exercise{})
	err = templates.Layout(contents, "FitHub | Workout Page", true).Render(r.Context(), w)
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to render workout page", slog.String("error", err.Error()))
		return
	}
}
