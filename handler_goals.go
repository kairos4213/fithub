package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kairos4213/fithub/internal/auth"
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

func (cfg *apiConfig) handlerGoalsCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		GoalDate    string `json:"goal_date"`
		Notes       string `json:"notes"`
	}

	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Missing JWT", err)
		return
	}

	userID, err := auth.ValidateJWT(accessToken, cfg.publicKey)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid JWT", err)
		return
	}

	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding request", err)
		return
	}

	goalDate, err := time.Parse(time.DateOnly, params.GoalDate)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error parsing date", err)
		return
	}

	goal, err := cfg.db.CreateGoal(r.Context(), database.CreateGoalParams{
		Name:        strings.ToLower(params.Name),
		Description: params.Description,
		GoalDate:    goalDate.UTC(),
		Notes:       sql.NullString{String: strings.ToLower(params.Notes)},
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
		Name:        goal.Name,
		Description: goal.Description,
		GoalDate:    goal.GoalDate,
		Notes:       goal.Notes.String,
		UserID:      goal.UserID,
	})
}

func (cfg *apiConfig) handlerGoalsGetAll(w http.ResponseWriter, r *http.Request) {
	accesToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Missing JWT", err)
		return
	}

	userID, err := auth.ValidateJWT(accesToken, cfg.publicKey)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid JWT", err)
		return
	}

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
			Name:           goal.Name,
			Description:    goal.Description,
			GoalDate:       goal.GoalDate,
			CompletionDate: goal.CompletionDate.Time,
			Notes:          goal.Notes.String,
			Status:         goal.Status,
			UserID:         userID,
		})
	}
	respondWithJSON(w, http.StatusOK, response)
}
