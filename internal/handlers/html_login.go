package handlers

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/kairos4213/fithub/internal/auth"
	"github.com/kairos4213/fithub/internal/database"
	"github.com/kairos4213/fithub/internal/templates"
)

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
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

		user, err := h.cfg.DB.GetUser(r.Context(), email)
		if err != nil {
			HandleInternalServerError(w, r)
			h.cfg.Logger.Error("failed to fetch user", slog.String("error", err.Error()))
			return
		}

		match, err := auth.CheckPasswordHash(password, user.HashedPassword)
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
			MaxAge:   60 * 15, // 15 minutes
		})
		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    refreshToken,
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteDefaultMode,
			MaxAge:   60 * 60 * 24 * 60, // 60 days
		})
		w.Header().Set("Content-type", "text/html")

		if user.IsAdmin {
			w.Header().Set("HX-Location", `{"path": "/admin"}`)
			w.WriteHeader(http.StatusAccepted)
		}

		w.Header().Set("HX-Location", `{"path": "/workouts"}`)
		w.WriteHeader(http.StatusAccepted)
	}
}
