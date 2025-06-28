package handlers

import (
	"database/sql"
	"net/http"

	"github.com/google/uuid"
	"github.com/kairos4213/fithub/internal/cntx"
	"github.com/kairos4213/fithub/internal/database"
	"github.com/kairos4213/fithub/internal/templates"
)

func (h *Handler) GetAdminExercisesPage(w http.ResponseWriter, r *http.Request) {
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

	exercises, err := h.DB.GetAllExercises(r.Context())
	if err != nil {
		return // TODO: handle err
	}

	contents := templates.AdminExercisesPage(exercises)
	templates.AdminLayout(contents, "FitHub-Admin | Home", true).Render(r.Context(), w)
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

func (h *Handler) EditExercise(w http.ResponseWriter, r *http.Request) {
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

	exerciseID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		return // TODO: Handle error
	}

	err = r.ParseForm()
	if err != nil {
		return // TODO: handle error
	}

	reqName := r.FormValue("exercise-name")
	reqDescription := r.FormValue("exercise-description")
	reqPrimaryMG := r.FormValue("primary-muscle-group")
	reqSecondaryMG := r.FormValue("secondary-muscle-group")

	description := sql.NullString{Valid: false}
	if reqDescription != "" {
		description.String = reqDescription
		description.Valid = true
	}

	primaryMG := sql.NullString{Valid: false}
	if reqPrimaryMG != "" {
		primaryMG.String = reqPrimaryMG
		primaryMG.Valid = true
	}

	secondaryMG := sql.NullString{Valid: false}
	if reqSecondaryMG != "" {
		secondaryMG.String = reqSecondaryMG
		secondaryMG.Valid = true
	}

	updatedExercise, err := h.DB.UpdateExercise(r.Context(), database.UpdateExerciseParams{
		Name:                 reqName,
		Description:          description,
		PrimaryMuscleGroup:   primaryMG,
		SecondaryMuscleGroup: secondaryMG,
		ID:                   exerciseID,
	})
	if err != nil {
		return // TODO: handle err
	}

	templates.ExerciseDataRow(updatedExercise).Render(r.Context(), w)
}

func (h *Handler) DeleteExercise(w http.ResponseWriter, r *http.Request) {
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

	exerciseID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		return // TODO: handle err
	}

	err = h.DB.DeleteExercise(r.Context(), exerciseID)
	if err != nil {
		return // TODO: handle err
	}

	w.WriteHeader(http.StatusOK)
}
