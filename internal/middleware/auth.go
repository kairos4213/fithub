package middleware

import (
	"context"
	"fmt"
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

		accessCookie, err := r.Cookie("access_token")
		if err != nil {
			// Access cookie missing
			// Attempt to refresh access
			accessToken, errReason := mw.refreshAccessToken(w, r)
			if errReason != "" {
				http.Redirect(w, r, fmt.Sprintf("/unauthorized?reason=%v", errReason), http.StatusSeeOther)
			}
			// Successful access token refresh
			utils.SetAccessCookie(w, accessToken)
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

		accessToken := accessCookie.Value
		claims, err := auth.ValidateJWT(accessToken, mw.cfg.TokenSecret)
		if err != nil {
			// Access token expired
			if strings.Contains(err.Error(), "token is expired") {
				accessToken, errReason := mw.refreshAccessToken(w, r)
				if errReason != "" {
					http.Redirect(w, r, fmt.Sprintf("/unauthorized?reason=%v", errReason), http.StatusSeeOther)
				}
				// Successful access token refresh
				utils.SetAccessCookie(w, accessToken)
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

func (mw *Middleware) refreshAccessToken(w http.ResponseWriter, r *http.Request) (accessToken string, errReason string) {
	// Check for refresh token existence
	refreshCookie, err := r.Cookie("refresh_token")
	if err != nil {
		mw.cfg.Logger.Info("missing refresh token", slog.String("error", err.Error()))
		return "", "invalid_missing"
	}

	// Check for valid refresh token in db
	refreshToken := refreshCookie.Value
	user, err := mw.cfg.DB.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		utils.ClearCookies(w, refreshCookie)
		mw.cfg.Logger.Info("unable to fetch valid refresh token", slog.String("error", err.Error()))
		return "", "expired"
	}

	// Try to make new access token
	accessToken, err = auth.MakeJWT(user.ID, user.IsAdmin, mw.cfg.TokenSecret)
	if err != nil {
		// Error making access token -> need to remove any existing cookies and
		// revoke valid refresh token
		utils.ClearCookies(w, refreshCookie)

		revokeErr := mw.cfg.DB.RevokeRefreshToken(r.Context(), refreshToken)
		if revokeErr != nil {
			mw.cfg.Logger.Error("unable to revoke refresh token", slog.String("error", err.Error()))
			return "", "internal_error"
		}

		mw.cfg.Logger.Error("unable to make JWT", slog.String("error", err.Error()))
		return "", "internal_error"
	}

	return accessToken, ""
}
