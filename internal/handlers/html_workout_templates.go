package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kairos4213/fithub/internal/cntx"
	"github.com/kairos4213/fithub/internal/database"
	"github.com/kairos4213/fithub/internal/templates"
	"github.com/kairos4213/fithub/internal/utils"
	"github.com/kairos4213/fithub/internal/validate"
)

// exerciseSlot represents one exercise entry in the template JSONB.
type exerciseSlot struct {
	Sets       int `json:"sets"`
	RepsPerSet int `json:"reps_per_set"`
}

func (h *Handler) GetAllWorkoutTemplates(w http.ResponseWriter, r *http.Request) {
	workoutTemplates, err := h.cfg.DB.GetAllWorkoutTemplates(r.Context())
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to get workout templates", slog.String("error", err.Error()))
		return
	}

	cards := make([]templates.TemplateCardData, 0, len(workoutTemplates))
	for _, t := range workoutTemplates {
		var exerciseData map[string][]exerciseSlot
		if err := json.Unmarshal(t.ExerciseSetReps, &exerciseData); err != nil {
			h.cfg.Logger.Error("failed to parse template JSONB", slog.String("template", t.TemplateName), slog.String("error", err.Error()))
			continue
		}

		muscleGroups := make([]string, 0, len(exerciseData))
		exerciseCount := 0
		for group, slots := range exerciseData {
			muscleGroups = append(muscleGroups, group)
			exerciseCount += len(slots)
		}

		cards = append(cards, templates.TemplateCardData{
			ID:              t.ID,
			Name:            utils.TitleString(t.TemplateName),
			Description:     t.Description,
			DurationMinutes: t.DurationMinutes,
			MuscleGroups:    muscleGroups,
			ExerciseCount:   exerciseCount,
		})
	}

	err = templates.Layout(templates.WorkoutTemplatesPage(cards), "FitHub | Templates", true).Render(r.Context(), w)
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to render workout templates", slog.String("error", err.Error()))
		return
	}
}

func (h *Handler) GetTemplatePreview(w http.ResponseWriter, r *http.Request) {
	_, ok := cntx.UserID(r.Context())
	if !ok {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("missing user id in context")
		return
	}

	templateID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		HandleBadRequest(w, r, "invalid template id")
		h.cfg.Logger.Info("failed to parse template id", slog.String("error", err.Error()))
		return
	}

	tmpl, err := h.cfg.DB.GetWorkoutTemplateByID(r.Context(), templateID)
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to get workout template", slog.String("error", err.Error()))
		return
	}

	var exerciseData map[string][]exerciseSlot
	if err := json.Unmarshal(tmpl.ExerciseSetReps, &exerciseData); err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to parse template JSONB", slog.String("error", err.Error()))
		return
	}

	exercises, err := h.resolveTemplateExercises(r, exerciseData)
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to resolve template exercises", slog.String("error", err.Error()))
		return
	}

	previewData := templates.TemplatePreviewData{
		TemplateID:  tmpl.ID,
		Title:       utils.TitleString(tmpl.TemplateName),
		Description: tmpl.Description,
		Duration:    tmpl.DurationMinutes,
		Exercises:   exercises,
	}

	err = templates.Layout(templates.TemplatePreviewPage(previewData), "FitHub | Preview Template", true).Render(r.Context(), w)
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to render template preview", slog.String("error", err.Error()))
		return
	}
}

// resolveTemplateExercises picks random exercises for each muscle group slot.
func (h *Handler) resolveTemplateExercises(r *http.Request, exerciseData map[string][]exerciseSlot) ([]templates.PreviewExercise, error) {
	var exercises []templates.PreviewExercise

	for group, slots := range exerciseData {
		count := int32(len(slots))
		dbExercises, err := h.cfg.DB.GetRandomExercisesByMuscleGroup(r.Context(), database.GetRandomExercisesByMuscleGroupParams{
			PrimaryMuscleGroup: sql.NullString{String: group, Valid: true},
			Limit:              count,
		})
		if err != nil {
			return nil, fmt.Errorf("fetching exercises for %s: %w", group, err)
		}

		for i, slot := range slots {
			if i >= len(dbExercises) {
				break
			}
			exercises = append(exercises, templates.PreviewExercise{
				ExerciseID:  dbExercises[i].ID,
				Name:        utils.TitleString(dbExercises[i].Name),
				MuscleGroup: group,
				Sets:        int32(slot.Sets),
				RepsPerSet:  int32(slot.RepsPerSet),
				Weight:      0,
			})
		}
	}

	return exercises, nil
}

