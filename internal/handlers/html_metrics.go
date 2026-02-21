package handlers

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/kairos4213/fithub/internal/cntx"
	"github.com/kairos4213/fithub/internal/database"
	"github.com/kairos4213/fithub/internal/templates"
	"github.com/kairos4213/fithub/internal/validate"
)

func (h *Handler) GetAllMetrics(w http.ResponseWriter, r *http.Request) {
	userID, ok := cntx.UserID(r.Context())
	if !ok {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("missing user id in context")
		return
	}

	tab := r.URL.Query().Get("tab")
	if tab != "muscleMasses" && tab != "bfPercents" {
		tab = "bodyweights"
	}

	var bodyweights []database.BodyWeight
	var muscleMasses []database.MuscleMass
	var bfPercents []database.BodyFatPercent
	var err error

	switch tab {
	case "bodyweights":
		bodyweights, err = h.cfg.DB.GetAllBodyWeights(r.Context(), userID)
	case "muscleMasses":
		muscleMasses, err = h.cfg.DB.GetAllMuscleMasses(r.Context(), userID)
	case "bfPercents":
		bfPercents, err = h.cfg.DB.GetAllBodyFatPercs(r.Context(), userID)
	}
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to get metrics", slog.String("error", err.Error()))
		return
	}

	// HTMX tab switch — return just the content fragment
	target := r.Header.Get("HX-Target")
	if target == "metrics-content" {
		var renderErr error
		switch tab {
		case "bodyweights":
			renderErr = templates.BodyweightsContent(bodyweights).Render(r.Context(), w)
		case "muscleMasses":
			renderErr = templates.MuscleMassesContent(muscleMasses).Render(r.Context(), w)
		case "bfPercents":
			renderErr = templates.BfPercentsContent(bfPercents).Render(r.Context(), w)
		}
		if renderErr != nil {
			HandleInternalServerError(w, r)
			h.cfg.Logger.Error("failed to render metrics content", slog.String("error", renderErr.Error()))
		}
		return
	}

	// Full page render
	contents := templates.MetricsPage(tab, bodyweights, muscleMasses, bfPercents)
	err = templates.Layout(contents, "Fithub | Metrics", true).Render(r.Context(), w)
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to render metrics page", slog.String("error", err.Error()))
		return
	}
}

