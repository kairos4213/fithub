package handlers

import (
	"log/slog"
	"net/http"

	"github.com/kairos4213/fithub/internal/templates"
)

func (h *Handler) GetAllWorkoutTemplates(w http.ResponseWriter, r *http.Request) {
	workoutTemplates, err := h.cfg.DB.GetAllWorkoutTemplates(r.Context())
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to get workout templates", slog.String("error", err.Error()))
		return
	}

	err = templates.Layout(templates.WorkoutTemplatesPage(workoutTemplates), "FitHub | Templates", true).Render(r.Context(), w)
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to render workout templates", slog.String("error", err.Error()))
		return
	}
}
