package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/kairos4213/fithub/internal/cntx"
	"github.com/kairos4213/fithub/internal/database"
	"github.com/kairos4213/fithub/internal/templates"
)

func (h *Handler) GetAllGoals(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(cntx.UserIDKey).(uuid.UUID)
	goals, err := h.DB.GetAllUserGoals(r.Context(), userID)
	if err != nil {
		return // TODO: handle error
	}

	contents := templates.Goals(goals)
	templates.Layout(contents, "Fithub | Goals", true).Render(r.Context(), w)
}

func (h *Handler) AddNewGoal(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(cntx.UserIDKey).(uuid.UUID)

	err := r.ParseForm()
	if err != nil {
		return // TODO: handle error
	}

	reqGoalName := r.FormValue("goal-name")
	reqDescription := r.FormValue("description")
	reqGoalDate := r.FormValue("goal-date")
	reqNotes := r.FormValue("notes")

	goalDate, err := time.Parse(time.DateOnly, reqGoalDate)
	if err != nil {
		return // TODO: handle err
	}

	notes := sql.NullString{Valid: false}
	if reqNotes != "" {
		notes.String = reqNotes
		notes.Valid = true
	}

	newGoal, err := h.DB.CreateGoal(r.Context(), database.CreateGoalParams{
		GoalName:    reqGoalName,
		Description: reqDescription,
		GoalDate:    goalDate,
		Notes:       notes,
		UserID:      userID,
	})
	if err != nil {
		return // TODO: handle error
	}

	templates.NewGoal(newGoal).Render(r.Context(), w)
}