func (h *Handler) LogMetrics(w http.ResponseWriter, r *http.Request) {
	userID, ok := cntx.UserID(r.Context())
	if !ok {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("missing user id in context")
		return
	}

	metricType := r.PathValue("type")
	switch metricType {
	case "bodyweights":
		entry := r.FormValue("bodyweight")
		if errs := validate.Fields(
			validate.Required(entry, "bodyweight"),
			validate.Numeric(entry, "bodyweight"),
		); errs != nil {
			HandleFieldErrors(w, r, h.cfg.Logger, errs, []string{"bodyweight"}, "")
			return
		}

		bw, err := h.cfg.DB.AddBodyWeight(r.Context(), database.AddBodyWeightParams{UserID: userID, Measurement: entry})
		if err != nil {
			HandleInternalServerError(w, r)
			h.cfg.Logger.Error("failed to add body weight", slog.String("error", err.Error()))
			return
		}

		w.Header().Set("HX-Trigger", "close-log-bw-card")
		err = templates.BWRow(bw).Render(r.Context(), w)
		if err != nil {
			HandleInternalServerError(w, r)
			h.cfg.Logger.Error("failed to render body weights", slog.String("error", err.Error()))
			return
		}
		err = templates.MetricsEmptyOOB(false, "").Render(r.Context(), w)
		if err != nil {
			HandleInternalServerError(w, r)
			h.cfg.Logger.Error("failed to render metrics empty oob", slog.String("error", err.Error()))
			return
		}
	case "muscleMasses":
		entry := r.FormValue("muscle-mass")
		if errs := validate.Fields(
			validate.Required(entry, "muscle mass"),
			validate.Numeric(entry, "muscle mass"),
		); errs != nil {
			HandleFieldErrors(w, r, h.cfg.Logger, errs, []string{"muscle-mass"}, "")
			return
		}

		mm, err := h.cfg.DB.AddMuscleMass(r.Context(), database.AddMuscleMassParams{UserID: userID, Measurement: entry})
		if err != nil {
			HandleInternalServerError(w, r)
			h.cfg.Logger.Error("failed to add muscle mass", slog.String("error", err.Error()))
			return
		}

		w.Header().Set("HX-Trigger", "close-log-mm-card")
		err = templates.MMRow(mm).Render(r.Context(), w)
		if err != nil {
			HandleInternalServerError(w, r)
			h.cfg.Logger.Error("failed to render muscle masses", slog.String("error", err.Error()))
			return
		}
		err = templates.MetricsEmptyOOB(false, "").Render(r.Context(), w)
		if err != nil {
			HandleInternalServerError(w, r)
			h.cfg.Logger.Error("failed to render metrics empty oob", slog.String("error", err.Error()))
			return
		}
	case "bfPercents":
		entry := r.FormValue("bf-percent")
		if errs := validate.Fields(
			validate.Required(entry, "body fat percent"),
			validate.Numeric(entry, "body fat percent"),
		); errs != nil {
			HandleFieldErrors(w, r, h.cfg.Logger, errs, []string{"body-fat-percent"}, "")
			return
		}

		bf, err := h.cfg.DB.AddBodyFatPerc(r.Context(), database.AddBodyFatPercParams{UserID: userID, Measurement: entry})
		if err != nil {
			HandleInternalServerError(w, r)
			h.cfg.Logger.Error("failed to add body fat", slog.String("error", err.Error()))
			return
		}

		w.Header().Set("HX-Trigger", "close-log-bf-card")
		err = templates.BFRow(bf).Render(r.Context(), w)
		if err != nil {
			HandleInternalServerError(w, r)
			h.cfg.Logger.Error("failed to render body fats", slog.String("error", err.Error()))
			return
		}
		err = templates.MetricsEmptyOOB(false, "").Render(r.Context(), w)
		if err != nil {
			HandleInternalServerError(w, r)
			h.cfg.Logger.Error("failed to render metrics empty oob", slog.String("error", err.Error()))
			return
		}
	}
}

