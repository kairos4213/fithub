package handlers

import (
	"log/slog"
	"net/http"

	"github.com/kairos4213/fithub/internal/templates"
)

func (h *Handler) GetAdminHome(w http.ResponseWriter, r *http.Request) {
	exercises, err := h.cfg.DB.GetAllExercises(r.Context())
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to fetch exercises", slog.String("error", err.Error()))
		return
	}

	contents := templates.AdminExercisesPage(exercises)
	err = templates.AdminLayout(contents, "FitHub-Admin | Home", true).Render(r.Context(), w)
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to render admin home", slog.String("error", err.Error()))
		return
	}
}
