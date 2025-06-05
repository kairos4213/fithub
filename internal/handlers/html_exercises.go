package handlers

import (
	"net/http"

	"github.com/kairos4213/fithub/internal/templates"
)

func (h *Handler) GetAddExerciseForm(w http.ResponseWriter, r *http.Request) {
	contents := templates.AddExerciseForm()
	templates.Layout(contents, "FitHub-Admin | Exercises", true).Render(r.Context(), w)
}
