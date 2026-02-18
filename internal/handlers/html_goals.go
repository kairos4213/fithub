package handlers

import (
	"database/sql"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kairos4213/fithub/internal/cntx"
	"github.com/kairos4213/fithub/internal/database"
	"github.com/kairos4213/fithub/internal/templates"
	"github.com/kairos4213/fithub/internal/validate"
)

func (h *Handler) GetAllGoals(w http.ResponseWriter, r *http.Request) {
	userID, ok := cntx.UserID(r.Context())
	if !ok {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("missing user id in context")
		return
	}

	// Tab-based filtering
	tab := r.URL.Query().Get("tab")
	if tab != "completed" {
		tab = "in_progress"
	}

	var goals []database.Goal
	var err error
	if tab == "completed" {
		goals, err = h.cfg.DB.GetCompletedGoals(r.Context(), userID)
	} else {
		goals, err = h.cfg.DB.GetInProgressGoals(r.Context(), userID)
	}
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to get user goals", slog.String("error", err.Error()))
		return
	}

	// HTMX tab switch — return just the card grid fragment
	target := r.Header.Get("HX-Target")
	if target == "goals-content" {
		err = templates.GoalsCardGrid(goals, tab).Render(r.Context(), w)
		if err != nil {
			HandleInternalServerError(w, r)
			h.cfg.Logger.Error("failed to render goals card grid", slog.String("error", err.Error()))
			return
		}
		return
	}

	// Full page render
	contents := templates.GoalsPage(goals, tab)
	err = templates.Layout(contents, "Fithub | Goals", true).Render(r.Context(), w)
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to render goals", slog.String("error", err.Error()))
		return
	}
}

func (h *Handler) AddNewGoal(w http.ResponseWriter, r *http.Request) {
	userID, ok := cntx.UserID(r.Context())
	if !ok {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("missing user id in context")
		return
	}

	reqGoalName := r.FormValue("goal-name")
	reqDescription := r.FormValue("description")
	reqGoalDate := r.FormValue("goal-date")
	reqNotes := r.FormValue("notes")

	inputFields := []string{"goal-name", "description", "goal-date", "notes"}
	if errs := validate.Fields(
		validate.Required(reqGoalName, "goal name"),
		validate.Required(reqDescription, "description"),
		validate.Required(reqGoalDate, "goal date"),
		validate.MaxLen(reqGoalName, 100, "goal name"),
		validate.MaxLen(reqDescription, 500, "description"),
		validate.MaxLen(reqNotes, 500, "notes"),
	); errs != nil {
		HandleFieldErrors(w, r, h.cfg.Logger, errs, inputFields, "")
		return
	}

	goalDate, err := time.Parse(time.DateOnly, reqGoalDate)
	if err != nil {
		HandleFieldErrors(w, r, h.cfg.Logger, []validate.FieldError{{Field: "goal date", Message: "goal date must be in YYYY-MM-DD format"}}, inputFields, "")
		return
	}

	notes := sql.NullString{Valid: false}
	if reqNotes != "" {
		notes.String = reqNotes
		notes.Valid = true
	}

	_, err = h.cfg.DB.CreateGoal(r.Context(), database.CreateGoalParams{
		GoalName:    reqGoalName,
		Description: reqDescription,
		GoalDate:    goalDate,
		Notes:       notes,
		UserID:      userID,
	})
	if err != nil {
		if strings.Contains(err.Error(), "pq: duplicate key value violates unique constraint") {
			HandleFieldErrors(w, r, h.cfg.Logger, []validate.FieldError{{Field: "goal name", Message: "goal already exists"}}, inputFields, "")
			return
		}
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to create goal", slog.String("error", err.Error()))
		return
	}

	goals, err := h.cfg.DB.GetInProgressGoals(r.Context(), userID)
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to fetch in-progress goals", slog.String("error", err.Error()))
		return
	}

	w.Header().Set("HX-Trigger", "close-create-goal")
	w.Header().Set("HX-Push-Url", "/goals?tab=in_progress")
	err = templates.GoalsCardGrid(goals, "in_progress").Render(r.Context(), w)
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to render goals card grid", slog.String("error", err.Error()))
		return
	}
}

func (h *Handler) EditGoal(w http.ResponseWriter, r *http.Request) {
	userID, ok := cntx.UserID(r.Context())
	if !ok {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("missing user id in context")
		return
	}
	goalID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to parse goal id", slog.String("error", err.Error()))
		return
	}

	reqGoalName := r.FormValue("goal-name")
	reqStatus := r.FormValue("status")
	reqDescription := r.FormValue("description")
	reqGoalDate := r.FormValue("goal-date")
	reqNotes := r.FormValue("notes")

	prefix := goalID.String() + "-"
	editFields := []string{prefix + "goal-name", prefix + "description", prefix + "goal-date", prefix + "notes", prefix + "status"}

	if errs := validate.Fields(
		validate.Required(reqGoalName, "goal name"),
		validate.Required(reqStatus, "status"),
		validate.Required(reqDescription, "description"),
		validate.Required(reqGoalDate, "goal date"),
		validate.MaxLen(reqGoalName, 100, "goal name"),
		validate.MaxLen(reqDescription, 500, "description"),
		validate.MaxLen(reqNotes, 500, "notes"),
	); errs != nil {
		HandleFieldErrors(w, r, h.cfg.Logger, errs, editFields, prefix)
		return
	}

	goalDate, err := time.Parse(time.DateOnly, reqGoalDate)
	if err != nil {
		HandleFieldErrors(w, r, h.cfg.Logger, []validate.FieldError{{Field: "goal date", Message: "goal date must be in YYYY-MM-DD format"}}, editFields, prefix)
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
		if strings.Contains(err.Error(), "pq: duplicate key value violates unique constraint") {
			HandleFieldErrors(w, r, h.cfg.Logger, []validate.FieldError{{Field: "goal name", Message: "goal already exists"}}, editFields, prefix)
			return
		}
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to update goal", slog.String("error", err.Error()))
		return
	}

	err = templates.GoalCard(updatedGoal).Render(r.Context(), w)
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to render goal card", slog.String("error", err.Error()))
		return
	}
}

func (h *Handler) DeleteGoal(w http.ResponseWriter, r *http.Request) {
	userID, ok := cntx.UserID(r.Context())
	if !ok {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("missing user id in context")
		return
	}
	goalID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to parse goal id", slog.String("error", err.Error()))
		return
	}

	err = h.cfg.DB.DeleteGoal(r.Context(), database.DeleteGoalParams{
		ID:     goalID,
		UserID: userID,
	})
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to delete goal", slog.String("error", err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
}
