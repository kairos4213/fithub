package handlers

import (
	"database/sql"
	"log/slog"
	"net/http"
	"strings"

	"github.com/kairos4213/fithub/internal/auth"
	"github.com/kairos4213/fithub/internal/database"
	"github.com/kairos4213/fithub/internal/templates"
	"github.com/kairos4213/fithub/internal/validate"
)

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// TODO: Make this page inaccessible unless user is logged out
		w.Header().Set("Content-type", "text/html")
		contents := templates.RegisterPage()
		err := templates.Layout(contents, "FitHub | Register", false).Render(r.Context(), w)
		if err != nil {
			HandleInternalServerError(w, r)
			h.cfg.Logger.Error("failed to render register page", slog.String("error", err.Error()))
			return
		}
		return
	}

	if r.Method == http.MethodPost {
		firstName := r.FormValue("first_name")
		lastName := r.FormValue("last_name")
		email := r.FormValue("email")
		password := r.FormValue("password")

		if errs := validate.Fields(
			validate.Required(firstName, "first name"),
			validate.Required(lastName, "last name"),
			validate.Required(email, "email"),
			validate.Required(password, "password"),
			validate.MinLen(password, 10, "password"),
			validate.MaxLen(firstName, 100, "first name"),
			validate.MaxLen(lastName, 100, "last name"),
			validate.MaxLen(email, 255, "email"),
		); errs != nil {
			HandleFieldErrors(w, r, h.cfg.Logger, errs, []string{"first-name", "last-name", "email", "password"}, "")
			return
		}

		hashedPassword, err := auth.HashPassword(password)
		if err != nil {
			HandleInternalServerError(w, r)
			h.cfg.Logger.Error("failed to hash password", slog.String("error", err.Error()))
			return
		}

		user, err := h.cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
			FirstName:      firstName,
			LastName:       lastName,
			Email:          email,
			HashedPassword: sql.NullString{String: hashedPassword, Valid: true},
		})
		if err != nil {
			if strings.Contains(err.Error(), "users_email_key") {
				HandleRegPageEmailAlert(w, r)
				h.cfg.Logger.Info("duplicate db email", slog.String("error", err.Error()))
				return
			}
			HandleInternalServerError(w, r)
			h.cfg.Logger.Error("failed to create user", slog.String("error", err.Error()))
			return
		}

		_, _, err = h.issueSessionTokens(r.Context(), w, user.ID, user.IsAdmin)
		if err != nil {
			HandleInternalServerError(w, r)
			h.cfg.Logger.Error("failed to issue session tokens", slog.String("error", err.Error()))
			return
		}

		w.Header().Set("Content-type", "text/html")
		w.Header().Set("HX-Location", `{"path": "/workouts"}`)
		w.WriteHeader(http.StatusCreated)
	}
}

func (h *Handler) CheckUserEmail(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")

	user, err := h.cfg.DB.GetUser(r.Context(), email)
	if err == nil {
		HandleRegPageEmailAlert(w, r)
		h.cfg.Logger.Info("email already exists alert", slog.String("email", user.Email))
		return
	}

	err = templates.RegPageEmailAlert(templates.HtmlErr{}).Render(r.Context(), w)
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to render registration page email alert", slog.String("error", err.Error()))
		return
	}
}
