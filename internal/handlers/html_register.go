package handlers

import (
	"log"
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
		templates.Layout(contents, "FitHub | Register", false).Render(r.Context(), w)
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
			log.Printf("%v", err)
			return
		}

		user, err := h.DB.CreateUser(r.Context(), database.CreateUserParams{
			FirstName:      firstName,
			LastName:       lastName,
			Email:          email,
			HashedPassword: hashedPassword,
		})
		if err != nil {
			if strings.Contains(err.Error(), "users_email_key") {
				w.Header().Set("Content-type", "text/html")
				w.WriteHeader(http.StatusConflict)

				regErr := templates.HtmlErr{Code: http.StatusConflict, Msg: "That email already exists! Please try again."}
				templates.RegPageEmailAlert(regErr, email).Render(r.Context(), w)

				log.Printf("DB Duplicate email error: %v", err)
				return
			}
		}

		accessToken, err := auth.MakeJWT(user.ID, user.IsAdmin, h.PrivateKey)
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
		w.Header().Set("Content-type", "text/html")
		w.Header().Set("HX-Location", `{"path": "/workouts"}`)
		w.WriteHeader(http.StatusCreated)
	}
}

func (h *Handler) CheckUserEmail(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")

	_, err := h.DB.GetUser(r.Context(), email)
	if err == nil {
		w.Header().Set("Content-type", "text/html")
		w.WriteHeader(http.StatusConflict)

		htmlErr := templates.HtmlErr{Code: http.StatusConflict, Msg: "That email already exists!"}
		templates.RegPageEmailAlert(htmlErr, email).Render(r.Context(), w)

		log.Print("User email already exists")
		return
	}

	templates.RegPageEmailAlert(templates.HtmlErr{}, email).Render(r.Context(), w)
}
