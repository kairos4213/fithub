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
	"github.com/kairos4213/fithub/internal/validate"
)

func (h *Handler) GetUserWorkouts(w http.ResponseWriter, r *http.Request) {
	userID, ok := cntx.UserID(r.Context())
	if !ok {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("missing user id in context")
		return
	}

	// Exercises page carousel — keep existing behavior
	currentURL := r.Header.Get("HX-Current-URL")
	target := r.Header.Get("HX-Target")
	if strings.Contains(currentURL, "/exercises/") && target == "user-workouts" {
		workouts, err := h.cfg.DB.GetAllUserWorkouts(r.Context(), userID)
		if err != nil {
			HandleInternalServerError(w, r)
			h.cfg.Logger.Error("failed to get user workouts", slog.String("error", err.Error()))
			return
		}
		err = templates.UserWorkoutsHTML(workouts).Render(r.Context(), w)
		if err != nil {
			HandleInternalServerError(w, r)
			h.cfg.Logger.Error("failed to render user workouts html", slog.String("error", err.Error()))
			return
		}
		return
	}

	// Tab-based filtering
	tab := r.URL.Query().Get("tab")
	if tab != "completed" {
		tab = "upcoming"
	}

	var workouts []database.Workout
	var err error
	if tab == "completed" {
		workouts, err = h.cfg.DB.GetCompletedUserWorkouts(r.Context(), userID)
	} else {
		workouts, err = h.cfg.DB.GetUpcomingUserWorkouts(r.Context(), userID)
	}
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to get user workouts", slog.String("error", err.Error()))
		return
	}

	// HTMX tab switch — return just the card grid fragment
	if target == "workouts-content" {
		err = templates.WorkoutsCardGrid(workouts, tab).Render(r.Context(), w)
		if err != nil {
			HandleInternalServerError(w, r)
			h.cfg.Logger.Error("failed to render workouts card grid", slog.String("error", err.Error()))
			return
		}
		return
	}

	// Full page render
	contents := templates.WorkoutsPage(workouts, tab)
	err = templates.Layout(contents, "Fithub | Workouts", true).Render(r.Context(), w)
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to render user workouts", slog.String("error", err.Error()))
		return
	}
}

func (h *Handler) CreateUserWorkout(w http.ResponseWriter, r *http.Request) {
	userID, ok := cntx.UserID(r.Context())
	if !ok {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("missing user id in context")
		return
	}

	reqTitle := r.FormValue("title")
	reqDescription := r.FormValue("workout-description")
	reqDuration := r.FormValue("duration")
	reqPlannedDate := r.FormValue("planned-date")

	if errs := validate.Fields(
		validate.Required(reqTitle, "title"),
		validate.Required(reqDuration, "duration"),
		validate.Required(reqPlannedDate, "planned date"),
		validate.MaxLen(reqTitle, 100, "title"),
		validate.MaxLen(reqDescription, 500, "description"),
	); errs != nil {
		HandleBadRequest(w, r, errs[0].Error())
		h.cfg.Logger.Info("invalid form field", slog.Any("fields", errs))
		return
	}

	description := sql.NullString{Valid: false}
	if reqDescription != "" {
		description.String = reqDescription
		description.Valid = true
	}

	duration, err := strconv.ParseInt(reqDuration, 10, 32)
	if err != nil {
		HandleBadRequest(w, r, "duration must be a number")
		h.cfg.Logger.Info("invalid duration input", slog.String("value", reqDuration), slog.String("error", err.Error()))
		return
	}

	plannedDate, err := time.Parse(time.DateOnly, reqPlannedDate)
	if err != nil {
		HandleBadRequest(w, r, "planned date must be in YYYY-MM-DD format")
		h.cfg.Logger.Info("invalid planned date input", slog.String("value", reqPlannedDate), slog.String("error", err.Error()))
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

	workouts, err := h.cfg.DB.GetUpcomingUserWorkouts(r.Context(), userID)
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to fetch upcoming workouts", slog.String("error", err.Error()))
		return
	}

	w.Header().Set("HX-Push-Url", "/workouts?tab=upcoming")
	err = templates.WorkoutsCardGrid(workouts, "upcoming").Render(r.Context(), w)
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to render workouts card grid", slog.String("error", err.Error()))
		return
	}
}

func (h *Handler) EditUserWorkout(w http.ResponseWriter, r *http.Request) {
	userID, ok := cntx.UserID(r.Context())
	if !ok {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("missing user id in context")
		return
	}
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

	if errs := validate.Fields(
		validate.Required(reqTitle, "title"),
		validate.Required(reqDuration, "duration"),
		validate.Required(reqPlannedDate, "planned date"),
		validate.MaxLen(reqTitle, 100, "title"),
		validate.MaxLen(reqDescription, 500, "description"),
	); errs != nil {
		HandleBadRequest(w, r, errs[0].Error())
		return
	}

	description := sql.NullString{Valid: false}
	if reqDescription != "" {
		description.String = reqDescription
		description.Valid = true
	}

	duration, err := strconv.ParseInt(reqDuration, 10, 32)
	if err != nil {
		HandleBadRequest(w, r, "duration must be a number")
		h.cfg.Logger.Info("invalid duration input", slog.String("value", reqDuration), slog.String("error", err.Error()))
		return
	}

	plannedDate, err := time.Parse(time.DateOnly, reqPlannedDate)
	if err != nil {
		HandleBadRequest(w, r, "planned date must be in YYYY-MM-DD format")
		h.cfg.Logger.Info("invalid planned date input", slog.String("value", reqPlannedDate), slog.String("error", err.Error()))
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

	err = templates.WorkoutCard(updatedWorkout).Render(r.Context(), w)
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to render workout card", slog.String("error", err.Error()))
		return
	}
}

func (h *Handler) DeleteUserWorkout(w http.ResponseWriter, r *http.Request) {
	userID, ok := cntx.UserID(r.Context())
	if !ok {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("missing user id in context")
		return
	}
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
	userID, ok := cntx.UserID(r.Context())
	if !ok {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("missing user id in context")
		return
	}
	workoutID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to parse workout id", slog.String("error", err.Error()))
		return
	}

	workout, err := h.cfg.DB.GetWorkoutByID(r.Context(), database.GetWorkoutByIDParams{
		ID:     workoutID,
		UserID: userID,
	})
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to fetch workout", slog.String("error", err.Error()))
		return
	}

	workoutExercises, err := h.cfg.DB.WorkoutAndExercises(r.Context(), database.WorkoutAndExercisesParams{
		WorkoutID: workoutID,
		UserID:    userID,
	})
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
