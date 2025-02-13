package main

import (
	"database/sql"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kairos4213/fithub/internal/database"
)

type Goal struct {
	ID             uuid.UUID `json:"goal_id,omitempty"`
	CreatedAt      time.Time `json:"created_at,omitempty"`
	UpdatedAt      time.Time `json:"updated_at,omitempty"`
	Name           string    `json:"name,omitempty"`
	Description    string    `json:"description,omitempty"`
	GoalDate       time.Time `json:"goal_date,omitempty"`
	CompletionDate time.Time `json:"completion_date,omitempty"`
	Notes          string    `json:"notes,omitempty"`
	Status         string    `json:"status,omitempty"`
	UserID         uuid.UUID `json:"user_id,omitempty"`
}

func (cfg *apiConfig) createGoalsHandler(w http.ResponseWriter, r *http.Request) {
	type requestParams struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		GoalDate    string `json:"goal_date"`
		Notes       string `json:"notes"`
	}

	userID := r.Context().Value(userIDKey).(uuid.UUID)

	reqParams := requestParams{}
	if err := parseJSON(r, &reqParams); err != nil {
		respondWithError(w, http.StatusBadRequest, "malformed request", err)
		return
	}

	goalDate, err := time.Parse(time.DateOnly, reqParams.GoalDate)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error parsing date", err)
		return
	}

	goal, err := cfg.db.CreateGoal(r.Context(), database.CreateGoalParams{
		GoalName:    strings.ToLower(reqParams.Name),
		Description: reqParams.Description,
		GoalDate:    goalDate.UTC(),
		Notes:       sql.NullString{String: strings.ToLower(reqParams.Notes)},
		UserID:      userID,
	})
	if err != nil {
		if strings.Contains(err.Error(), `pq: duplicate key value violates unique constraint "goals_name_user_id_key"`) {
			respondWithError(w, http.StatusBadRequest, "Cannot have duplicate goal names", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Error saving goal", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, Goal{
		ID:          goal.ID,
		CreatedAt:   goal.CreatedAt,
		UpdatedAt:   goal.UpdatedAt,
		Name:        goal.GoalName,
		Description: goal.Description,
		GoalDate:    goal.GoalDate,
		Notes:       goal.Notes.String,
		UserID:      goal.UserID,
	})
}

func (cfg *apiConfig) getAllGoalsHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIDKey).(uuid.UUID)
	goals, err := cfg.db.GetAllUserGoals(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error getting goals", err)
		return
	}

	response := []Goal{}
	for _, goal := range goals {
		response = append(response, Goal{
			ID:             goal.ID,
			CreatedAt:      goal.CreatedAt,
			UpdatedAt:      goal.UpdatedAt,
			Name:           goal.GoalName,
			Description:    goal.Description,
			GoalDate:       goal.GoalDate,
			CompletionDate: goal.CompletionDate.Time,
			Notes:          goal.Notes.String,
			Status:         goal.Status,
			UserID:         goal.UserID,
		})
	}
	respondWithJSON(w, http.StatusOK, response)
}

func (cfg *apiConfig) updateGoalsHandler(w http.ResponseWriter, r *http.Request) {
	type requestParams struct {
		ID             uuid.UUID `json:"goal_id"`
		Name           string    `json:"goal_name"`
		Description    string    `json:"description"`
		GoalDate       string    `json:"goal_date"`
		CompletionDate string    `json:"completion_date"`
		Notes          string    `json:"notes"`
		Status         string    `json:"status"`
	}

	reqParams := requestParams{}
	if err := parseJSON(r, &reqParams); err != nil {
		respondWithError(w, http.StatusBadRequest, "malformed request", err)
		return
	}

	updateGoalParams := database.UpdateGoalParams{
		ID:          reqParams.ID,
		GoalName:    reqParams.Name,
		Description: reqParams.Description,
		Status:      reqParams.Status,
	}

	goalDate, err := time.Parse(time.DateOnly, reqParams.GoalDate)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error parsing date", err)
		return
	}
	updateGoalParams.GoalDate = goalDate

	if reqParams.CompletionDate == "" {
		updateGoalParams.CompletionDate = sql.NullTime{Valid: false}
	} else {
		completionDate, err := time.Parse(time.DateOnly, reqParams.CompletionDate)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error parsing date", err)
			return
		}
		updateGoalParams.CompletionDate = sql.NullTime{Time: completionDate, Valid: true}
	}

	if reqParams.Notes == "" {
		updateGoalParams.Notes = sql.NullString{Valid: false}
	} else {
		updateGoalParams.Notes = sql.NullString{String: reqParams.Notes, Valid: true}
	}

	goal, err := cfg.db.UpdateGoal(r.Context(), updateGoalParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error updating goal", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Goal{
		ID:             goal.ID,
		CreatedAt:      goal.CreatedAt,
		UpdatedAt:      goal.UpdatedAt,
		Name:           goal.GoalName,
		Description:    goal.Description,
		GoalDate:       goal.GoalDate,
		CompletionDate: goal.CompletionDate.Time,
		Notes:          goal.Notes.String,
		Status:         goal.Status,
		UserID:         goal.UserID,
	})
}

func (cfg *apiConfig) deleteGoalsHandler(w http.ResponseWriter, r *http.Request) {
	goal_id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error parsing goal id", err)
		return
	}

	err = cfg.db.DeleteGoal(r.Context(), goal_id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error deleting goal", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, Goal{})
}
