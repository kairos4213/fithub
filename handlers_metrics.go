package main

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/kairos4213/fithub/internal/database"
)

type Metric struct {
	ID          uuid.UUID `json:"metric_id,omitempty"`
	MetricType  string    `json:"metric_type,omitempty"`
	Measurement string    `json:"measurement,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
	UserID      uuid.UUID `json:"user_id,omitempty"`
}

func (cfg *apiConfig) addMetricsHandler(w http.ResponseWriter, r *http.Request) {
	type requestParams struct {
		Measurement string `json:"measurement"`
	}
	userID := r.Context().Value(userIDKey).(uuid.UUID)

	reqParams := requestParams{}
	if err := parseJSON(r, &reqParams); err != nil {
		respondWithError(w, http.StatusBadRequest, "malformed request", err)
		return
	}

	metricType := r.PathValue("type")
	switch metricType {
	case "body_weight":
		bodyWeightEntry, err := cfg.db.AddBodyWeight(r.Context(), database.AddBodyWeightParams{
			UserID:      userID,
			Measurement: reqParams.Measurement,
		})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error saving body weight metric", err)
			return
		}

		respondWithJSON(w, http.StatusCreated, Metric{
			ID:          bodyWeightEntry.ID,
			MetricType:  metricType,
			Measurement: bodyWeightEntry.Measurement,
			CreatedAt:   bodyWeightEntry.CreatedAt,
			UpdatedAt:   bodyWeightEntry.UpdatedAt,
			UserID:      userID,
		})
	}
}
