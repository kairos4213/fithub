package handlers

import (
	"database/sql"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kairos4213/fithub/internal/cntx"
	"github.com/kairos4213/fithub/internal/database"
	"github.com/kairos4213/fithub/internal/utils"
)

type Goal struct {
	ID             string `json:"goal_id,omitempty"`
	CreatedAt      string `json:"created_at,omitempty"`
	UpdatedAt      string `json:"updated_at,omitempty"`
	Name           string `json:"name,omitempty"`
	Description    string `json:"description,omitempty"`
	GoalDate       string `json:"goal_date,omitempty"`
	CompletionDate string `json:"completion_date,omitempty"`
	Notes          string `json:"notes,omitempty"`
	Status         string `json:"status,omitempty"`
	UserID         string `json:"user_id,omitempty"`
}

func (h *Handler) CreateGoal(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(cntx.UserIDKey).(uuid.UUID)

	reqParams := Goal{}
	if err := utils.ParseJSON(r, &reqParams); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "malformed request", err)
		return
	}

	goalDate, err := time.Parse(time.DateOnly, reqParams.GoalDate)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error parsing date", err)
		return
	}

	goalNotes := sql.NullString{Valid: false}
	if reqParams.Notes != "" {
		goalNotes.Valid = true
		goalNotes.String = reqParams.Notes
	}

	goal, err := h.DB.CreateGoal(r.Context(), database.CreateGoalParams{
		GoalName:    strings.ToLower(reqParams.Name),
		Description: reqParams.Description,
		GoalDate:    goalDate.UTC(),
		Notes:       goalNotes,
		UserID:      userID,
	})
	if err != nil {
		if strings.Contains(err.Error(), `pq: duplicate key value violates unique constraint "goals_name_user_id_key"`) {
			utils.RespondWithError(w, http.StatusBadRequest, "Cannot have duplicate goal names", err)
			return
		}
		utils.RespondWithError(w, http.StatusInternalServerError, "Error saving goal", err)
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, Goal{
		ID:          goal.ID.String(),
		CreatedAt:   goal.CreatedAt.Format(time.RFC822),
		UpdatedAt:   goal.UpdatedAt.Format(time.RFC822),
		Name:        goal.GoalName,
		Description: goal.Description,
		GoalDate:    goal.GoalDate.Format(time.DateOnly),
		Notes:       goal.Notes.String,
		UserID:      goal.UserID.String(),
	})
}

func (h *Handler) GetAllUserGoals(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(cntx.UserIDKey).(uuid.UUID)
	goals, err := h.DB.GetAllUserGoals(r.Context(), userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error getting goals", err)
		return
	}

	response := []Goal{}
	for _, goal := range goals {
		response = append(response, Goal{
			ID:             goal.ID.String(),
			CreatedAt:      goal.CreatedAt.Format(time.RFC822),
			UpdatedAt:      goal.UpdatedAt.Format(time.RFC822),
			Name:           goal.GoalName,
			Description:    goal.Description,
			GoalDate:       goal.GoalDate.String(),
			CompletionDate: goal.CompletionDate.Time.Format(time.DateOnly),
			Notes:          goal.Notes.String,
			Status:         goal.Status,
			UserID:         goal.UserID.String(),
		})
	}
	utils.RespondWithJSON(w, http.StatusOK, response)
}

func (h *Handler) UpdateGoal(w http.ResponseWriter, r *http.Request) {
	goalID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid goal id", err)
		return
	}

	reqParams := Goal{}
	if err := utils.ParseJSON(r, &reqParams); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "malformed request", err)
		return
	}

	updateGoalParams := database.UpdateGoalParams{
		ID:          goalID,
		GoalName:    reqParams.Name,
		Description: reqParams.Description,
		Status:      reqParams.Status,
	}

	goalDate, err := time.Parse(time.DateOnly, reqParams.GoalDate)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error parsing date", err)
		return
	}
	updateGoalParams.GoalDate = goalDate

	if reqParams.CompletionDate == "" {
		updateGoalParams.CompletionDate = sql.NullTime{Valid: false}
	} else {
		completionDate, err := time.Parse(time.DateOnly, reqParams.CompletionDate)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error parsing date", err)
			return
		}
		updateGoalParams.CompletionDate = sql.NullTime{Time: completionDate, Valid: true}
	}

	if reqParams.Notes == "" {
		updateGoalParams.Notes = sql.NullString{Valid: false}
	} else {
		updateGoalParams.Notes = sql.NullString{String: reqParams.Notes, Valid: true}
	}

	goal, err := h.DB.UpdateGoal(r.Context(), updateGoalParams)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error updating goal", err)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, Goal{
		ID:             goal.ID.String(),
		CreatedAt:      goal.CreatedAt.Format(time.RFC822),
		UpdatedAt:      goal.UpdatedAt.Format(time.RFC822),
		Name:           goal.GoalName,
		Description:    goal.Description,
		GoalDate:       goal.GoalDate.Format(time.DateOnly),
		CompletionDate: goal.CompletionDate.Time.Format(time.DateOnly),
		Notes:          goal.Notes.String,
		Status:         goal.Status,
		UserID:         goal.UserID.String(),
	})
}

func (h *Handler) DeleteGoal(w http.ResponseWriter, r *http.Request) {
	goalID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error parsing goal id", err)
		return
	}

	if err := h.DB.DeleteGoal(r.Context(), goalID); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error deleting goal", err)
		return
	}

	utils.RespondWithJSON(w, http.StatusNoContent, Goal{})
}

func (h *Handler) DeleteAllUserGoals(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(cntx.UserIDKey).(uuid.UUID)

	if err := h.DB.DeleteAllUserGoals(r.Context(), userID); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error deleting goals", err)
		return
	}

	utils.RespondWithJSON(w, http.StatusNoContent, Goal{})
}
