package handlers

import (
	"net/http"
	"time"
)

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteDefaultMode,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
	})
	w.Header().Set("Content-type", "text/html")
	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusOK)
}
