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
	muscleMassesResp := []Metric{}
	bfPercentsResp := []Metric{}

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

	muscleMasses, err := cfg.db.GetAllMuscleMasses(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error retrieving muscle mass metrics", err)
		return
	}
	for _, muscleMass := range muscleMasses {
		muscleMassesResp = append(muscleMassesResp, Metric{
			ID:          muscleMass.ID,
			MetricType:  "muscle_mass",
			Measurement: muscleMass.Measurement,
			CreatedAt:   muscleMass.CreatedAt,
			UpdatedAt:   muscleMass.UpdatedAt,
			UserID:      muscleMass.UserID,
		})
	}

	bfPercents, err := cfg.db.GetAllBodyFatPercs(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error retrieving body fat percentages", err)
		return
	}
	for _, bfPercent := range bfPercents {
		bfPercentsResp = append(bfPercentsResp, Metric{
			ID:          bfPercent.ID,
			MetricType:  "body_fat_percentage",
			Measurement: bfPercent.Measurement,
			CreatedAt:   bfPercent.CreatedAt,
			UpdatedAt:   bfPercent.UpdatedAt,
			UserID:      userID,
		})
	}

	resp := map[string][]Metric{
		"body_weights":         bodyWeightsResp,
		"muscle_masses":        muscleMassesResp,
		"body_fat_percentages": bfPercentsResp,
	}
	respondWithJSON(w, http.StatusOK, resp)
}

func (cfg *apiConfig) updateMetricsHandler(w http.ResponseWriter, r *http.Request) {
	type requestParams struct {
		Measurement string `json:"measurement"`
	}
	metricType := r.PathValue("type")
	userID := r.Context().Value(userIDKey).(uuid.UUID)
	metricID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error parsing metric id", err)
		return
	}

	reqParams := requestParams{}
	if err := parseJSON(r, &reqParams); err != nil {
		respondWithError(w, http.StatusBadRequest, "malformed request", err)
		return
	}

	switch metricType {
	case "body_weights":
		bodyWeightEntry, err := cfg.db.UpdateBodyWeight(r.Context(), database.UpdateBodyWeightParams{
			ID:          metricID,
			UserID:      userID,
			Measurement: reqParams.Measurement,
		})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "error updating body weight entry", err)
			return
		}

		respondWithJSON(w, http.StatusAccepted, Metric{
			ID:          bodyWeightEntry.ID,
			MetricType:  "body_weight",
			Measurement: bodyWeightEntry.Measurement,
			CreatedAt:   bodyWeightEntry.CreatedAt,
			UpdatedAt:   bodyWeightEntry.UpdatedAt,
			UserID:      bodyWeightEntry.UserID,
		})
	case "muscle_masses":
		muscleMassEntry, err := cfg.db.UpdateMuscleMass(r.Context(), database.UpdateMuscleMassParams{
			ID:          metricID,
			UserID:      userID,
			Measurement: reqParams.Measurement,
		})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "error updating body weight entry", err)
			return
		}

		respondWithJSON(w, http.StatusAccepted, Metric{
			ID:          muscleMassEntry.ID,
			MetricType:  "body_weight",
			Measurement: muscleMassEntry.Measurement,
			CreatedAt:   muscleMassEntry.CreatedAt,
			UpdatedAt:   muscleMassEntry.UpdatedAt,
			UserID:      muscleMassEntry.UserID,
		})
	case "body_fat_percentages":
		bfPercentEntry, err := cfg.db.UpdateBodyFatPerc(r.Context(), database.UpdateBodyFatPercParams{
			ID:          metricID,
			UserID:      userID,
			Measurement: reqParams.Measurement,
		})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "error updating body weight entry", err)
			return
		}

		respondWithJSON(w, http.StatusAccepted, Metric{
			ID:          bfPercentEntry.ID,
			MetricType:  "body_fat_percentage",
			Measurement: bfPercentEntry.Measurement,
			CreatedAt:   bfPercentEntry.CreatedAt,
			UpdatedAt:   bfPercentEntry.UpdatedAt,
			UserID:      bfPercentEntry.UserID,
		})
	}
}

func (cfg *apiConfig) deleteMetricsHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIDKey).(uuid.UUID)
	metricType := r.PathValue("type")
	metricID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "error parsing metric id", err)
		return
	}

	switch metricType {
	case "body_weights":
		err := cfg.db.DeleteBodyWeight(r.Context(), database.DeleteBodyWeightParams{
			ID:     metricID,
			UserID: userID,
		})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "error deleting body weight entry", err)
			return
		}

		respondWithJSON(w, http.StatusNoContent, Metric{})
	case "muscle_masses":
		err := cfg.db.DeleteMuscleMass(r.Context(), database.DeleteMuscleMassParams{
			ID:     metricID,
			UserID: userID,
		})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "error deleting muscle mass entry", err)
			return
		}
		respondWithJSON(w, http.StatusNoContent, Metric{})
	case "body_fat_percentages":
		err := cfg.db.DeleteBodyFatPerc(r.Context(), database.DeleteBodyFatPercParams{
			ID:     metricID,
			UserID: userID,
		})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "error deleting body fat percentage entry", err)
			return
		}
		respondWithJSON(w, http.StatusNoContent, Metric{})
	}
}

func (cfg *apiConfig) deleteAllMetricsHandler(w http.ResponseWriter, r *http.Request) {
	metricType := r.PathValue("type")
	userID := r.Context().Value(userIDKey).(uuid.UUID)

	switch metricType {
	case "body_weights":
		if err := cfg.db.DeleteAllBodyWeights(r.Context(), userID); err != nil {
			respondWithError(w, http.StatusInternalServerError, "error deleting body weight entries", err)
			return
		}
		respondWithJSON(w, http.StatusNoContent, Metric{})
	case "muscle_masses":
		if err := cfg.db.DeleteAllMuscleMasses(r.Context(), userID); err != nil {
			respondWithError(w, http.StatusInternalServerError, "error deleting muscle mass entries", err)
			return
		}
		respondWithJSON(w, http.StatusNoContent, Metric{})
	case "body_fat_percentages":
		if err := cfg.db.DeleteAllBodyFatPercs(r.Context(), userID); err != nil {
			respondWithError(w, http.StatusInternalServerError, "error deleting body fat percentage entries", err)
			return
		}
		respondWithJSON(w, http.StatusNoContent, Metric{})
	}
}
