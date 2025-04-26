package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/kairos4213/fithub/internal/auth"
	"github.com/kairos4213/fithub/internal/database"
	"github.com/kairos4213/fithub/internal/templates"
)

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		contents := templates.RegisterPage()
		templates.Layout(contents, "FitHub | Register").Render(r.Context(), w)
		return
	}

	if r.Method == http.MethodPost {
		firstName := r.FormValue("first_name")
		lastName := r.FormValue("last_name")
		email := r.FormValue("email")
		password := r.FormValue("password")

		log.Default().Printf("%v, %v, %v, %v", firstName, lastName, email, password)

		hashedPassword, err := auth.HashPassword(password)
		if err != nil {
			http.Error(w, "Error hashing password", http.StatusInternalServerError)
			return
		}

		user, err := h.DB.CreateUser(r.Context(), database.CreateUserParams{
			FirstName:      firstName,
			LastName:       lastName,
			Email:          email,
			HashedPassword: hashedPassword,
		})
		if err != nil {
			http.Error(w, fmt.Sprintf("Error creating user: %v", err), http.StatusInternalServerError)
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
		w.Header().Set("HX-Redirect", "/workouts")
		w.WriteHeader(http.StatusCreated)
	}
}
