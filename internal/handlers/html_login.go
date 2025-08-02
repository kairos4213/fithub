package handlers

import (
	"net/http"
	"time"

	"github.com/kairos4213/fithub/internal/auth"
	"github.com/kairos4213/fithub/internal/database"
	"github.com/kairos4213/fithub/internal/templates"
)

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		contents := templates.LoginPage()
		templates.Layout(contents, "FitHub | Login", false).Render(r.Context(), w)
		return
	}

	if r.Method == "POST" {
		email := r.FormValue("email")
		password := r.FormValue("password")

		user, err := h.DB.GetUser(r.Context(), email)
		if err != nil {
			templates.LoginPage().Render(r.Context(), w)
			return
		}

		if err = auth.CheckPasswordHash(password, user.HashedPassword); err != nil {
			templates.LoginPage().Render(r.Context(), w)
			return
		}

		accessToken, err := auth.MakeJWT(user.ID, h.PrivateKey)
		if err != nil {
			http.Error(w, "Issue creating access token", http.StatusInternalServerError)
			return
		}

		refreshToken, err := auth.MakeRefreshToken()
		if err != nil {
			http.Error(w, "Issue creating refresh token", http.StatusInternalServerError)
			return
		}

		err = h.DB.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
			Token:     refreshToken,
			UserID:    user.ID,
			ExpiresAt: time.Now().UTC().AddDate(0, 0, 60),
		})
		if err != nil {
			http.Error(w, "Issue storing refresh token", http.StatusInternalServerError)
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

		if user.IsAdmin {
			w.Header().Set("HX-Location", `{"path": "/admin"}`)
			w.WriteHeader(http.StatusAccepted)
		}

		w.Header().Set("HX-Location", `{"path": "/workouts"}`)
		w.WriteHeader(http.StatusAccepted)
	}
}
