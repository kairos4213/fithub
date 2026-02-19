// Package handlers creates http handlers for server to use
package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/kairos4213/fithub/internal/auth"
	"github.com/kairos4213/fithub/internal/config"
	"github.com/kairos4213/fithub/internal/database"
	"github.com/kairos4213/fithub/internal/utils"
)

type Handler struct {
	cfg *config.Config
}

func New(cfg *config.Config) *Handler {
	return &Handler{cfg: cfg}
}

// issueSessionTokens creates a JWT access token and refresh token, stores the
// refresh token in the database, and sets both as HTTP cookies.
func (h *Handler) issueSessionTokens(ctx context.Context, w http.ResponseWriter, userID uuid.UUID, isAdmin bool) (accessToken, refreshToken string, err error) {
	accessToken, err = auth.MakeJWT(userID, isAdmin, h.cfg.TokenSecret)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = auth.MakeRefreshToken()
	if err != nil {
		return "", "", err
	}

	err = h.cfg.DB.CreateRefreshToken(ctx, database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    userID,
		ExpiresAt: time.Now().UTC().AddDate(0, 0, 60),
	})
	if err != nil {
		return "", "", err
	}

	utils.SetAccessCookie(w, accessToken)
	utils.SetRefreshCookie(w, refreshToken)
	return accessToken, refreshToken, nil
}