func (h *Handler) EditMetrics(w http.ResponseWriter, r *http.Request) {
	userID, ok := cntx.UserID(r.Context())
	if !ok {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("missing user id in context")
		return
	}

	metricType := r.PathValue("type")
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to parse metric id", slog.String("error", err.Error()))
		return
	}
	prefix := fmt.Sprintf("%v-", id)
	switch metricType {
	case "bodyweights":
		entry := r.FormValue("bodyweight")
		if errs := validate.Fields(
			validate.Required(entry, "bodyweight"),
			validate.Numeric(entry, "bodyweight"),
		); errs != nil {
			HandleScopedFieldErrors(w, r, h.cfg.Logger, errs, []string{prefix + "bodyweight"}, prefix, fmt.Sprintf("form-error-bw-%v", id))
			return
		}

		updatedBW, err := h.cfg.DB.UpdateBodyWeight(r.Context(), database.UpdateBodyWeightParams{Measurement: entry, ID: id, UserID: userID})
		if err != nil {
			HandleInternalServerError(w, r)
			h.cfg.Logger.Error("failed to update body weight", slog.String("error", err.Error()))
			return
		}

		err = templates.BWRow(updatedBW).Render(r.Context(), w)
		if err != nil {
			HandleInternalServerError(w, r)
			h.cfg.Logger.Error("failed to render body weight", slog.String("error", err.Error()))
			return
		}
	case "muscleMasses":
		entry := r.FormValue("muscle-mass")
		if errs := validate.Fields(
			validate.Required(entry, "muscle mass"),
			validate.Numeric(entry, "muscle mass"),
		); errs != nil {
			HandleScopedFieldErrors(w, r, h.cfg.Logger, errs, []string{prefix + "muscle-mass"}, prefix, fmt.Sprintf("form-error-mm-%v", id))
			return
		}

		updatedMM, err := h.cfg.DB.UpdateMuscleMass(r.Context(), database.UpdateMuscleMassParams{Measurement: entry, ID: id, UserID: userID})
		if err != nil {
			HandleInternalServerError(w, r)
			h.cfg.Logger.Error("failed to update muscle mass", slog.String("error", err.Error()))
			return
		}

		err = templates.MMRow(updatedMM).Render(r.Context(), w)
		if err != nil {
			HandleInternalServerError(w, r)
			h.cfg.Logger.Error("failed to render muscle mass", slog.String("error", err.Error()))
			return
		}
	case "bfPercents":
		entry := r.FormValue("bf-percent")
		if errs := validate.Fields(
			validate.Required(entry, "body fat percent"),
			validate.Numeric(entry, "body fat percent"),
		); errs != nil {
			HandleScopedFieldErrors(w, r, h.cfg.Logger, errs, []string{prefix + "body-fat-percent"}, prefix, fmt.Sprintf("form-error-bf-%v", id))
			return
		}

		updatedBF, err := h.cfg.DB.UpdateBodyFatPerc(r.Context(), database.UpdateBodyFatPercParams{Measurement: entry, ID: id, UserID: userID})
		if err != nil {
			HandleInternalServerError(w, r)
			h.cfg.Logger.Error("failed to update body fat", slog.String("error", err.Error()))
			return
		}

		err = templates.BFRow(updatedBF).Render(r.Context(), w)
		if err != nil {
			HandleInternalServerError(w, r)
			h.cfg.Logger.Error("failed to render body fat", slog.String("error", err.Error()))
			return
		}
	}
}

func (h *Handler) DeleteMetrics(w http.ResponseWriter, r *http.Request) {
	userID, ok := cntx.UserID(r.Context())
	if !ok {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("missing user id in context")
		return
	}

	metricType := r.PathValue("type")
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to parse metric id", slog.String("error", err.Error()))
		return
	}
	var count int64
	var emptyMessage string

	switch metricType {
	case "bodyweights":
		count, err = h.cfg.DB.DeleteBodyWeight(r.Context(), database.DeleteBodyWeightParams{ID: id, UserID: userID})
		if err != nil {
			HandleInternalServerError(w, r)
			h.cfg.Logger.Error("failed to delete body weight", slog.String("error", err.Error()))
			return
		}
		emptyMessage = "No body weight entries yet."
	case "muscleMasses":
		count, err = h.cfg.DB.DeleteMuscleMass(r.Context(), database.DeleteMuscleMassParams{ID: id, UserID: userID})
		if err != nil {
			HandleInternalServerError(w, r)
			h.cfg.Logger.Error("failed to delete muscle mass", slog.String("error", err.Error()))
			return
		}
		emptyMessage = "No muscle mass entries yet."
	case "bfPercents":
		count, err = h.cfg.DB.DeleteBodyFatPerc(r.Context(), database.DeleteBodyFatPercParams{ID: id, UserID: userID})
		if err != nil {
			HandleInternalServerError(w, r)
			h.cfg.Logger.Error("failed to delete body fat", slog.String("error", err.Error()))
			return
		}
		emptyMessage = "No body fat entries yet."
	}

	if count <= 1 {
		err = templates.MetricsEmptyOOB(true, emptyMessage).Render(r.Context(), w)
		if err != nil {
			HandleInternalServerError(w, r)
			h.cfg.Logger.Error("failed to render metrics empty oob", slog.String("error", err.Error()))
			return
		}
		return
	}
	w.WriteHeader(http.StatusOK)
}
