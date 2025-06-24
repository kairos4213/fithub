package handlers

import (
	"database/sql"
	"net/http"

	"github.com/google/uuid"
	"github.com/kairos4213/fithub/internal/cntx"
	"github.com/kairos4213/fithub/internal/database"
	"github.com/kairos4213/fithub/internal/templates"
)

func (h *Handler) GetAddExerciseForm(w http.ResponseWriter, r *http.Request) {
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
	templates.AdminLayout(contents, "FitHub-Admin | Exercises", true).Render(r.Context(), w)
}

func (h *Handler) AddExercise(w http.ResponseWriter, r *http.Request) {
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

	name := r.FormValue("exercise-name")

	exerciseDescription := r.FormValue("exercise-description")
	description := sql.NullString{Valid: false}
	if exerciseDescription != "" {
		description.String = exerciseDescription
		description.Valid = true
	}

	primaryMuscleGroup := r.FormValue("primary-muscle-group")
	primeMG := sql.NullString{Valid: false}
	if primaryMuscleGroup != "" {
		primeMG.String = primaryMuscleGroup
		primeMG.Valid = true
	}

	secondaryMuscleGroup := r.FormValue("secondary-muscle-group")
	secMG := sql.NullString{Valid: false}
	if secondaryMuscleGroup != "" {
		secMG.String = secondaryMuscleGroup
		secMG.Valid = true
	}

	_, err = h.DB.CreateExercise(r.Context(), database.CreateExerciseParams{
		Name:                 name,
		Description:          description,
		PrimaryMuscleGroup:   primeMG,
		SecondaryMuscleGroup: secMG,
	})
	if err != nil {
		return // TODO: handle err
	}

	w.Header().Set("Content-type", "text/html")
	w.Header().Set("HX-Location", `{"path": "/admin/exercises"}`)
	w.WriteHeader(http.StatusCreated)
}
