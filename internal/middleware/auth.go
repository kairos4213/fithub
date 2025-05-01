package middleware

import (
	"context"
	"net/http"

	"github.com/kairos4213/fithub/internal/auth"
	"github.com/kairos4213/fithub/internal/cntx"
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
			// TODO: Make an unauthorized / please login page
			w.Header().Set("Content-type", "text/html")
			w.Header().Set("HX-Location", `{"path": "/"}`)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		accessToken := cookie.Value
		userID, err := auth.ValidateJWT(accessToken, mw.PublicKey)
		if err != nil {
			// TODO: Make an unauthorized / please login page
			w.Header().Set("Content-type", "text/html")
			w.Header().Set("HX-Location", `{"path": "/"}`)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), cntx.UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
