package handlers

import (
	"database/sql"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/kairos4213/fithub/internal/cntx"
	"github.com/kairos4213/fithub/internal/database"
	"github.com/kairos4213/fithub/internal/templates"
	"github.com/kairos4213/fithub/internal/validate"
)

func (h *Handler) GetAllGoals(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(cntx.UserIDKey).(uuid.UUID)
	goals, err := h.cfg.DB.GetAllUserGoals(r.Context(), userID)
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to get user goals", slog.String("error", err.Error()))
		return
	}

	contents := templates.Goals(goals)
	err = templates.Layout(contents, "Fithub | Goals", true).Render(r.Context(), w)
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to render goals", slog.String("error", err.Error()))
		return
	}
}

func (h *Handler) AddNewGoal(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(cntx.UserIDKey).(uuid.UUID)

	reqGoalName := r.FormValue("goal-name")
	reqDescription := r.FormValue("description")
	reqGoalDate := r.FormValue("goal-date")
	reqNotes := r.FormValue("notes")

	if errs := validate.Fields(
		validate.Required(reqGoalName, "goal name"),
		validate.Required(reqDescription, "description"),
		validate.Required(reqGoalDate, "goal date"),
		validate.MaxLen(reqGoalName, 100, "goal name"),
		validate.MaxLen(reqDescription, 500, "description"),
		validate.MaxLen(reqNotes, 500, "notes"),
	); errs != nil {
		HandleBadRequest(w, r, errs[0].Error())
		h.cfg.Logger.Info("invalid form field", slog.Any("fields", errs))
		return
	}

	goalDate, err := time.Parse(time.DateOnly, reqGoalDate)
	if err != nil {
		HandleBadRequest(w, r, "goal date must be in YYYY-MM-DD format")
		h.cfg.Logger.Info("invalid goal date input", slog.String("value", reqGoalDate), slog.String("error", err.Error()))
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
		h.cfg.Logger.Error("failed to create goal", slog.String("error", err.Error()))
		return
	}

	err = templates.GoalDataRow(newGoal).Render(r.Context(), w)
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to render goal row", slog.String("error", err.Error()))
		return
	}
}

func (h *Handler) EditGoal(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(cntx.UserIDKey).(uuid.UUID)
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

	if errs := validate.Fields(
		validate.Required(reqGoalName, "goal name"),
		validate.Required(reqStatus, "status"),
		validate.Required(reqDescription, "description"),
		validate.Required(reqGoalDate, "goal date"),
		validate.MaxLen(reqGoalName, 100, "goal name"),
		validate.MaxLen(reqDescription, 500, "description"),
		validate.MaxLen(reqNotes, 500, "notes"),
	); errs != nil {
		HandleBadRequest(w, r, errs[0].Error())
		h.cfg.Logger.Info("invalid form field", slog.Any("fields", errs))
		return
	}

	goalDate, err := time.Parse(time.DateOnly, reqGoalDate)
	if err != nil {
		HandleBadRequest(w, r, "goal date must be in YYYY-MM-DD format")
		h.cfg.Logger.Info("invalid goal date input", slog.String("value", reqGoalDate), slog.String("error", err.Error()))
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
		h.cfg.Logger.Error("failed to update goal", slog.String("error", err.Error()))
		return
	}

	err = templates.GoalDataRow(updatedGoal).Render(r.Context(), w)
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to render goal row", slog.String("error", err.Error()))
		return
	}
}

func (h *Handler) DeleteGoal(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(cntx.UserIDKey).(uuid.UUID)
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
