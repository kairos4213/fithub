package handlers

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/kairos4213/fithub/internal/cntx"
	"github.com/kairos4213/fithub/internal/database"
	"github.com/kairos4213/fithub/internal/templates"
)

func (h *Handler) GetAllMetrics(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(cntx.UserIDKey).(uuid.UUID)

	bodyweights, err := h.DB.GetAllBodyWeights(r.Context(), userID)
	if err != nil {
		log.Printf("error: %v", err)
		// TODO: html error response
		return
	}
	muscleMasses, err := h.DB.GetAllMuscleMasses(r.Context(), userID)
	if err != nil {
		// TODO: html error response
		return
	}

	bfPercents, err := h.DB.GetAllBodyFatPercs(r.Context(), userID)
	if err != nil {
		// TODO: html error response
		return
	}

	contents := templates.Metrics(bodyweights, muscleMasses, bfPercents)
	templates.Layout(contents, "Fithub | Metrics", true).Render(r.Context(), w)
}

func (h *Handler) LogMetrics(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(cntx.UserIDKey).(uuid.UUID)

	metricType := r.PathValue("type")
	switch metricType {
	case "bodyweights":
		entry := r.FormValue("bodyweight")
		bw, err := h.DB.AddBodyWeight(r.Context(), database.AddBodyWeightParams{UserID: userID, Measurement: entry})
		if err != nil {
			// TODO: handle error
			return
		}

		templates.BWDataRow(bw).Render(r.Context(), w)
	case "muscleMasses":
		entry := r.FormValue("muscleMass")
		_, err := h.DB.AddMuscleMass(r.Context(), database.AddMuscleMassParams{UserID: userID, Measurement: entry})
		if err != nil {
			return
		}

		muscleMasses, err := h.DB.GetAllMuscleMasses(r.Context(), userID)
		if err != nil {
			return
		}
		templates.MuscleMassesSect(muscleMasses).Render(r.Context(), w)
	case "bfPercents":
		entry := r.FormValue("bfPercent")
		_, err := h.DB.AddBodyFatPerc(r.Context(), database.AddBodyFatPercParams{UserID: userID, Measurement: entry})
		if err != nil {
			return
		}

		bfPercents, err := h.DB.GetAllBodyFatPercs(r.Context(), userID)
		if err != nil {
			return
		}
		templates.BfPercentsSect(bfPercents).Render(r.Context(), w)
	}
}

func (h *Handler) GetEditMetricsForm(w http.ResponseWriter, r *http.Request) {
	metricType := r.PathValue("type")
	id := r.PathValue("id")
	switch metricType {
	case "bodyweights":
		entry := r.FormValue("bw-entry")
		date := r.FormValue("bw-date")
		templates.EditBWForm(entry, date, id).Render(r.Context(), w)
	}
}

func (h *Handler) EditMetrics(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(cntx.UserIDKey).(uuid.UUID)

	metricType := r.PathValue("type")
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		// TODO: send error
		return
	}
	switch metricType {
	case "bodyweights":
		entry := r.FormValue("bodyweight")

		updatedBW, err := h.DB.UpdateBodyWeight(r.Context(), database.UpdateBodyWeightParams{Measurement: entry, ID: id, UserID: userID})
		if err != nil {
			return // TODO: send error
		}

		templates.BWDataRow(updatedBW).Render(r.Context(), w)
	}
}

func (h *Handler) DeleteMetrics(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(cntx.UserIDKey).(uuid.UUID)

	metricType := r.PathValue("type")
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		// TODO: send error
		return
	}
	switch metricType {
	case "bodyweights":
		err := h.DB.DeleteBodyWeight(r.Context(), database.DeleteBodyWeightParams{ID: id, UserID: userID})
		if err != nil {
			return // TODO: send error
		}

		w.WriteHeader(http.StatusOK)
	}
}
