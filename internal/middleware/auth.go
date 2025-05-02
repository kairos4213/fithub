package middleware

import (
	"context"
	"net/http"

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
			contents := templates.LoginError("Please login")
			templates.Layout(contents, "FitHub", false).Render(r.Context(), w)
			return
		}

		accessToken := cookie.Value
		userID, err := auth.ValidateJWT(accessToken, mw.PublicKey)
		if err != nil {
			w.Header().Set("Content-type", "text/html")
			w.WriteHeader(http.StatusUnauthorized)
			contents := templates.LoginError("You are not authorized to view this page")
			templates.Layout(contents, "FitHub", false).Render(r.Context(), w)
			return
		}

		ctx := context.WithValue(r.Context(), cntx.UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
