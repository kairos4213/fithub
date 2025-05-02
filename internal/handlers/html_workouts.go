package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/kairos4213/fithub/internal/cntx"
	"github.com/kairos4213/fithub/internal/templates"
)

func (h *Handler) GetUserWorkouts(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(cntx.UserIDKey).(uuid.UUID)
	workouts, err := h.DB.GetAllUserWorkouts(r.Context(), userID)
	if err != nil {
		http.Error(w, "error fetching workouts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	contents := templates.Workouts(workouts)
	templates.Layout(contents, "Fithub | Workouts", true).Render(r.Context(), w)
}

func (h *Handler) NewUserWorkout(w http.ResponseWriter, r *http.Request) {
	// userID := r.Context().Value(cntx.UserIDKey).(uuid.UUID)

	if r.Method == "GET" {
		contents := templates.CreateWorkout()
		templates.Layout(contents, "FitHub | New Workout", true).Render(r.Context(), w)
		return
	}
}
