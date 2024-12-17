package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/kairos4213/fithub/internal/auth"
	"github.com/kairos4213/fithub/internal/database"
)

type Goal struct {
	ID             uuid.UUID `json:"goal_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	GoalDate       time.Time `json:"goal_date"`
	CompletionDate time.Time `json:"completion_date"`
	Notes          string    `json:"notes"`
	Status         string    `json:"status"`
	UserID         uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerGoalsCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
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

	goal, err := cfg.db.CreateGoal(r.Context(), database.CreateGoalParams{})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error saving goal", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, Goal{
		
	})
}