func (h *Handler) ApplyTemplate(w http.ResponseWriter, r *http.Request) {
	userID, ok := cntx.UserID(r.Context())
	if !ok {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("missing user id in context")
		return
	}

	reqTitle := r.FormValue("title")
	reqDescription := r.FormValue("description")
	reqDuration := r.FormValue("duration")
	reqPlannedDate := r.FormValue("planned-date")
	reqExerciseCount := r.FormValue("exercise_count")

	templateFields := []string{"title", "duration", "planned-date", "description"}
	if errs := validate.Fields(
		validate.Required(reqTitle, "title"),
		validate.Required(reqDuration, "duration"),
		validate.Required(reqPlannedDate, "planned date"),
		validate.Required(reqExerciseCount, "exercise count"),
		validate.MaxLen(reqTitle, 100, "title"),
		validate.MaxLen(reqDescription, 500, "description"),
		validate.Numeric(reqDuration, "duration"),
		validate.Numeric(reqExerciseCount, "exercise count"),
	); errs != nil {
		HandleFieldErrors(w, r, h.cfg.Logger, errs, templateFields, "")
		return
	}

	description := sql.NullString{Valid: false}
	if reqDescription != "" {
		description.String = reqDescription
		description.Valid = true
	}

	duration, err := strconv.ParseInt(reqDuration, 10, 32)
	if err != nil {
		HandleBadRequest(w, r, "duration must be a number")
		h.cfg.Logger.Info("invalid duration", slog.String("value", reqDuration), slog.String("error", err.Error()))
		return
	}

	plannedDate, err := time.Parse(time.DateOnly, reqPlannedDate)
	if err != nil {
		HandleBadRequest(w, r, "planned date must be in YYYY-MM-DD format")
		h.cfg.Logger.Info("invalid planned date", slog.String("value", reqPlannedDate), slog.String("error", err.Error()))
		return
	}

	exerciseCount, err := strconv.Atoi(reqExerciseCount)
	if err != nil {
		HandleBadRequest(w, r, "invalid exercise count")
		h.cfg.Logger.Info("invalid exercise count", slog.String("value", reqExerciseCount), slog.String("error", err.Error()))
		return
	}

	// Parse each exercise from indexed form fields
	type exerciseInput struct {
		exerciseID uuid.UUID
		sets       int32
		reps       []int32
		weights    []int32
	}
	exerciseInputs := make([]exerciseInput, 0, exerciseCount)

	for i := 0; i < exerciseCount; i++ {
		prefix := strconv.Itoa(i)

		exerciseID, err := uuid.Parse(r.FormValue("exercise_id_" + prefix))
		if err != nil {
			HandleBadRequest(w, r, "invalid exercise id")
			h.cfg.Logger.Info("invalid exercise id at index", slog.String("index", prefix), slog.String("error", err.Error()))
			return
		}

		setsStr := r.FormValue("sets_" + prefix)
		sets, err := strconv.ParseInt(setsStr, 10, 32)
		if err != nil || sets < 1 {
			HandleBadRequest(w, r, "sets must be a positive number")
			h.cfg.Logger.Info("invalid sets at index", slog.String("index", prefix), slog.String("value", setsStr))
			return
		}

		// Parse per-set reps array (reps_0[], reps_1[], etc.)
		repsSlice := r.PostForm[fmt.Sprintf("reps_%s[]", prefix)]
		reps := make([]int32, 0, len(repsSlice))
		for _, v := range repsSlice {
			rep, err := strconv.ParseInt(v, 10, 32)
			if err != nil || rep < 1 {
				HandleBadRequest(w, r, "reps must be positive numbers")
				h.cfg.Logger.Info("invalid rep value at index", slog.String("index", prefix), slog.String("value", v))
				return
			}
			reps = append(reps, int32(rep))
		}

		// Parse per-set weights array (weight_0[], weight_1[], etc.)
		weightsSlice := r.PostForm[fmt.Sprintf("weight_%s[]", prefix)]
		weights := make([]int32, 0, len(weightsSlice))
		for _, v := range weightsSlice {
			wt, err := strconv.ParseInt(v, 10, 32)
			if err != nil {
				HandleBadRequest(w, r, "weights must be numbers")
				h.cfg.Logger.Info("invalid weight value at index", slog.String("index", prefix), slog.String("value", v))
				return
			}
			weights = append(weights, int32(wt))
		}

		exerciseInputs = append(exerciseInputs, exerciseInput{
			exerciseID: exerciseID,
			sets:       int32(sets),
			reps:       reps,
			weights:    weights,
		})
	}

	// Create workout + exercises in a transaction
	ctx := r.Context()
	tx, err := h.cfg.RawDB.BeginTx(ctx, nil)
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to begin transaction", slog.String("error", err.Error()))
		return
	}
	defer tx.Rollback()

	qtx := h.cfg.DB.WithTx(tx)

	workout, err := qtx.CreateWorkout(ctx, database.CreateWorkoutParams{
		UserID:          userID,
		Title:           reqTitle,
		Description:     description,
		DurationMinutes: int32(duration),
		PlannedDate:     plannedDate,
	})
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to create workout from template", slog.String("error", err.Error()))
		return
	}

	for _, ex := range exerciseInputs {
		_, err := qtx.AddExerciseToWorkout(ctx, database.AddExerciseToWorkoutParams{
			WorkoutID:           workout.ID,
			ExerciseID:          ex.exerciseID,
			SetsPlanned:         ex.sets,
			RepsPerSetPlanned:   ex.reps,
			SetsCompleted:       0,
			RepsPerSetCompleted: []int32{},
			WeightsPlannedLbs:   ex.weights,
			WeightsCompletedLbs: []int32{},
		})
		if err != nil {
			HandleInternalServerError(w, r)
			h.cfg.Logger.Error("failed to add exercise to workout", slog.String("error", err.Error()))
			return
		}
	}

	if err := tx.Commit(); err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to commit transaction", slog.String("error", err.Error()))
		return
	}

	w.Header().Set("HX-Redirect", "/workouts/"+workout.ID.String())
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) RerollExercise(w http.ResponseWriter, r *http.Request) {
	_, ok := cntx.UserID(r.Context())
	if !ok {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("missing user id in context")
		return
	}

	muscleGroup := r.URL.Query().Get("muscle_group")
	indexStr := r.URL.Query().Get("index")
	excludeStr := r.URL.Query().Get("exclude")
	setsStr := r.URL.Query().Get("sets")
	repsStr := r.URL.Query().Get("reps")

	if muscleGroup == "" || indexStr == "" {
		HandleBadRequest(w, r, "muscle_group and index are required")
		h.cfg.Logger.Info("missing reroll params", slog.String("muscle_group", muscleGroup), slog.String("index", indexStr))
		return
	}

	index, err := strconv.Atoi(indexStr)
	if err != nil {
		HandleBadRequest(w, r, "invalid index")
		h.cfg.Logger.Info("invalid reroll index", slog.String("value", indexStr), slog.String("error", err.Error()))
		return
	}

	sets, _ := strconv.ParseInt(setsStr, 10, 32)
	if sets < 1 {
		sets = 3
	}
	reps, _ := strconv.ParseInt(repsStr, 10, 32)
	if reps < 1 {
		reps = 10
	}

	var excludeIDs []uuid.UUID
	if excludeStr != "" {
		for _, idStr := range strings.Split(excludeStr, ",") {
			if id, err := uuid.Parse(strings.TrimSpace(idStr)); err == nil {
				excludeIDs = append(excludeIDs, id)
			}
		}
	}

	exercise, err := h.cfg.DB.GetRandomExerciseExcluding(r.Context(), database.GetRandomExerciseExcludingParams{
		PrimaryMuscleGroup: sql.NullString{String: muscleGroup, Valid: true},
		ExcludeIds:         excludeIDs,
	})
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to get random exercise", slog.String("muscle_group", muscleGroup), slog.String("error", err.Error()))
		return
	}

	previewExercise := templates.PreviewExercise{
		ExerciseID:  exercise.ID,
		Name:        utils.TitleString(exercise.Name),
		MuscleGroup: muscleGroup,
		Sets:        int32(sets),
		RepsPerSet:  int32(reps),
		Weight:      0,
	}

	err = templates.TemplateExerciseRow(previewExercise, index).Render(r.Context(), w)
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to render exercise row", slog.String("error", err.Error()))
		return
	}
}
