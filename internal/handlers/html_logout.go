package handlers

import (
	"log/slog"
	"net/http"

	"github.com/kairos4213/fithub/internal/utils"
)

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	accessCookie, err := r.Cookie("access_token")
	if err != nil {
		h.cfg.Logger.Info("failed to find access token", slog.String("error", err.Error()))
	}

	refreshCookie, err := r.Cookie("refresh_token")
	if err != nil {
		h.cfg.Logger.Info("failed to find refresh token", slog.String("error", err.Error()))
	}

	if refreshCookie != nil {
		err = h.cfg.DB.RevokeRefreshToken(r.Context(), refreshCookie.Value)
		if err != nil {
			h.cfg.Logger.Error("failed to revoke refresh token", slog.String("refresh_token", refreshCookie.Value), slog.String("error", err.Error()))
		}
	}

	utils.ClearCookies(w, accessCookie, refreshCookie)

	w.Header().Set("Content-type", "text/html")
	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusOK)
}
