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
		templates.Layout(contents, "FitHub | Login").Render(r.Context(), w)
		return
	}

	if r.Method == "POST" {
		email := r.FormValue("email")
		password := r.FormValue("password")

		if email == "" || password == "" {
			http.Error(w, "email or password cannot be blank", http.StatusBadRequest)
			return
		}

		user, err := h.DB.GetUser(r.Context(), email)
		if err != nil {
			http.Error(w, "User does not exist", http.StatusUnauthorized)
			return
		}

		if err = auth.CheckPasswordHash(password, user.HashedPassword); err != nil {
			http.Error(w, "Incorrect password", http.StatusUnauthorized)
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
			Name:     "token",
			Value:    accessToken,
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteDefaultMode,
		})
		w.Header().Set("Content-type", "text/html")
		w.Header().Set("HX-Location", `{"path": "/workouts"}`)
		w.WriteHeader(http.StatusAccepted)
	}
}
