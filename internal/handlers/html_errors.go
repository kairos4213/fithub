package handlers

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/kairos4213/fithub/internal/templates"
	"github.com/kairos4213/fithub/internal/validate"
)

const (
	ServerErrMsg      = "Something went wrong. Please try later"
	AccessForbidden   = "You don't have permission to access this resource"
	NoAccessMsg       = "You don't have access to this! Please login, or register!"
	AccessExpiredMsg  = "Access Expired. Please login."
	LoginFailMsg      = "Username and/or password are incorrect. Please try again."
	DuplicateEmailMsg = "That email already exists!"
)

func HandleInternalServerError(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")
	w.Header().Set("HX-Retarget", "body")
	w.Header().Set("HX-Reswap", "outerHTML")
	w.WriteHeader(http.StatusInternalServerError)

	htmlErr := templates.HtmlErr{Code: http.StatusInternalServerError, Msg: ServerErrMsg}
	err := templates.ErrorDisplay(htmlErr).Render(r.Context(), w)
	if err != nil {
		http.Error(w, ServerErrMsg, http.StatusInternalServerError)
	}
}

func GetForbiddenPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")
	w.Header().Set("HX-Retarget", "body")
	w.Header().Set("HX-Reswap", "outerHTML")
	w.WriteHeader(http.StatusForbidden)

	htmlErr := templates.HtmlErr{Code: http.StatusForbidden, Msg: AccessForbidden}
	contents := templates.ErrorDisplay(htmlErr)
	err := templates.Layout(contents, "FitHub", true).Render(r.Context(), w)
	if err != nil {
		HandleInternalServerError(w, r)
		return
	}
}

func GetUnauthorizedPage(w http.ResponseWriter, r *http.Request) {
	errStr := r.URL.Query().Get("reason")
	errMsg := NoAccessMsg
	if errStr == "expired" {
		errMsg = AccessExpiredMsg
	}

	w.Header().Set("Content-type", "text/html")
	w.Header().Set("HX-Retarget", "body")
	w.Header().Set("HX-Reswap", "outerHTML")
	w.WriteHeader(http.StatusUnauthorized)

	htmlErr := templates.HtmlErr{Code: http.StatusUnauthorized, Msg: errMsg}
	contents := templates.ErrorDisplay(htmlErr)
	err := templates.Layout(contents, "FitHub", false).Render(r.Context(), w)
	if err != nil {
		HandleInternalServerError(w, r)
		return
	}
}

func HandleLoginFailure(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")
	w.WriteHeader(http.StatusUnprocessableEntity)

	htmlErr := templates.HtmlErr{Code: http.StatusUnprocessableEntity, Msg: LoginFailMsg}
	err := templates.LoginFailure(htmlErr).Render(r.Context(), w)
	if err != nil {
		HandleInternalServerError(w, r)
		return
	}
}

func HandleRegPageEmailAlert(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")
	w.WriteHeader(http.StatusConflict)

	htmlErr := templates.HtmlErr{Code: http.StatusConflict, Msg: DuplicateEmailMsg}
	err := templates.RegPageEmailAlert(htmlErr).Render(r.Context(), w)
	if err != nil {
		HandleInternalServerError(w, r)
		return
	}
}

func HandleBadRequest(w http.ResponseWriter, r *http.Request, errMsg string) {
	w.Header().Set("Content-type", "text/html")
	w.WriteHeader(http.StatusBadRequest)

	htmlErr := templates.HtmlErr{Code: http.StatusBadRequest, Msg: errMsg}
	err := templates.FormError(htmlErr).Render(r.Context(), w)
	if err != nil {
		HandleInternalServerError(w, r)
		return
	}
}

// HandleScopedFieldErrors renders field errors and retargets the response to a
// specific form-error container. Pass an empty formErrorID to skip retargeting
// (e.g. when the form already has a global #form-error target).
func HandleScopedFieldErrors(w http.ResponseWriter, r *http.Request, logger *slog.Logger, errs []validate.FieldError, fields []string, prefix string, formErrorID string) {
	if formErrorID != "" {
		w.Header().Set("HX-Retarget", "#"+formErrorID)
		w.Header().Set("HX-Reswap", "innerHTML")
	}
	HandleFieldErrors(w, r, logger, errs, fields, prefix)
}

func HandleFieldErrors(w http.ResponseWriter, r *http.Request, logger *slog.Logger, errs []validate.FieldError, fields []string, prefix string) {
	errMap := make(map[string]string)
	for _, e := range errs {
		slug := prefix + strings.ReplaceAll(e.Field, " ", "-")
		if _, exists := errMap[slug]; !exists {
			errMap[slug] = e.Message
		}
	}

	w.Header().Set("Content-type", "text/html")
	w.WriteHeader(http.StatusBadRequest)

	err := templates.FieldErrors(errMap, fields).Render(r.Context(), w)
	if err != nil {
		logger.Error("failed to render field errors", slog.String("error", err.Error()))
	}
}
