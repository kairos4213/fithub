package handlers

import (
	"log/slog"
	"net/http"

	"github.com/kairos4213/fithub/internal/templates"
)

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("access_token")
	if err != nil {
		contents := templates.Index()
		err = templates.Layout(contents, "FitHub", false).Render(r.Context(), w)
		if err != nil {
			HandleInternalServerError(w, r)
			h.cfg.Logger.Error("failed to render index page", slog.String("error", err.Error()))
			return
		}
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("Content-type", "text/html")
		w.Header().Set("HX-Location", `{"path": "/workouts"}`)
		w.WriteHeader(http.StatusFound)
	} else {
		http.Redirect(w, r, "/workouts", http.StatusSeeOther)
	}
}
