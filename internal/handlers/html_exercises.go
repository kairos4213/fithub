package handlers

import (
	"database/sql"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/kairos4213/fithub/internal/database"
	"github.com/kairos4213/fithub/internal/templates"
)

func (h *Handler) GetAdminExercisesPage(w http.ResponseWriter, r *http.Request) {
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
		h.cfg.Logger.Error("failed to render admin exercises page", slog.String("error", err.Error()))
		return
	}
}

func (h *Handler) AddDBExercise(w http.ResponseWriter, r *http.Request) {
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

	_, err := h.cfg.DB.CreateExercise(r.Context(), database.CreateExerciseParams{
		Name:                 name,
		Description:          description,
		PrimaryMuscleGroup:   primeMG,
		SecondaryMuscleGroup: secMG,
	})
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to create exercise", slog.String("error", err.Error()))
		return
	}

	w.Header().Set("Content-type", "text/html")
	w.Header().Set("HX-Location", `{"path": "/admin/exercises"}`)
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) EditDBExercise(w http.ResponseWriter, r *http.Request) {
	exerciseID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to parse exercise id", slog.String("error", err.Error()))
		return
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

	updatedExercise, err := h.cfg.DB.UpdateExercise(r.Context(), database.UpdateExerciseParams{
		Name:                 reqName,
		Description:          description,
		PrimaryMuscleGroup:   primaryMG,
		SecondaryMuscleGroup: secondaryMG,
		ID:                   exerciseID,
	})
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to update exercise", slog.String("error", err.Error()))
		return
	}

	err = templates.ExerciseDataRow(updatedExercise).Render(r.Context(), w)
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to render exercise row", slog.String("error", err.Error()))
		return
	}
}

func (h *Handler) DeleteDBExercise(w http.ResponseWriter, r *http.Request) {
	exerciseID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to parse exercise id", slog.String("error", err.Error()))
		return
	}

	err = h.cfg.DB.DeleteExercise(r.Context(), exerciseID)
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to delete exercise", slog.String("error", err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetExerciseByKeyword(w http.ResponseWriter, r *http.Request) {
	exerciseSearch := strings.ToLower(r.FormValue("exercise-search"))
	workoutID, err := uuid.Parse(r.FormValue("workoutID"))
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to parse workout id", slog.String("error", err.Error()))
		return
	}

	if exerciseSearch == "" {
		templates.ExerciseQuickSearchResults([]database.Exercise{}, workoutID).Render(r.Context(), w)
		return
	}

	exercises, err := h.cfg.DB.GetExerciseByKeyword(r.Context(), exerciseSearch)
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to fetch searched exercise", slog.String("error", err.Error()))
		return
	}

	err = templates.ExerciseQuickSearchResults(exercises, workoutID).Render(r.Context(), w)
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to render exercise search results", slog.String("error", err.Error()))
		return
	}
}

