package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/kairos4213/fithub/internal/auth"
	"github.com/kairos4213/fithub/internal/database"
	"github.com/kairos4213/fithub/internal/templates"
)

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	// TODO: start utilizing refresh tokens
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
			HandleLoginFailure(w, r)
			log.Print("User email does not exist")
			return
		}

		match, err := auth.CheckPasswordHash(password, user.HashedPassword)
		if !match {
			HandleLoginFailure(w, r)
			log.Print("Incorrect password entered")
			return
		}
		if err != nil {
			HandleLoginFailure(w, r)
			log.Print("Bad Request: invalid hash")
			return
		}

		accessToken, err := auth.MakeJWT(user.ID, user.IsAdmin, h.TokenSecret)
		if err != nil {
			HandleInternalServerError(w, r)
			log.Printf("%v", err)
			return
		}

		refreshToken, err := auth.MakeRefreshToken()
		if err != nil {
			HandleInternalServerError(w, r)
			log.Printf("%v", err)
			return
		}

		err = h.DB.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
			Token:     refreshToken,
			UserID:    user.ID,
			ExpiresAt: time.Now().UTC().AddDate(0, 0, 60),
		})
		if err != nil {
			HandleInternalServerError(w, r)
			log.Printf("%v", err)
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
		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    refreshToken,
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
