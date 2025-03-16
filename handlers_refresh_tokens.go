package main

import (
	"net/http"

	"github.com/kairos4213/fithub/internal/auth"
)

func (cfg api) refreshHandler(w http.ResponseWriter, r *http.Request) {
	type response struct {
		AccesToken string `json:"access_token"`
	}

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error getting bearer token", err)
		return
	}

	user, err := cfg.db.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't match user with refresh token", err)
		return
	}

	accessToken, err := auth.MakeJWT(user.ID, cfg.privateKey)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error making JWT", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{AccesToken: accessToken})
}

func (cfg api) revokeHandler(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error getting bearer token", err)
		return
	}

	err = cfg.db.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error revoking token", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
