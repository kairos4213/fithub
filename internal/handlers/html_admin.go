package handlers

import (
	"log"
	"net/http"

	"github.com/kairos4213/fithub/internal/templates"
)

func (h *Handler) GetAdminHome(w http.ResponseWriter, r *http.Request) {
	exercises, err := h.DB.GetAllExercises(r.Context())
	if err != nil {
		HandleInternalServerError(w, r)
		log.Printf("Server Error: %v", err)
		return
	}

	contents := templates.AdminExercisesPage(exercises)
	templates.AdminLayout(contents, "FitHub-Admin | Home", true).Render(r.Context(), w)
}
