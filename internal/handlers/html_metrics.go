package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/kairos4213/fithub/internal/cntx"
	"github.com/kairos4213/fithub/internal/templates"
)

func (h *Handler) GetAllMetrics(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(cntx.UserIDKey).(uuid.UUID)

	bodyweights, err := h.DB.GetAllBodyWeights(r.Context(), userID)
	if err != nil {
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
