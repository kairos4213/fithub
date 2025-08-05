package handlers

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/kairos4213/fithub/internal/cntx"
	"github.com/kairos4213/fithub/internal/templates"
)

func (h *Handler) GetAdminHome(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(cntx.UserIDKey).(uuid.UUID)
	user, err := h.DB.GetUserByID(r.Context(), userID)
	if err != nil {
		w.Header().Set("Content-type", "text/html")
		w.WriteHeader(http.StatusInternalServerError)

		htmlErr := templates.HtmlErr{Code: http.StatusInternalServerError, Msg: "Something went wrong. Please try later"}
		contents := templates.ErrorDisplay(htmlErr)
		templates.Layout(contents, "FitHub", false).Render(r.Context(), w)

		log.Printf("Server Error: %v", err)
		return
	}

	if !user.IsAdmin {
		w.Header().Set("Content-type", "text/html")
		w.WriteHeader(http.StatusForbidden)

		htmlErr := templates.HtmlErr{Code: http.StatusForbidden, Msg: "You don't have permission to access this resource"}
		contents := templates.ErrorDisplay(htmlErr)
		templates.Layout(contents, "FitHub", true).Render(r.Context(), w)

		log.Println("Unauthorized admin GET request:")
		log.Printf("\tUser ID: %v", user.ID)
		log.Printf("\tUser Email: %v", user.Email)
		return
	}

	exercises, err := h.DB.GetAllExercises(r.Context())
	if err != nil {
		w.Header().Set("Content-type", "text/html")
		w.WriteHeader(http.StatusInternalServerError)

		htmlErr := templates.HtmlErr{Code: http.StatusInternalServerError, Msg: "Something went wrong. Please try later"}
		contents := templates.ErrorDisplay(htmlErr)
		templates.Layout(contents, "FitHub", false).Render(r.Context(), w)

		log.Printf("Server Error: %v", err)
		return
	}

	contents := templates.AdminExercisesPage(exercises)
	templates.AdminLayout(contents, "FitHub-Admin | Home", true).Render(r.Context(), w)
}
