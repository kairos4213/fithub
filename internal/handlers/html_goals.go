package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/kairos4213/fithub/internal/cntx"
	"github.com/kairos4213/fithub/internal/database"
	"github.com/kairos4213/fithub/internal/templates"
)

func (h *Handler) GetAllGoals(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(cntx.UserIDKey).(uuid.UUID)
	goals, err := h.cfg.DB.GetAllUserGoals(r.Context(), userID)
	if err != nil {
		HandleInternalServerError(w, r)
		log.Printf("%v", err)
		return
	}

	contents := templates.Goals(goals)
	templates.Layout(contents, "Fithub | Goals", true).Render(r.Context(), w)
}

func (h *Handler) AddNewGoal(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(cntx.UserIDKey).(uuid.UUID)

	reqGoalName := r.FormValue("goal-name")
	reqDescription := r.FormValue("description")
	reqGoalDate := r.FormValue("goal-date")
	reqNotes := r.FormValue("notes")

	goalDate, err := time.Parse(time.DateOnly, reqGoalDate)
	if err != nil {
		HandleInternalServerError(w, r)
		log.Printf("%v", err)
		return
	}

	notes := sql.NullString{Valid: false}
	if reqNotes != "" {
		notes.String = reqNotes
		notes.Valid = true
	}

	newGoal, err := h.cfg.DB.CreateGoal(r.Context(), database.CreateGoalParams{
		GoalName:    reqGoalName,
		Description: reqDescription,
		GoalDate:    goalDate,
		Notes:       notes,
		UserID:      userID,
	})
	if err != nil {
		HandleInternalServerError(w, r)
		log.Printf("%v", err)
		return
	}

	templates.GoalDataRow(newGoal).Render(r.Context(), w)
}

func (h *Handler) EditGoal(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(cntx.UserIDKey).(uuid.UUID)
	goalID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		HandleInternalServerError(w, r)
		log.Printf("%v", err)
		return
	}

	reqGoalName := r.FormValue("goal-name")
	reqStatus := r.FormValue("status")
	reqDescription := r.FormValue("description")
	reqGoalDate := r.FormValue("goal-date")
	reqNotes := r.FormValue("notes")

	goalDate, err := time.Parse(time.DateOnly, reqGoalDate)
	if err != nil {
		HandleInternalServerError(w, r)
		log.Printf("%v", err)
		return
	}

	notes := sql.NullString{Valid: false}
	if reqNotes != "" {
		notes.String = reqNotes
		notes.Valid = true
	}

	completionDate := sql.NullTime{Valid: false}
	if reqStatus == "completed" {
		completionDate.Time = time.Now().UTC()
		completionDate.Valid = true
	}

	updatedGoal, err := h.cfg.DB.UpdateGoal(r.Context(), database.UpdateGoalParams{
		GoalName:       reqGoalName,
		Description:    reqDescription,
		GoalDate:       goalDate,
		CompletionDate: completionDate,
		Notes:          notes,
		Status:         reqStatus,
		ID:             goalID,
		UserID:         userID,
	})
	if err != nil {
		HandleInternalServerError(w, r)
		log.Printf("%v", err)
		return
	}

	templates.GoalDataRow(updatedGoal).Render(r.Context(), w)
}

func (h *Handler) DeleteGoal(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(cntx.UserIDKey).(uuid.UUID)
	goalID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		HandleInternalServerError(w, r)
		log.Printf("%v", err)
		return
	}

	err = h.cfg.DB.DeleteGoal(r.Context(), database.DeleteGoalParams{
		ID:     goalID,
		UserID: userID,
	})
	if err != nil {
		HandleInternalServerError(w, r)
		log.Printf("%v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
