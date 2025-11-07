package middleware

import (
	"context"
	"log"
	"log/slog"
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

			claims, err := auth.ValidateJWT(accessToken, mw.cfg.TokenSecret)
			if err != nil {
				utils.RespondWithError(w, http.StatusUnauthorized, "Invalid JWT", err)
				return
			}

			ctx := context.WithValue(r.Context(), cntx.UserIDKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		// TODO: create refresh handler that gets redirected - Success -> Home
		//  Fail -> revoke and request re-login
		accessCookie, err := r.Cookie("access_token")
		if err != nil {
			http.Redirect(w, r, "/unauthorized?reason=invalid_missing", http.StatusSeeOther)
			mw.cfg.Logger.Info("missing access token", slog.String("error", err.Error()))
			return
		}
		accessToken := accessCookie.Value
		claims, err := auth.ValidateJWT(accessToken, mw.cfg.TokenSecret)
		if err != nil {
			// Access token expired
			if strings.Contains(err.Error(), "token is expired") {
				// Check for refresh token existence
				refreshCookie, err := r.Cookie("refresh_token")
				if err != nil {
					utils.ClearCookies(w, accessCookie)
					http.Redirect(w, r, "/unauthorized?reason=invalid_missing", http.StatusSeeOther)
					mw.cfg.Logger.Info("missing refresh token", slog.String("error", err.Error()))
					return
				}

				// Check for valid refresh token in db
				refreshToken := refreshCookie.Value
				user, err := mw.cfg.DB.GetUserFromRefreshToken(r.Context(), refreshToken)
				if err != nil {
					utils.ClearCookies(w, accessCookie, refreshCookie)
					http.Redirect(w, r, "/unauthorized?reason=expired", http.StatusSeeOther)
					mw.cfg.Logger.Info("unable to fetch valid refresh token", slog.String("error", err.Error()))
					return
				}

				// Try to make new access token
				accessToken, err = auth.MakeJWT(user.ID, user.IsAdmin, mw.cfg.TokenSecret)
				if err != nil {
					// Error making access token -> need to remove any existing cookies and
					// revoke valid refresh token
					utils.ClearCookies(w, accessCookie, refreshCookie)

					revokeErr := mw.cfg.DB.RevokeRefreshToken(r.Context(), refreshToken)
					if revokeErr != nil {
						http.Redirect(w, r, "/unauthorized?reason=internal_error", http.StatusSeeOther)
						mw.cfg.Logger.Error("unable to revoke refresh token", slog.String("error", err.Error()))
						return
					}

					http.Redirect(w, r, "/unauthorized?reason=internal_error", http.StatusSeeOther)
					mw.cfg.Logger.Error("unable to make JWT", slog.String("error", err.Error()))
					return
				}
				// Successful access token refresh
				http.SetCookie(w, &http.Cookie{
					Name:     "access_token",
					Value:    accessToken,
					Path:     "/",
					HttpOnly: true,
					Secure:   true,
					SameSite: http.SameSiteDefaultMode,
				})
				// Attempt creating claims
				claims, err := auth.ValidateJWT(accessToken, mw.cfg.TokenSecret)
				if err != nil {
					// clear access cookie and redirect to try again
					utils.ClearCookies(w, accessCookie)
					http.Redirect(w, r, "/unauthorized?reason=internal_error", http.StatusSeeOther)
					mw.cfg.Logger.Error("unable to validate JWT", slog.String("error", err.Error()))
					return
				}
				ctx := context.WithValue(r.Context(), cntx.UserIDKey, claims.UserID)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
			// Invalid access token -> clear access token cookie
			utils.ClearCookies(w, accessCookie)
			http.Redirect(w, r, "/unauthorized?reason=invalid_missing", http.StatusSeeOther)
			mw.cfg.Logger.Info("invalid access token", slog.String("error", err.Error()))
			return
		}

		// Valid access token
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

			claims, err := auth.ValidateJWT(accessToken, mw.cfg.TokenSecret)
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

		accessCookie, err := r.Cookie("access_token")
		if err != nil {
			http.Redirect(w, r, "/unauthorized?reason=invalid_missing", http.StatusSeeOther)
			mw.cfg.Logger.Info("missing access token", slog.String("error", err.Error()))
			return
		}

		accessToken := accessCookie.Value
		claims, err := auth.ValidateJWT(accessToken, mw.cfg.TokenSecret)
		if err != nil {
			// Access token expired
			if strings.Contains(err.Error(), "token is expired") {
				// Check for refresh token existence
				refreshCookie, err := r.Cookie("refresh_token")
				if err != nil {
					utils.ClearCookies(w, accessCookie)
					http.Redirect(w, r, "/unauthorized?reason=invalid_missing", http.StatusSeeOther)
					mw.cfg.Logger.Info("missing refresh token", slog.String("error", err.Error()))
					return
				}

				// Check for valid refresh token in db
				refreshToken := refreshCookie.Value
				user, err := mw.cfg.DB.GetUserFromRefreshToken(r.Context(), refreshToken)
				if err != nil {
					utils.ClearCookies(w, accessCookie, refreshCookie)
					http.Redirect(w, r, "/unauthorized?reason=expired", http.StatusSeeOther)
					mw.cfg.Logger.Info("unable to fetch valid refresh token", slog.String("error", err.Error()))
					return
				}

				// Try to make new access token
				accessToken, err = auth.MakeJWT(user.ID, user.IsAdmin, mw.cfg.TokenSecret)
				if err != nil {
					// Error making access token -> need to remove any existing cookies and
					// revoke valid refresh token
					utils.ClearCookies(w, accessCookie, refreshCookie)

					revokeErr := mw.cfg.DB.RevokeRefreshToken(r.Context(), refreshToken)
					if revokeErr != nil {
						http.Redirect(w, r, "/unauthorized?reason=internal_error", http.StatusSeeOther)
						mw.cfg.Logger.Error("unable to revoke refresh token", slog.String("error", err.Error()))
						return
					}

					http.Redirect(w, r, "/unauthorized?reason=internal_error", http.StatusSeeOther)
					mw.cfg.Logger.Error("unable to make JWT", slog.String("error", err.Error()))
					return
				}
				// Successful access token refresh
				http.SetCookie(w, &http.Cookie{
					Name:     "access_token",
					Value:    accessToken,
					Path:     "/",
					HttpOnly: true,
					Secure:   true,
					SameSite: http.SameSiteDefaultMode,
				})
				// Attempt creating claims
				claims, err := auth.ValidateJWT(accessToken, mw.cfg.TokenSecret)
				if err != nil {
					// clear access cookie and redirect to try again
					utils.ClearCookies(w, accessCookie)
					http.Redirect(w, r, "/unauthorized?reason=internal_error", http.StatusSeeOther)
					mw.cfg.Logger.Error("unable to validate JWT", slog.String("error", err.Error()))
					return
				}

				if !claims.IsAdmin {
					http.Redirect(w, r, "/forbidden", http.StatusSeeOther)
					mw.cfg.Logger.Warn("unauthorized admin request",
						slog.String("user", claims.UserID.String()),
						slog.String("method", r.Method),
						slog.String("path", r.URL.Path),
					)
					return
				}

				ctx := context.WithValue(r.Context(), cntx.UserIDKey, claims.UserID)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
			// Invalid access token -> clear access token cookie
			utils.ClearCookies(w, accessCookie)
			http.Redirect(w, r, "/unauthorized?reason=invalid_missing", http.StatusSeeOther)
			mw.cfg.Logger.Info("invalid access token", slog.String("error", err.Error()))
			return
		}
		// Valid access token

		if !claims.IsAdmin {
			http.Redirect(w, r, "/forbidden", http.StatusSeeOther)
			mw.cfg.Logger.Warn("unauthorized admin request",
				slog.String("user", claims.UserID.String()),
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
			)
			return
		}

		ctx := context.WithValue(r.Context(), cntx.UserIDKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
