package handlers

import (
	"log"
	"net/http"

	"github.com/kairos4213/fithub/internal/utils"
)

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	accessCookie, err := r.Cookie("access_token")
	if err != nil {
		HandleInternalServerError(w, r)
		log.Printf("%v", err)
		return
	}

	refreshCookie, err := r.Cookie("refresh_token")
	if err != nil {
		HandleInternalServerError(w, r)
		log.Printf("%v", err)
		return
	}
	err = h.DB.RevokeRefreshToken(r.Context(), refreshCookie.Value)
	if err != nil {
		HandleInternalServerError(w, r)
		log.Printf("%v", err)
		return
	}

	utils.ClearCookies(w, accessCookie, refreshCookie)

	w.Header().Set("Content-type", "text/html")
	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusOK)
}
