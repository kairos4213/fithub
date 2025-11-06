package handlers

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/kairos4213/fithub/internal/cntx"
	"github.com/kairos4213/fithub/internal/database"
	"github.com/kairos4213/fithub/internal/utils"
)

type Metric struct {
	ID          string `json:"metric_id,omitempty"`
	MetricType  string `json:"metric_type,omitempty"`
	Measurement string `json:"measurement,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
	UserID      string `json:"user_id,omitempty"`
}

func (h *Handler) AddMetric(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(cntx.UserIDKey).(uuid.UUID)

	reqParams := Metric{}
	if err := utils.ParseJSON(r, &reqParams); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "malformed request", err)
		return
	}

	metricType := r.PathValue("type")
	switch metricType {
	case "body_weights":
		bodyWeightEntry, err := h.cfg.DB.AddBodyWeight(r.Context(), database.AddBodyWeightParams{
			UserID:      userID,
			Measurement: reqParams.Measurement,
		})
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error saving body weight metric", err)
			return
		}

		utils.RespondWithJSON(w, http.StatusCreated, Metric{
			ID:          bodyWeightEntry.ID.String(),
			MetricType:  "body_weight",
			Measurement: bodyWeightEntry.Measurement,
			CreatedAt:   bodyWeightEntry.CreatedAt.Format(time.RFC822),
			UpdatedAt:   bodyWeightEntry.UpdatedAt.Format(time.RFC822),
			UserID:      userID.String(),
		})
	case "muscle_masses":
		muscleMassEntry, err := h.cfg.DB.AddMuscleMass(r.Context(), database.AddMuscleMassParams{
			UserID:      userID,
			Measurement: reqParams.Measurement,
		})
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error saving muscle mass metric", err)
			return
		}

		utils.RespondWithJSON(w, http.StatusCreated, Metric{
			ID:          muscleMassEntry.ID.String(),
			MetricType:  "muscle_mass",
			Measurement: muscleMassEntry.Measurement,
			CreatedAt:   muscleMassEntry.CreatedAt.Format(time.RFC822),
			UpdatedAt:   muscleMassEntry.UpdatedAt.Format(time.RFC822),
			UserID:      userID.String(),
		})
	case "body_fat_percents":
		bfPercentEntry, err := h.cfg.DB.AddBodyFatPerc(r.Context(), database.AddBodyFatPercParams{
			UserID:      userID,
			Measurement: reqParams.Measurement,
		})
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error saving body fat percentage", err)
			return
		}

		utils.RespondWithJSON(w, http.StatusCreated, Metric{
			ID:          bfPercentEntry.ID.String(),
			MetricType:  "body_fat_percentage",
			Measurement: bfPercentEntry.Measurement,
			CreatedAt:   bfPercentEntry.CreatedAt.Format(time.RFC822),
			UpdatedAt:   bfPercentEntry.UpdatedAt.Format(time.RFC822),
			UserID:      userID.String(),
		})
	}
}

func (h *Handler) GetAllUserMetrics(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(cntx.UserIDKey).(uuid.UUID)

	bodyWeightsResp := []Metric{}
	muscleMassesResp := []Metric{}
	bfPercentsResp := []Metric{}

	bodyWeights, err := h.cfg.DB.GetAllBodyWeights(r.Context(), userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error retrieving body weights", err)
		return
	}
	for _, bodyWeight := range bodyWeights {
		bodyWeightsResp = append(bodyWeightsResp, Metric{
			ID:          bodyWeight.ID.String(),
			MetricType:  "body_weight",
			Measurement: bodyWeight.Measurement,
			CreatedAt:   bodyWeight.CreatedAt.Format(time.RFC822),
			UpdatedAt:   bodyWeight.UpdatedAt.Format(time.RFC822),
			UserID:      userID.String(),
		})
	}

	muscleMasses, err := h.cfg.DB.GetAllMuscleMasses(r.Context(), userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error retrieving muscle mass metrics", err)
		return
	}
	for _, muscleMass := range muscleMasses {
		muscleMassesResp = append(muscleMassesResp, Metric{
			ID:          muscleMass.ID.String(),
			MetricType:  "muscle_mass",
			Measurement: muscleMass.Measurement,
			CreatedAt:   muscleMass.CreatedAt.Format(time.RFC822),
			UpdatedAt:   muscleMass.UpdatedAt.Format(time.RFC822),
			UserID:      userID.String(),
		})
	}

	bfPercents, err := h.cfg.DB.GetAllBodyFatPercs(r.Context(), userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error retrieving body fat percentages", err)
		return
	}
	for _, bfPercent := range bfPercents {
		bfPercentsResp = append(bfPercentsResp, Metric{
			ID:          bfPercent.ID.String(),
			MetricType:  "body_fat_percentage",
			Measurement: bfPercent.Measurement,
			CreatedAt:   bfPercent.CreatedAt.Format(time.RFC822),
			UpdatedAt:   bfPercent.UpdatedAt.Format(time.RFC822),
			UserID:      userID.String(),
		})
	}

	resp := map[string][]Metric{
		"body_weights":         bodyWeightsResp,
		"muscle_masses":        muscleMassesResp,
		"body_fat_percentages": bfPercentsResp,
	}
	utils.RespondWithJSON(w, http.StatusOK, resp)
}

func (h *Handler) UpdateMetric(w http.ResponseWriter, r *http.Request) {
	metricType := r.PathValue("type")
	userID := r.Context().Value(cntx.UserIDKey).(uuid.UUID)
	metricID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "error parsing metric id", err)
		return
	}

	reqParams := Metric{}
	if err := utils.ParseJSON(r, &reqParams); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "malformed request", err)
		return
	}

	switch metricType {
	case "body_weights":
		bodyWeightEntry, err := h.cfg.DB.UpdateBodyWeight(r.Context(), database.UpdateBodyWeightParams{
			ID:          metricID,
			UserID:      userID,
			Measurement: reqParams.Measurement,
		})
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "error updating body weight entry", err)
			return
		}

		utils.RespondWithJSON(w, http.StatusAccepted, Metric{
			ID:          bodyWeightEntry.ID.String(),
			MetricType:  "body_weight",
			Measurement: bodyWeightEntry.Measurement,
			CreatedAt:   bodyWeightEntry.CreatedAt.Format(time.RFC822),
			UpdatedAt:   bodyWeightEntry.UpdatedAt.Format(time.RFC822),
			UserID:      userID.String(),
		})
	case "muscle_masses":
		muscleMassEntry, err := h.cfg.DB.UpdateMuscleMass(r.Context(), database.UpdateMuscleMassParams{
			ID:          metricID,
			UserID:      userID,
			Measurement: reqParams.Measurement,
		})
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "error updating body weight entry", err)
			return
		}

		utils.RespondWithJSON(w, http.StatusAccepted, Metric{
			ID:          muscleMassEntry.ID.String(),
			MetricType:  "body_weight",
			Measurement: muscleMassEntry.Measurement,
			CreatedAt:   muscleMassEntry.CreatedAt.Format(time.RFC822),
			UpdatedAt:   muscleMassEntry.UpdatedAt.Format(time.RFC822),
			UserID:      userID.String(),
		})
	case "body_fat_percentages":
		bfPercentEntry, err := h.cfg.DB.UpdateBodyFatPerc(r.Context(), database.UpdateBodyFatPercParams{
			ID:          metricID,
			UserID:      userID,
			Measurement: reqParams.Measurement,
		})
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "error updating body weight entry", err)
			return
		}

		utils.RespondWithJSON(w, http.StatusAccepted, Metric{
			ID:          bfPercentEntry.ID.String(),
			MetricType:  "body_fat_percentage",
			Measurement: bfPercentEntry.Measurement,
			CreatedAt:   bfPercentEntry.CreatedAt.Format(time.RFC822),
			UpdatedAt:   bfPercentEntry.UpdatedAt.Format(time.RFC822),
			UserID:      bfPercentEntry.UserID.String(),
		})
	}
}

func (h *Handler) DeleteMetric(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(cntx.UserIDKey).(uuid.UUID)
	metricType := r.PathValue("type")
	metricID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "error parsing metric id", err)
		return
	}

	switch metricType {
	case "body_weights":
		err := h.cfg.DB.DeleteBodyWeight(r.Context(), database.DeleteBodyWeightParams{
			ID:     metricID,
			UserID: userID,
		})
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "error deleting body weight entry", err)
			return
		}

		utils.RespondWithJSON(w, http.StatusNoContent, Metric{})
	case "muscle_masses":
		err := h.cfg.DB.DeleteMuscleMass(r.Context(), database.DeleteMuscleMassParams{
			ID:     metricID,
			UserID: userID,
		})
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "error deleting muscle mass entry", err)
			return
		}
		utils.RespondWithJSON(w, http.StatusNoContent, Metric{})
	case "body_fat_percentages":
		err := h.cfg.DB.DeleteBodyFatPerc(r.Context(), database.DeleteBodyFatPercParams{
			ID:     metricID,
			UserID: userID,
		})
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "error deleting body fat percentage entry", err)
			return
		}
		utils.RespondWithJSON(w, http.StatusNoContent, Metric{})
	}
}

func (h *Handler) DeleteAllUserMetrics(w http.ResponseWriter, r *http.Request) {
	metricType := r.PathValue("type")
	userID := r.Context().Value(cntx.UserIDKey).(uuid.UUID)

	switch metricType {
	case "body_weights":
		if err := h.cfg.DB.DeleteAllBodyWeights(r.Context(), userID); err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "error deleting body weight entries", err)
			return
		}
		utils.RespondWithJSON(w, http.StatusNoContent, Metric{})
	case "muscle_masses":
		if err := h.cfg.DB.DeleteAllMuscleMasses(r.Context(), userID); err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "error deleting muscle mass entries", err)
			return
		}
		utils.RespondWithJSON(w, http.StatusNoContent, Metric{})
	case "body_fat_percentages":
		if err := h.cfg.DB.DeleteAllBodyFatPercs(r.Context(), userID); err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "error deleting body fat percentage entries", err)
			return
		}
		utils.RespondWithJSON(w, http.StatusNoContent, Metric{})
	}
}
