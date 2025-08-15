package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

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

			claims, err := auth.ValidateJWT(accessToken, mw.PublicKey)
			if err != nil {
				utils.RespondWithError(w, http.StatusUnauthorized, "Invalid JWT", err)
				return
			}

			ctx := context.WithValue(r.Context(), cntx.UserIDKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		cookie, err := r.Cookie("access_token")
		if err != nil {
			http.Redirect(w, r, "/unauthorized?reason=invalid_missing", http.StatusSeeOther)
			log.Printf("%v", err)
			return
		}
		accessToken := cookie.Value
		claims, err := auth.ValidateJWT(accessToken, mw.PublicKey)
		if err != nil {
			if strings.Contains(err.Error(), "token is expired") {
				http.Redirect(w, r, "/unauthorized?reason=expired", http.StatusSeeOther)
				log.Printf("%v", err)
				return
			}

			http.Redirect(w, r, "/unauthorized?reason=invalid_missing", http.StatusSeeOther)
			log.Printf("%v", err)
			return
		}

		ctx := context.WithValue(r.Context(), cntx.UserIDKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (mw *Middleware) AdminAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Accept")
		if header == "application/json" {
			accessToken, err := auth.GetBearerToken(r.Header)
			if err != nil {
				utils.RespondWithError(w, http.StatusUnauthorized, "Missing JWT", err)
				return
			}

			claims, err := auth.ValidateJWT(accessToken, mw.PublicKey)
			if err != nil {
				utils.RespondWithError(w, http.StatusUnauthorized, "Invalid JWT", err)
				return
			}

			if !claims.IsAdmin {
				utils.RespondWithError(w, http.StatusForbidden, "You don't have permission to view this!", err)
				log.Println("Unauthorized admin request:")
				log.Printf("\tUser ID: %v", claims.UserID)
				log.Printf("\tRequest type: %v", r.Method)
				log.Printf("\tRequest body: %v", r.Body)
				return
			}

			ctx := context.WithValue(r.Context(), cntx.UserIDKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		cookie, err := r.Cookie("access_token")
		if err != nil {
			http.Redirect(w, r, "/unauthorized?reason=invalid_missing", http.StatusSeeOther)
			log.Printf("%v", err)
			return
		}

		accessToken := cookie.Value
		claims, err := auth.ValidateJWT(accessToken, mw.PublicKey)
		if err != nil {
			if strings.Contains(err.Error(), "token expired") {
				http.Redirect(w, r, "/unauthorized?reason=expired", http.StatusSeeOther)
				log.Printf("%v", err)
				return
			}

			http.Redirect(w, r, "/unauthorized?reason=invalid_missing", http.StatusSeeOther)
			log.Printf("%v", err)
			return
		}

		if !claims.IsAdmin {
			http.Redirect(w, r, "/forbidden", http.StatusSeeOther)
			log.Println("Unauthorized admin request:")
			log.Printf("\tUser ID: %v", claims.UserID)
			log.Printf("\tRequest type: %v", r.Method)
			log.Printf("\tRequest body: %v", r.Body)
			return
		}

		ctx := context.WithValue(r.Context(), cntx.UserIDKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
