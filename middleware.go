package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/kairos4213/fithub/internal/auth"
)

type authedHandler func(http.ResponseWriter, *http.Request, uuid.UUID)

func (cfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
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

		handler(w, r, userID)
	})
}
