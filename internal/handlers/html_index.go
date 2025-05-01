package handlers

import (
	"net/http"

	"github.com/kairos4213/fithub/internal/templates"
)

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("access_token")
	if err != nil {
		contents := templates.Index()
		templates.Layout(contents, "FitHub", false).Render(r.Context(), w)
		return
	}

	w.Header().Set("Content-type", "text/html")
	w.Header().Set("HX-Location", `{"path": "/workouts"}`)
	w.WriteHeader(http.StatusFound)
}
