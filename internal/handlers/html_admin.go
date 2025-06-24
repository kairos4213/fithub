package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/kairos4213/fithub/internal/cntx"
	"github.com/kairos4213/fithub/internal/templates"
)

func (h *Handler) GetAdminHome(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(cntx.UserIDKey).(uuid.UUID)
	user, err := h.DB.GetUserByID(r.Context(), userID)
	if err != nil {
		return // TODO: Handle err
	}

	if !user.IsAdmin {
		w.Header().Set("Content-type", "text/html")
		w.WriteHeader(http.StatusUnauthorized)
		return // TODO: Finish error handler
	}

	contents := templates.AddExerciseForm()
	templates.AdminLayout(contents, "FitHub-Admin | Home", true).Render(r.Context(), w)
}
