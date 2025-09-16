package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/kairos4213/fithub/internal/database"
	"github.com/kairos4213/fithub/internal/templates"
)

func (h *Handler) GetAdminExercisesPage(w http.ResponseWriter, r *http.Request) {
	exercises, err := h.DB.GetAllExercises(r.Context())
	if err != nil {
		HandleInternalServerError(w, r)
		log.Printf("Server Error: %v", err)
		return
	}

	contents := templates.AdminExercisesPage(exercises)
	templates.AdminLayout(contents, "FitHub-Admin | Home", true).Render(r.Context(), w)
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

	_, err := h.DB.CreateExercise(r.Context(), database.CreateExerciseParams{
		Name:                 name,
		Description:          description,
		PrimaryMuscleGroup:   primeMG,
		SecondaryMuscleGroup: secMG,
	})
	if err != nil {
		HandleInternalServerError(w, r)
		log.Printf("Server Error: %v", err)
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
		log.Printf("Server Error: %v", err)
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

	updatedExercise, err := h.DB.UpdateExercise(r.Context(), database.UpdateExerciseParams{
		Name:                 reqName,
		Description:          description,
		PrimaryMuscleGroup:   primaryMG,
		SecondaryMuscleGroup: secondaryMG,
		ID:                   exerciseID,
	})
	if err != nil {
		HandleInternalServerError(w, r)
		log.Printf("Server Error: %v", err)
		return
	}

	templates.ExerciseDataRow(updatedExercise).Render(r.Context(), w)
}

func (h *Handler) DeleteDBExercise(w http.ResponseWriter, r *http.Request) {
	exerciseID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		HandleInternalServerError(w, r)
		log.Printf("Server Error: %v", err)
		return
	}

	err = h.DB.DeleteExercise(r.Context(), exerciseID)
	if err != nil {
		HandleInternalServerError(w, r)
		log.Printf("Server Error: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetExerciseByKeyword(w http.ResponseWriter, r *http.Request) {
	exerciseSearch := strings.ToLower(r.FormValue("exercise-search"))
	workoutID, err := uuid.Parse(r.FormValue("workoutID"))
	if err != nil {
		HandleInternalServerError(w, r)
		log.Printf("Server Error: %v", err)
		return
	}

	if exerciseSearch == "" {
		templates.ExerciseSearchResults([]database.Exercise{}, workoutID).Render(r.Context(), w)
		return
	}

	exercises, err := h.DB.GetExerciseByKeyword(r.Context(), exerciseSearch)
	if err != nil {
		HandleInternalServerError(w, r)
		log.Printf("Server Error: %v", err)
		return
	}

	templates.ExerciseSearchResults(exercises, workoutID).Render(r.Context(), w)
}

func (h *Handler) AddExerciseToWorkout(w http.ResponseWriter, r *http.Request) {
	workoutID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		HandleInternalServerError(w, r)
		log.Printf("Server err: %v", err)
		return
	}

	exerciseName := r.FormValue("exercise-name")
	exercise, err := h.DB.GetExerciseByName(r.Context(), exerciseName)
	if err != nil {
		HandleInternalServerError(w, r)
		log.Printf("Server err: %v", err)
		return
	}

	plannedSets, err := strconv.ParseInt(r.FormValue("planned-sets"), 10, 32)
	if err != nil {
		HandleInternalServerError(w, r)
		log.Printf("Server err: %v", err)
		return
	}

	plannedRepsSlice := r.PostForm["planned-reps[]"]
	plannedReps := make([]int32, len(plannedRepsSlice))
	for i, plannedRep := range plannedRepsSlice {
		rep, err := strconv.ParseInt(plannedRep, 10, 32)
		if err != nil {
			HandleInternalServerError(w, r)
			log.Printf("Server err: %v", err)
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
			log.Printf("Server err: %v", err)
			return
		}
		plannedWeights[i] = int32(weight)
	}

	workoutExercise, err := h.DB.AddExerciseToWorkout(r.Context(), database.AddExerciseToWorkoutParams{
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
		log.Printf("Server err: %v", err)
		return
	}

	exerciseForWorkout := database.WorkoutAndExercisesRow{
		WorkoutsExercise: workoutExercise,
		Exercise:         exercise,
	}

	templates.WorkoutExercisesTableDataRow(exerciseForWorkout).Render(r.Context(), w)
}

func (h *Handler) UpdateWorkoutExercise(w http.ResponseWriter, r *http.Request) {
	workoutID, err := uuid.Parse(r.PathValue("workoutID"))
	if err != nil {
		HandleInternalServerError(w, r)
		log.Printf("Server Error: %v", err)
		return
	}
	workoutExerciseID, err := uuid.Parse(r.PathValue("workoutExerciseID"))
	if err != nil {
		HandleInternalServerError(w, r)
		log.Printf("Server Error: %v", err)
		return
	}

	exerciseID, err := uuid.Parse(r.FormValue("exercise"))
	if err != nil {
		HandleInternalServerError(w, r)
		log.Printf("Server Error with parsing exercise id: %v", err)
		return
	}

	exercise, err := h.DB.GetExerciseByID(r.Context(), exerciseID)
	if err != nil {
		HandleInternalServerError(w, r)
		log.Printf("Server err with fetching exercise by id: %v", err)
		return
	}

	plannedSets, err := strconv.ParseInt(r.FormValue("planned-sets"), 10, 32)
	if err != nil {
		HandleInternalServerError(w, r)
		log.Printf("Server err: %v", err)
		return
	}

	plannedRepsSlice := r.PostForm["planned-reps[]"]
	plannedReps := make([]int32, len(plannedRepsSlice))
	for i, plannedRep := range plannedRepsSlice {
		rep, err := strconv.ParseInt(plannedRep, 10, 32)
		if err != nil {
			HandleInternalServerError(w, r)
			log.Printf("Server err: %v", err)
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
			log.Printf("Server err: %v", err)
			return
		}
		plannedWeights[i] = int32(weight)
	}

	completedSets, err := strconv.ParseInt(r.FormValue("completed-sets"), 10, 32)
	if err != nil {
		HandleInternalServerError(w, r)
		log.Printf("Server err: %v", err)
		return
	}

	completedRepsSlice := r.PostForm["completed-reps[]"]
	completedReps := make([]int32, len(completedRepsSlice))
	for i, completedRep := range completedRepsSlice {
		rep, err := strconv.ParseInt(completedRep, 10, 32)
		if err != nil {
			HandleInternalServerError(w, r)
			log.Printf("Server err: %v", err)
			return
		}
		completedReps[i] = int32(rep)
	}

	completedWeightsSlice := r.PostForm["completed-weights[]"]
	completedWeights := make([]int32, len(completedWeightsSlice))
	for i, completedWeight := range completedRepsSlice {
		weight, err := strconv.ParseInt(completedWeight, 10, 32)
		if err != nil {
			HandleInternalServerError(w, r)
			log.Printf("Server err: %v", err)
			return
		}
		completedWeights[i] = int32(weight)
	}

	updatedWorkoutExercise, err := h.DB.UpdateWorkoutExercise(r.Context(), database.UpdateWorkoutExerciseParams{
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
		log.Printf("Server err: %v", err)
		return
	}

	exerciseForWorkout := database.WorkoutAndExercisesRow{
		WorkoutsExercise: updatedWorkoutExercise,
		Exercise:         exercise,
	}

	templates.WorkoutExercisesTableDataRow(exerciseForWorkout).Render(r.Context(), w)
}

func (h *Handler) UpdateWorkoutExercisesSortOrder(w http.ResponseWriter, r *http.Request) {
	workoutID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		HandleInternalServerError(w, r)
		log.Printf("Server Error: %v", err)
		return
	}

	id, err := uuid.Parse(r.FormValue("id"))
	if err != nil {
		HandleInternalServerError(w, r)
		log.Printf("Server Error: %v", err)
		return
	}

	sortOrder, err := strconv.ParseInt(r.FormValue("sort-order"), 10, 32)
	if err != nil {
		HandleInternalServerError(w, r)
		log.Printf("Server Error: %v", err)
		return
	}

	err = h.DB.UpdateWorkoutExercisesSortOrder(r.Context(), database.UpdateWorkoutExercisesSortOrderParams{
		SortOrder: int32(sortOrder),
		ID:        id,
		WorkoutID: workoutID,
	})
	if err != nil {
		HandleInternalServerError(w, r)
		log.Printf("Server Error: %v", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DeleteExerciseFromWorkout(w http.ResponseWriter, r *http.Request) {
	workoutExerciseID, err := uuid.Parse(r.PathValue("workoutExerciseID"))
	if err != nil {
		HandleInternalServerError(w, r)
		log.Printf("Server Error: %v", err)
		return
	}

	err = h.DB.DeleteExerciseFromWorkout(r.Context(), workoutExerciseID)
	if err != nil {
		HandleInternalServerError(w, r)
		log.Printf("Server Error: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
