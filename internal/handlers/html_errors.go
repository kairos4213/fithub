package handlers

import (
	"net/http"

	"github.com/kairos4213/fithub/internal/templates"
)

const (
	ServerErrMsg     = "Something went wrong. Please try later"
	AccessForbidden  = "You don't have permission to access this resource"
	NoAccessMsg      = "You don't have access to this! Please login, or register!"
	AccessExpiredMsg = "Access Expired. Please login."
)

func HandleInternalServerError(w http.ResponseWriter, r *http.Request) {
	htmlErr := templates.HtmlErr{Code: http.StatusInternalServerError, Msg: ServerErrMsg}
	contents := templates.ErrorDisplay(htmlErr)
	templates.Layout(contents, "FitHub", false).Render(r.Context(), w)
}

func HandleAccessForbiddenError(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")
	w.WriteHeader(http.StatusForbidden)

	htmlErr := templates.HtmlErr{Code: http.StatusForbidden, Msg: AccessForbidden}
	contents := templates.ErrorDisplay(htmlErr)
	templates.Layout(contents, "FitHub", true).Render(r.Context(), w)
}

func HandleUnauthorizedError(w http.ResponseWriter, r *http.Request, errMsg string) {
	htmlErr := templates.HtmlErr{Code: http.StatusUnauthorized, Msg: errMsg}
	contents := templates.ErrorDisplay(htmlErr)
	templates.Layout(contents, "FitHub", false).Render(r.Context(), w)
}
