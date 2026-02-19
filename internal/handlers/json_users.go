package handlers

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/kairos4213/fithub/internal/auth"
	"github.com/kairos4213/fithub/internal/cntx"
	"github.com/kairos4213/fithub/internal/database"
	"github.com/kairos4213/fithub/internal/utils"
	"github.com/kairos4213/fithub/internal/validate"
)

type User struct {
	ID         string `json:"id,omitempty"`
	FirstName  string `json:"first_name,omitempty"`
	MiddleName string `json:"middle_name,omitempty"`
	LastName   string `json:"last_name,omitempty"`
	Email      string `json:"email,omitempty"`
	Password   string `json:"password,omitempty"`
}

type response struct {
	User
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	reqParams := User{}
	if err := utils.ParseJSON(r, &reqParams); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "malformed request", err)
		return
	}

	if errs := validate.Fields(
		validate.Required(reqParams.FirstName, "first name"),
		validate.Required(reqParams.LastName, "last name"),
		validate.Required(reqParams.Email, "email"),
		validate.Required(reqParams.Password, "password"),
		validate.MinLen(reqParams.Password, 10, "password"),
		validate.MaxLen(reqParams.FirstName, 100, "first name"),
		validate.MaxLen(reqParams.LastName, 100, "last name"),
		validate.MaxLen(reqParams.Email, 255, "email"),
	); errs != nil {
		utils.RespondWithError(w, http.StatusBadRequest, errs[0].Error(), nil)
		return
	}

	hashedPassword, err := auth.HashPassword(reqParams.Password)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
		return
	}

	middleName := sql.NullString{Valid: false}
	if reqParams.MiddleName != "" {
		middleName.Valid = true
		middleName.String = strings.ToLower(reqParams.MiddleName)
	}

	user, err := h.cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		FirstName:      strings.ToLower(reqParams.FirstName),
		MiddleName:     middleName,
		LastName:       strings.ToLower(reqParams.LastName),
		Email:          strings.ToLower(reqParams.Email),
		HashedPassword: sql.NullString{String: hashedPassword, Valid: true},
	})
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error creating user in database", err)
		return
	}

	accessToken, refreshToken, err := h.issueSessionTokens(r.Context(), w, user.ID, user.IsAdmin)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error issuing session tokens", err)
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, response{
		User: User{
			ID:         user.ID.String(),
			FirstName:  user.FirstName,
			MiddleName: user.MiddleName.String,
			LastName:   user.LastName,
			Email:      user.Email,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	reqParams := User{}
	if err := utils.ParseJSON(r, &reqParams); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Malformed request", err)
		return
	}

	if errs := validate.Fields(
		validate.Required(reqParams.Email, "email"),
		validate.Required(reqParams.Password, "password"),
	); errs != nil {
		utils.RespondWithError(w, http.StatusBadRequest, errs[0].Error(), nil)
		return
	}

	user, err := h.cfg.DB.GetUser(r.Context(), reqParams.Email)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	if !user.HashedPassword.Valid {
		utils.RespondWithError(w, http.StatusUnauthorized, "Incorrect email or password", nil)
		return
	}

	match, err := auth.CheckPasswordHash(reqParams.Password, user.HashedPassword.String)
	if !match {
		utils.RespondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Bad Request: Invalid hash", err)
		return
	}

	accessToken, refreshToken, err := h.issueSessionTokens(r.Context(), w, user.ID, user.IsAdmin)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error issuing session tokens", err)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        user.ID.String(),
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	// TODO: Split this handler into two separate handlers
	// updateUsersHandlerPassword
	// updateUsersHandlerInfo (handles all other user information)
	userID, ok := cntx.UserID(r.Context())
	if !ok {
		utils.RespondWithError(w, http.StatusInternalServerError, "missing user id in context", nil)
		return
	}

	reqParams := User{}
	if err := utils.ParseJSON(r, &reqParams); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "malformed request", err)
		return
	}

	userParams := database.UpdateUserParams{
		ID: userID,
	}

	if reqParams.Password != "" {
		hashedPassword, err := auth.HashPassword(reqParams.Password)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error hashing password", err)
			return
		}
		userParams.HashedPassword.String = hashedPassword
		userParams.HashedPassword.Valid = true
	}

	if reqParams.Email != "" {
		userParams.Email.String = reqParams.Email
		userParams.Email.Valid = true
	}

	updatedUser, err := h.cfg.DB.UpdateUser(r.Context(), userParams)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error updating user", err)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, response{User: User{
		ID:        updatedUser.ID.String(),
		FirstName: updatedUser.FirstName,
		LastName:  updatedUser.LastName,
		Email:     updatedUser.Email,
	}})
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := cntx.UserID(r.Context())
	if !ok {
		utils.RespondWithError(w, http.StatusInternalServerError, "missing user id in context", nil)
		return
	}

	if err := h.cfg.DB.DeleteUser(r.Context(), userID); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error deleting user profile", err)
		return
	}

	utils.RespondWithJSON(w, http.StatusNoContent, User{})
}