func (h *Handler) AddExerciseToWorkout(w http.ResponseWriter, r *http.Request) {
	workoutID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to parse workout id", slog.String("error", err.Error()))
		return
	}

	exerciseName := r.FormValue("exercise-name")
	exercise, err := h.cfg.DB.GetExerciseByName(r.Context(), exerciseName)
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to fetch exercise by name", slog.String("error", err.Error()))
		return
	}

	plannedSets, err := strconv.ParseInt(r.FormValue("planned-sets"), 10, 32)
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to parse planned sets", slog.String("error", err.Error()))
		return
	}

	plannedRepsSlice := r.PostForm["planned-reps[]"]
	plannedReps := make([]int32, len(plannedRepsSlice))
	for i, plannedRep := range plannedRepsSlice {
		rep, err := strconv.ParseInt(plannedRep, 10, 32)
		if err != nil {
			HandleInternalServerError(w, r)
			h.cfg.Logger.Error("failed to parse planned rep", slog.String("error", err.Error()))
			return
		}
		plannedReps[i] = int32(rep)
	}

	plannedWeightsSlice := r.PostForm["planned-weights[]"]
	plannedWeights := make([]int32, len(plannedWeightsSlice))
	for i, plannedWeight := range plannedWeightsSlice {
		weight, err := strconv.ParseInt(plannedWeight, 10, 32)
		if err != nil {
			HandleInternalServerError(w, r)
			h.cfg.Logger.Error("failed to parse planned weight", slog.String("error", err.Error()))
			return
		}
		plannedWeights[i] = int32(weight)
	}

	workoutExercise, err := h.cfg.DB.AddExerciseToWorkout(r.Context(), database.AddExerciseToWorkoutParams{
		WorkoutID:           workoutID,
		ExerciseID:          exercise.ID,
		SetsPlanned:         int32(plannedSets),
		RepsPerSetPlanned:   plannedReps,
		SetsCompleted:       0,
		RepsPerSetCompleted: []int32{},
		WeightsPlannedLbs:   plannedWeights,
		WeightsCompletedLbs: []int32{},
	})
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to add exercise to workout", slog.String("error", err.Error()))
		return
	}

	exerciseForWorkout := database.WorkoutAndExercisesRow{
		WorkoutsExercise: workoutExercise,
		Exercise:         exercise,
	}

	err = templates.WorkoutExercisesTableDataRow(exerciseForWorkout).Render(r.Context(), w)
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to render workout exercises table row", slog.String("error", err.Error()))
		return
	}
}

func (h *Handler) UpdateWorkoutExercise(w http.ResponseWriter, r *http.Request) {
	workoutID, err := uuid.Parse(r.PathValue("workoutID"))
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to parse workout id", slog.String("error", err.Error()))
		return
	}
	workoutExerciseID, err := uuid.Parse(r.PathValue("workoutExerciseID"))
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to parse workout exercise id", slog.String("error", err.Error()))
		return
	}

	exerciseID, err := uuid.Parse(r.FormValue("exercise"))
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to parse exercise id", slog.String("error", err.Error()))
		return
	}

	exercise, err := h.cfg.DB.GetExerciseByID(r.Context(), exerciseID)
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to fetch exercise", slog.String("error", err.Error()))
		return
	}

	plannedSets, err := strconv.ParseInt(r.FormValue("planned-sets"), 10, 32)
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to parse planned sets", slog.String("error", err.Error()))
		return
	}

	plannedRepsSlice := r.PostForm["planned-reps[]"]
	plannedReps := make([]int32, len(plannedRepsSlice))
	for i, plannedRep := range plannedRepsSlice {
		rep, err := strconv.ParseInt(plannedRep, 10, 32)
		if err != nil {
			HandleInternalServerError(w, r)
			h.cfg.Logger.Error("failed to parse planned rep", slog.String("error", err.Error()))
			return
		}
		plannedReps[i] = int32(rep)
	}

	plannedWeightsSlice := r.PostForm["planned-weights[]"]
	plannedWeights := make([]int32, len(plannedWeightsSlice))
	for i, plannedWeight := range plannedWeightsSlice {
		weight, err := strconv.ParseInt(plannedWeight, 10, 32)
		if err != nil {
			HandleInternalServerError(w, r)
			h.cfg.Logger.Error("failed to parse planned weight", slog.String("error", err.Error()))
			return
		}
		plannedWeights[i] = int32(weight)
	}

	completedSets, err := strconv.ParseInt(r.FormValue("completed-sets"), 10, 32)
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to parse completed sets", slog.String("error", err.Error()))
		return
	}

	completedRepsSlice := r.PostForm["completed-reps[]"]
	completedReps := make([]int32, len(completedRepsSlice))
	for i, completedRep := range completedRepsSlice {
		rep, err := strconv.ParseInt(completedRep, 10, 32)
		if err != nil {
			HandleInternalServerError(w, r)
			h.cfg.Logger.Error("failed to parse completed rep", slog.String("error", err.Error()))
			return
		}
		completedReps[i] = int32(rep)
	}

	completedWeightsSlice := r.PostForm["completed-weights[]"]
	completedWeights := make([]int32, len(completedWeightsSlice))
	for i, completedWeight := range completedWeightsSlice {
		weight, err := strconv.ParseInt(completedWeight, 10, 32)
		if err != nil {
			HandleInternalServerError(w, r)
			h.cfg.Logger.Error("failed to parse completed weight", slog.String("error", err.Error()))
			return
		}
		completedWeights[i] = int32(weight)
	}

	updatedWorkoutExercise, err := h.cfg.DB.UpdateWorkoutExercise(r.Context(), database.UpdateWorkoutExerciseParams{
		SetsPlanned:         int32(plannedSets),
		RepsPerSetPlanned:   plannedReps,
		SetsCompleted:       int32(completedSets),
		RepsPerSetCompleted: completedReps,
		WeightsPlannedLbs:   plannedWeights,
		WeightsCompletedLbs: completedWeights,
		ID:                  workoutExerciseID,
		WorkoutID:           workoutID,
	})
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to update workout exercise", slog.String("error", err.Error()))
		return
	}

	exerciseForWorkout := database.WorkoutAndExercisesRow{
		WorkoutsExercise: updatedWorkoutExercise,
		Exercise:         exercise,
	}

	err = templates.WorkoutExercisesTableDataRow(exerciseForWorkout).Render(r.Context(), w)
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to render workout exercises table data row", slog.String("error", err.Error()))
		return
	}
}

