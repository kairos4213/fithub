package handlers

import (
	"log/slog"
	"net/http"

	"github.com/kairos4213/fithub/internal/auth"
	"github.com/kairos4213/fithub/internal/templates"
	"github.com/kairos4213/fithub/internal/validate"
)

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// TODO: Make this page inaccessible unless user is logged out
		contents := templates.LoginPage()
		err := templates.Layout(contents, "FitHub | Login", false).Render(r.Context(), w)
		if err != nil {
			HandleInternalServerError(w, r)
			h.cfg.Logger.Error("failed to render login page", slog.String("error", err.Error()))
			return
		}
		return
	}

	if r.Method == "POST" {
		email := r.FormValue("email")
		password := r.FormValue("password")

		if errs := validate.Fields(
			validate.Required(email, "email"),
			validate.Required(password, "password"),
		); errs != nil {
			HandleFieldErrors(w, r, h.cfg.Logger, errs, []string{"email", "password"}, "")
			return
		}

		user, err := h.cfg.DB.GetUser(r.Context(), email)
		if err != nil {
			HandleLoginFailure(w, r)
			h.cfg.Logger.Error("failed to fetch user", slog.String("error", err.Error()))
			return
		}

		if !user.HashedPassword.Valid {
			HandleLoginFailure(w, r)
			h.cfg.Logger.Info("password login attempted for OAuth-only user", slog.String("email", email))
			return
		}

		match, err := auth.CheckPasswordHash(password, user.HashedPassword.String)
		if err != nil {
			HandleLoginFailure(w, r)
			h.cfg.Logger.Error("bad request: invalid hash", slog.String("error", err.Error()))
			return
		}

		if !match {
			HandleLoginFailure(w, r)
			h.cfg.Logger.Info("incorrect password attempt", slog.String("user_email", user.Email), slog.String("ip", r.RemoteAddr))
			return
		}

		_, _, err = h.issueSessionTokens(r.Context(), w, user.ID, user.IsAdmin)
		if err != nil {
			HandleInternalServerError(w, r)
			h.cfg.Logger.Error("failed to issue session tokens", slog.String("error", err.Error()))
			return
		}

		w.Header().Set("Content-type", "text/html")

		if user.IsAdmin {
			w.Header().Set("HX-Location", `{"path": "/admin"}`)
			w.WriteHeader(http.StatusAccepted)
		}

		w.Header().Set("HX-Location", `{"path": "/workouts"}`)
		w.WriteHeader(http.StatusAccepted)
	}
}
