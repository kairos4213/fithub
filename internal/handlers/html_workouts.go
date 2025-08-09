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
	"github.com/kairos4213/fithub/internal/templates"
)

func (h *Handler) GetUserWorkouts(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(cntx.UserIDKey).(uuid.UUID)
	workouts, err := h.DB.GetAllUserWorkouts(r.Context(), userID)
	if err != nil {
		HandleInternalServerError(w, r)
		log.Printf("%v", err)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	contents := templates.Workouts(workouts)
	templates.Layout(contents, "Fithub | Workouts", true).Render(r.Context(), w)
}

func (h *Handler) CreateUserWorkout(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(cntx.UserIDKey).(uuid.UUID)

	reqTitle := r.FormValue("title")
	reqDescription := r.FormValue("workout-description")
	reqDuration := r.FormValue("duration")
	reqPlannedDate := r.FormValue("planned-date")

	description := sql.NullString{Valid: false}
	if reqDescription != "" {
		description.String = reqDescription
		description.Valid = true
	}

	duration, err := strconv.ParseInt(reqDuration, 10, 32)
	if err != nil {
		HandleInternalServerError(w, r)
		log.Printf("%v", err)
		return
	}

	plannedDate, err := time.Parse(time.DateOnly, reqPlannedDate)
	if err != nil {
		HandleInternalServerError(w, r)
		log.Printf("%v", err)
		return
	}

	newWorkout, err := h.DB.CreateWorkout(r.Context(), database.CreateWorkoutParams{
		UserID: userID, Title: reqTitle, Description: description, DurationMinutes: int32(duration), PlannedDate: plannedDate,
	})
	if err != nil {
		HandleInternalServerError(w, r)
		log.Printf("%v", err)
		return
	}

	templates.WorkoutsDataRow(newWorkout).Render(r.Context(), w)
}

func (h *Handler) EditUserWorkout(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(cntx.UserIDKey).(uuid.UUID)
	workoutID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		HandleInternalServerError(w, r)
		log.Printf("%v", err)
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
		log.Printf("%v", err)
		return
	}

	plannedDate, err := time.Parse(time.DateOnly, reqPlannedDate)
	if err != nil {
		HandleInternalServerError(w, r)
		log.Printf("%v", err)
		return
	}

	dateCompleted := sql.NullTime{Valid: false}
	if reqCompletionDate != "" {
		date, err := time.Parse(time.DateOnly, reqCompletionDate)
		if err != nil {
			HandleInternalServerError(w, r)
			log.Printf("%v", err)
			return
		}
		dateCompleted.Time = date
		dateCompleted.Valid = true
	}

	updatedWorkout, err := h.DB.UpdateWorkout(r.Context(), database.UpdateWorkoutParams{
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
		log.Printf("%v", err)
		return
	}

	templates.WorkoutsDataRow(updatedWorkout).Render(r.Context(), w)
}

func (h *Handler) DeleteUserWorkout(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(cntx.UserIDKey).(uuid.UUID)
	workoutID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		HandleInternalServerError(w, r)
		log.Printf("%v", err)
		return
	}

	err = h.DB.DeleteWorkout(r.Context(), database.DeleteWorkoutParams{
		ID:     workoutID,
		UserID: userID,
	})
	if err != nil {
		HandleInternalServerError(w, r)
		log.Printf("%v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetUserWorkoutExercises(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(cntx.UserIDKey).(uuid.UUID)
	workoutID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		HandleInternalServerError(w, r)
		log.Printf("%v", err)
		return
	}

	title := r.FormValue("title")
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
		log.Printf("%v", err)
		return
	}

	plannedDate, err := time.Parse(time.DateOnly, reqPlannedDate)
	if err != nil {
		HandleInternalServerError(w, r)
		log.Printf("%v", err)
		return
	}

	dateCompleted := sql.NullTime{Valid: false}
	if reqCompletionDate != "" {
		date, err := time.Parse(time.DateOnly, reqCompletionDate)
		if err != nil {
			HandleInternalServerError(w, r)
			log.Printf("%v", err)
			return
		}
		dateCompleted.Time = date
		dateCompleted.Valid = true
	}

	workout := database.Workout{
		ID:              workoutID,
		UserID:          userID,
		Title:           title,
		Description:     description,
		DurationMinutes: int32(duration),
		PlannedDate:     plannedDate,
		DateCompleted:   dateCompleted,
	}

	workoutExercises, err := h.DB.ExercisesForWorkout(r.Context(), workoutID)
	if err != nil {
		HandleInternalServerError(w, r)
		log.Printf("%v", err)
		return
	}

	contents := templates.WorkoutPage(workout, workoutExercises, []database.Exercise{})
	templates.Layout(contents, "FitHub | Workout Page", true).Render(r.Context(), w)
}
