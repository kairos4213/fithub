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
	if r.Method == "GET" {
		contents := templates.LoginPage(templates.HtmlErr{})
		templates.Layout(contents, "FitHub | Login", false).Render(r.Context(), w)
		return
	}

	if r.Method == "POST" {
		email := r.FormValue("email")
		password := r.FormValue("password")

		user, err := h.DB.GetUser(r.Context(), email)
		if err != nil {
			w.Header().Set("Content-type", "text/html")
			w.WriteHeader(http.StatusUnauthorized)

			htmlErr := templates.HtmlErr{Code: http.StatusUnauthorized, Msg: "Username and/or password are incorrect. Please try again."}
			templates.LoginPage(htmlErr).Render(r.Context(), w)

			log.Print("User email does not exist")
			return
		}

		if err = auth.CheckPasswordHash(password, user.HashedPassword); err != nil {
			w.Header().Set("Content-type", "text/html")
			w.WriteHeader(http.StatusUnauthorized)

			htmlErr := templates.HtmlErr{Code: http.StatusUnauthorized, Msg: "Username and/or password are incorrect. Please try again."}
			templates.LoginPage(htmlErr).Render(r.Context(), w)

			log.Print("Incorrect password entered")
			return
		}

		accessToken, err := auth.MakeJWT(user.ID, user.IsAdmin, h.PrivateKey)
		if err != nil {
			w.Header().Set("Content-type", "text/html")
			w.WriteHeader(http.StatusInternalServerError)

			htmlErr := templates.HtmlErr{Code: http.StatusInternalServerError, Msg: "Something went wrong! Please try later."}
			templates.ErrorDisplay(htmlErr).Render(r.Context(), w)

			log.Printf("%v", err)
			return
		}

		refreshToken, err := auth.MakeRefreshToken()
		if err != nil {
			w.Header().Set("Content-type", "text/html")
			w.WriteHeader(http.StatusInternalServerError)

			htmlErr := templates.HtmlErr{Code: http.StatusInternalServerError, Msg: "Something went wrong! Please try later."}
			templates.ErrorDisplay(htmlErr).Render(r.Context(), w)

			log.Printf("%v", err)
			return
		}

		err = h.DB.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
			Token:     refreshToken,
			UserID:    user.ID,
			ExpiresAt: time.Now().UTC().AddDate(0, 0, 60),
		})
		if err != nil {
			w.Header().Set("Content-type", "text/html")
			w.WriteHeader(http.StatusInternalServerError)

			htmlErr := templates.HtmlErr{Code: http.StatusInternalServerError, Msg: "Something went wrong! Please try later."}
			templates.ErrorDisplay(htmlErr).Render(r.Context(), w)

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

		if user.IsAdmin {
			w.Header().Set("HX-Location", `{"path": "/admin"}`)
			w.WriteHeader(http.StatusAccepted)
		}

		w.Header().Set("HX-Location", `{"path": "/workouts"}`)
		w.WriteHeader(http.StatusAccepted)
	}
}
