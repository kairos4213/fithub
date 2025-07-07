package handlers

import (
	"database/sql"
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
		return // TODO: handle err
	}

	w.Header().Set("Content-Type", "text/html")
	contents := templates.Workouts(workouts)
	templates.Layout(contents, "Fithub | Workouts", true).Render(r.Context(), w)
}

func (h *Handler) CreateUserWorkout(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(cntx.UserIDKey).(uuid.UUID)

	err := r.ParseForm()
	if err != nil {
		return // TODO: handle error
	}

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
		return // TODO: handle err
	}

	plannedDate, err := time.Parse(time.DateOnly, reqPlannedDate)
	if err != nil {
		return // TODO: handle err
	}

	newWorkout, err := h.DB.CreateWorkout(r.Context(), database.CreateWorkoutParams{
		UserID: userID, Title: reqTitle, Description: description, DurationMinutes: int32(duration), PlannedDate: plannedDate,
	})
	if err != nil {
		return // TODO: handle err
	}

	templates.WorkoutDataRow(newWorkout).Render(r.Context(), w)
}