func (h *Handler) UpdateWorkoutExercisesSortOrder(w http.ResponseWriter, r *http.Request) {
	workoutID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to parse workout id", slog.String("error", err.Error()))
		return
	}

	if err := r.ParseForm(); err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to parse workout exercises form", slog.String("error", err.Error()))
		return
	}
	sortOrder := r.PostForm["sort-order[]"]
	for index, workoutExerciseID := range sortOrder {
		id, err := uuid.Parse(workoutExerciseID)
		if err != nil {
			HandleInternalServerError(w, r)
			h.cfg.Logger.Error("failed to parse workout exercise id", slog.String("error", err.Error()))
			return
		}

		err = h.cfg.DB.UpdateWorkoutExercisesSortOrder(r.Context(), database.UpdateWorkoutExercisesSortOrderParams{
			SortOrder: int32(index),
			ID:        id,
			WorkoutID: workoutID,
		})
		if err != nil {
			HandleInternalServerError(w, r)
			h.cfg.Logger.Error("failed to update workout exercise sort order", slog.String("error", err.Error()))
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DeleteExerciseFromWorkout(w http.ResponseWriter, r *http.Request) {
	workoutExerciseID, err := uuid.Parse(r.PathValue("workoutExerciseID"))
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to parse workout exercise id", slog.String("error", err.Error()))
		return
	}

	err = h.cfg.DB.DeleteExerciseFromWorkout(r.Context(), workoutExerciseID)
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to delete exercise from workout", slog.String("error", err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetExercisesPage(w http.ResponseWriter, r *http.Request) {
	// TODO: Need to change this or the html side -- Either, or separate into
	// primary/secondary
	muscleGroups, err := h.cfg.DB.GetAllMuscleGroups(r.Context())
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to get muscle groups", slog.String("error", err.Error()))
		return
	}

	mgrps := []string{}
	for _, mg := range muscleGroups {
		if mg.Valid {
			mgrps = append(mgrps, mg.String)
		}
	}

	err = templates.Layout(templates.ExercisesPage(mgrps), "FitHub | Exercises", true).Render(r.Context(), w)
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to render exercises page", slog.String("error", err.Error()))
		return
	}
}

func (h *Handler) GetMGExercisesPage(w http.ResponseWriter, r *http.Request) {
	group := r.PathValue("group")
	exercises, err := h.cfg.DB.GetExercisesByPrimaryMG(r.Context(), sql.NullString{String: group, Valid: true})
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to get exercises by muscle group", slog.String("error", err.Error()))
		return
	}

	err = templates.Layout(templates.ExercisesByGroup(group, exercises), "FtiHub | Exercises", true).Render(r.Context(), w)
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to render exercises by muscle group", slog.String("error", err.Error()))
		return
	}
}
