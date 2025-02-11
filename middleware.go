package main

import (
	"context"
	"net/http"

	"github.com/kairos4213/fithub/internal/auth"
)

type contextKey string

const userIDKey contextKey = "userID"

func (cfg *apiConfig) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken, err := auth.GetBearerToken(r.Header)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Missing JWT", err)
			return
		}

		userID, err := auth.ValidateJWT(accessToken, cfg.publicKey)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Invalid JWT", err)
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
