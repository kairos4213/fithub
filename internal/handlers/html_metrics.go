package handlers

import (
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/kairos4213/fithub/internal/cntx"
	"github.com/kairos4213/fithub/internal/database"
	"github.com/kairos4213/fithub/internal/templates"
)

func (h *Handler) GetAllMetrics(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(cntx.UserIDKey).(uuid.UUID)

	bodyweights, err := h.cfg.DB.GetAllBodyWeights(r.Context(), userID)
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to get body weights", slog.String("error", err.Error()))
		return
	}
	muscleMasses, err := h.cfg.DB.GetAllMuscleMasses(r.Context(), userID)
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to get muscle masses", slog.String("error", err.Error()))
		return
	}

	bfPercents, err := h.cfg.DB.GetAllBodyFatPercs(r.Context(), userID)
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to get body fats", slog.String("error", err.Error()))
		return
	}

	contents := templates.Metrics(bodyweights, muscleMasses, bfPercents)
	err = templates.Layout(contents, "Fithub | Metrics", true).Render(r.Context(), w)
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to render metrics page", slog.String("error", err.Error()))
		return
	}
}

func (h *Handler) LogMetrics(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(cntx.UserIDKey).(uuid.UUID)

	metricType := r.PathValue("type")
	switch metricType {
	case "bodyweights":
		entry := r.FormValue("bodyweight")
		bw, err := h.cfg.DB.AddBodyWeight(r.Context(), database.AddBodyWeightParams{UserID: userID, Measurement: entry})
		if err != nil {
			HandleInternalServerError(w, r)
			h.cfg.Logger.Error("failed to add body weight", slog.String("error", err.Error()))
			return
		}

		err = templates.BWDataRow(bw).Render(r.Context(), w)
		if err != nil {
			HandleInternalServerError(w, r)
			h.cfg.Logger.Error("failed to render body weights", slog.String("error", err.Error()))
			return
		}
	case "muscleMasses":
		entry := r.FormValue("muscle-mass")
		mm, err := h.cfg.DB.AddMuscleMass(r.Context(), database.AddMuscleMassParams{UserID: userID, Measurement: entry})
		if err != nil {
			HandleInternalServerError(w, r)
			h.cfg.Logger.Error("failed to add muscle mass", slog.String("error", err.Error()))
			return
		}

		err = templates.MMDataRow(mm).Render(r.Context(), w)
		if err != nil {
			HandleInternalServerError(w, r)
			h.cfg.Logger.Error("failed to render muscle masses", slog.String("error", err.Error()))
			return
		}
	case "bfPercents":
		entry := r.FormValue("bf-percent")
		bf, err := h.cfg.DB.AddBodyFatPerc(r.Context(), database.AddBodyFatPercParams{UserID: userID, Measurement: entry})
		if err != nil {
			HandleInternalServerError(w, r)
			h.cfg.Logger.Error("failed to add body fat", slog.String("error", err.Error()))
			return
		}

		err = templates.BFPercentDataRow(bf).Render(r.Context(), w)
		if err != nil {
			HandleInternalServerError(w, r)
			h.cfg.Logger.Error("failed to render body fats", slog.String("error", err.Error()))
			return
		}
	}
}

func (h *Handler) EditMetrics(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(cntx.UserIDKey).(uuid.UUID)

	metricType := r.PathValue("type")
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to parse metric id", slog.String("error", err.Error()))
		return
	}
	switch metricType {
	case "bodyweights":
		entry := r.FormValue("bodyweight")

		updatedBW, err := h.cfg.DB.UpdateBodyWeight(r.Context(), database.UpdateBodyWeightParams{Measurement: entry, ID: id, UserID: userID})
		if err != nil {
			HandleInternalServerError(w, r)
			h.cfg.Logger.Error("failed to update body weight", slog.String("error", err.Error()))
			return
		}

		err = templates.BWDataRow(updatedBW).Render(r.Context(), w)
		if err != nil {
			HandleInternalServerError(w, r)
			h.cfg.Logger.Error("failed to render body weight", slog.String("error", err.Error()))
			return
		}
	case "muscleMasses":
		entry := r.FormValue("muscle-mass")

		updatedMM, err := h.cfg.DB.UpdateMuscleMass(r.Context(), database.UpdateMuscleMassParams{Measurement: entry, ID: id, UserID: userID})
		if err != nil {
			HandleInternalServerError(w, r)
			h.cfg.Logger.Error("failed to update muscle mass", slog.String("error", err.Error()))
			return
		}

		err = templates.MMDataRow(updatedMM).Render(r.Context(), w)
		if err != nil {
			HandleInternalServerError(w, r)
			h.cfg.Logger.Error("failed to render muscle mass", slog.String("error", err.Error()))
			return
		}
	case "bfPercents":
		entry := r.FormValue("bf-percent")

		updatedBF, err := h.cfg.DB.UpdateBodyFatPerc(r.Context(), database.UpdateBodyFatPercParams{Measurement: entry, ID: id, UserID: userID})
		if err != nil {
			HandleInternalServerError(w, r)
			h.cfg.Logger.Error("failed to update body fat", slog.String("error", err.Error()))
			return
		}

		err = templates.BFPercentDataRow(updatedBF).Render(r.Context(), w)
		if err != nil {
			HandleInternalServerError(w, r)
			h.cfg.Logger.Error("failed to render body fat", slog.String("error", err.Error()))
			return
		}
	}
}

func (h *Handler) DeleteMetrics(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(cntx.UserIDKey).(uuid.UUID)

	metricType := r.PathValue("type")
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to parse metric id", slog.String("error", err.Error()))
		return
	}
	switch metricType {
	case "bodyweights":
		err := h.cfg.DB.DeleteBodyWeight(r.Context(), database.DeleteBodyWeightParams{ID: id, UserID: userID})
		if err != nil {
			HandleInternalServerError(w, r)
			h.cfg.Logger.Error("failed to delete body weight", slog.String("error", err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
	case "muscleMasses":
		err := h.cfg.DB.DeleteMuscleMass(r.Context(), database.DeleteMuscleMassParams{ID: id, UserID: userID})
		if err != nil {
			HandleInternalServerError(w, r)
			h.cfg.Logger.Error("failed to delete muscle mass", slog.String("error", err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
	case "bfPercents":
		err := h.cfg.DB.DeleteBodyFatPerc(r.Context(), database.DeleteBodyFatPercParams{ID: id, UserID: userID})
		if err != nil {
			HandleInternalServerError(w, r)
			h.cfg.Logger.Error("failed to delete body fat", slog.String("error", err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
