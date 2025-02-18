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
	case "body_weights":
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
			MetricType:  "body_weight",
			Measurement: bodyWeightEntry.Measurement,
			CreatedAt:   bodyWeightEntry.CreatedAt,
			UpdatedAt:   bodyWeightEntry.UpdatedAt,
			UserID:      userID,
		})
	case "muscle_masses":
		muscleMassEntry, err := cfg.db.AddMuscleMass(r.Context(), database.AddMuscleMassParams{
			UserID:      userID,
			Measurement: reqParams.Measurement,
		})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error saving muscle mass metric", err)
			return
		}

		respondWithJSON(w, http.StatusCreated, Metric{
			ID:          userID,
			MetricType:  "muscle_mass",
			Measurement: muscleMassEntry.Measurement,
			CreatedAt:   muscleMassEntry.CreatedAt,
			UpdatedAt:   muscleMassEntry.UpdatedAt,
			UserID:      userID,
		})
	case "body_fat_percents":
		bfPercentEntry, err := cfg.db.AddBodyFatPerc(r.Context(), database.AddBodyFatPercParams{
			UserID:      userID,
			Measurement: reqParams.Measurement,
		})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error saving body fat percentage", err)
			return
		}

		respondWithJSON(w, http.StatusCreated, Metric{
			ID:          userID,
			MetricType:  "body_fat_percentage",
			Measurement: bfPercentEntry.Measurement,
			CreatedAt:   bfPercentEntry.CreatedAt,
			UpdatedAt:   bfPercentEntry.UpdatedAt,
			UserID:      userID,
		})
	}
}

func (cfg *apiConfig) getAllUserMetrics(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIDKey).(uuid.UUID)

	bodyWeightsResp := []Metric{}
	// muscleMassesResp := []Metric{}
	// bfPercentsResp := []Metric{}

	bodyWeights, err := cfg.db.GetAllBodyWeights(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error retrieving body weights", err)
		return
	}
	for _, bodyWeight := range bodyWeights {
		bodyWeightsResp = append(bodyWeightsResp, Metric{
			ID:          bodyWeight.ID,
			MetricType:  "body_weight",
			Measurement: bodyWeight.Measurement,
			CreatedAt:   bodyWeight.CreatedAt,
			UpdatedAt:   bodyWeight.UpdatedAt,
			UserID:      bodyWeight.UserID,
		})
	}

	resp := map[string][]Metric{
		"body_weights": bodyWeightsResp,
	}
	respondWithJSON(w, http.StatusOK, resp)
}
