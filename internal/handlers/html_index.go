package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/kairos4213/fithub/internal/cntx"
	"github.com/kairos4213/fithub/internal/templates"
)

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value(cntx.UserIDKey).(uuid.UUID)

	if ok {
		http.Redirect(w, r, "/home", http.StatusFound)
	}

	contents := templates.Index()
	templates.Layout(contents, "FitHub").Render(r.Context(), w)
}
