package utils

import (
	"net/http"
	"time"
)

func ClearCookies(w http.ResponseWriter, cookies ...*http.Cookie) {
	for _, cookie := range cookies {
		http.SetCookie(w, &http.Cookie{
			Name:     cookie.Name,
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteDefaultMode,
			Expires:  time.Unix(0, 0),
			MaxAge:   -1,
		})
	}
}
