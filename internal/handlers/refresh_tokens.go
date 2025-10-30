package handlers

import (
	"net/http"

	"github.com/kairos4213/fithub/internal/auth"
	"github.com/kairos4213/fithub/internal/utils"
)

func (h *Handler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	type response struct {
		AccesToken string `json:"access_token"`
	}

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Error getting bearer token", err)
		return
	}

	user, err := h.DB.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Couldn't match user with refresh token", err)
		return
	}

	accessToken, err := auth.MakeJWT(user.ID, user.IsAdmin, h.TokenSecret)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error making JWT", err)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, response{AccesToken: accessToken})
}

func (h *Handler) RevokeToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Error getting bearer token", err)
		return
	}

	err = h.DB.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error revoking token", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
