package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/kairos4213/fithub/internal/cntx"
	"github.com/kairos4213/fithub/internal/templates"
)

func (h *Handler) Workouts(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(cntx.UserIDKey).(uuid.UUID)

	if !ok {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	workouts, err := h.DB.GetAllUserWorkouts(r.Context(), userID)
	if err != nil {
		http.Error(w, "Error fetching workouts", http.StatusInternalServerError)
		return
	}

	contents := templates.Workouts(workouts)
	err = templates.Layout(contents, "FitHub", true).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
