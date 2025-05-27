package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/kairos4213/fithub/internal/cntx"
	"github.com/kairos4213/fithub/internal/templates"
)

func (h *Handler) GetAllGoals(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(cntx.UserIDKey).(uuid.UUID)
	goals, err := h.DB.GetAllUserGoals(r.Context(), userID)
	if err != nil {
		return // TODO: handle error
	}

	contents := templates.Goals(goals)
	templates.Layout(contents, "Fithub | Goals", true).Render(r.Context(), w)
}
