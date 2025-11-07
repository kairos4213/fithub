package handlers

import (
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/kairos4213/fithub/internal/auth"
	"github.com/kairos4213/fithub/internal/database"
	"github.com/kairos4213/fithub/internal/templates"
)

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
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
			HashedPassword: hashedPassword,
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

		accessToken, err := auth.MakeJWT(user.ID, user.IsAdmin, h.cfg.TokenSecret)
		if err != nil {
			HandleInternalServerError(w, r)
			h.cfg.Logger.Error("failed to make JWT", slog.String("error", err.Error()))
			return
		}

		refreshToken, err := auth.MakeRefreshToken()
		if err != nil {
			HandleInternalServerError(w, r)
			h.cfg.Logger.Error("failed to make refresh token", slog.String("error", err.Error()))
			return
		}

		err = h.cfg.DB.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
			Token:     refreshToken,
			UserID:    user.ID,
			ExpiresAt: time.Now().UTC().AddDate(0, 0, 60),
		})
		if err != nil {
			HandleInternalServerError(w, r)
			h.cfg.Logger.Error("failed to store refresh token", slog.String("error", err.Error()))
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    accessToken,
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteDefaultMode,
		})
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
