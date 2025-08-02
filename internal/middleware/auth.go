package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/kairos4213/fithub/internal/auth"
	"github.com/kairos4213/fithub/internal/cntx"
	"github.com/kairos4213/fithub/internal/templates"
	"github.com/kairos4213/fithub/internal/utils"
)

func (mw *Middleware) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Accept")
		if header == "application/json" {
			accessToken, err := auth.GetBearerToken(r.Header)
			if err != nil {
				utils.RespondWithError(w, http.StatusUnauthorized, "Missing JWT", err)
				return
			}

			userID, err := auth.ValidateJWT(accessToken, mw.PublicKey)
			if err != nil {
				utils.RespondWithError(w, http.StatusUnauthorized, "Invalid JWT", err)
				return
			}

			ctx := context.WithValue(r.Context(), cntx.UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		cookie, err := r.Cookie("access_token")
		if err != nil {
			w.Header().Set("Content-type", "text/html")
			w.WriteHeader(http.StatusUnauthorized)

			htmlErr := templates.HtmlErr{Code: http.StatusUnauthorized, Msg: "You don't have access to this! Please login, or register!"}
			contents := templates.ErrorDisplay(htmlErr)
			templates.Layout(contents, "FitHub", false).Render(r.Context(), w)

			log.Printf("%v", err)
			return
		}

		accessToken := cookie.Value
		userID, err := auth.ValidateJWT(accessToken, mw.PublicKey)
		if err != nil {
			w.Header().Set("Content-type", "text/html")
			w.WriteHeader(http.StatusUnauthorized)

			if strings.Contains(err.Error(), "token expired") {
				htmlErr := templates.HtmlErr{Code: http.StatusUnauthorized, Msg: "Access Expired. Please login again"}
				contents := templates.ErrorDisplay(htmlErr)
				templates.Layout(contents, "FitHub", false).Render(r.Context(), w)

				log.Printf("%v", err)
				return
			}

			htmlErr := templates.HtmlErr{Code: http.StatusUnauthorized, Msg: "You don't have access to this! Please login, or register!"}
			contents := templates.ErrorDisplay(htmlErr)
			templates.Layout(contents, "FitHub", false).Render(r.Context(), w)

			log.Printf("%v", err)
			return
		}

		ctx := context.WithValue(r.Context(), cntx.UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
