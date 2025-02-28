package main

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kairos4213/fithub/internal/auth"
	"github.com/kairos4213/fithub/internal/database"
)

type User struct {
	ID         uuid.UUID `json:"id,omitempty"`
	FirstName  string    `json:"first_name,omitempty"`
	MiddleName string    `json:"middle_name,omitempty"`
	LastName   string    `json:"last_name,omitempty"`
	Email      string    `json:"email,omitempty"`
}

type response struct {
	User
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

func (cfg *apiConfig) createUsersHandler(w http.ResponseWriter, r *http.Request) {
	type requestParams struct {
		FirstName  string `json:"first_name"`
		MiddleName string `json:"middle_name"`
		LastName   string `json:"last_name"`
		Email      string `json:"email"`
		Password   string `json:"password"`
	}

	reqParams := requestParams{}
	if err := parseJSON(r, &reqParams); err != nil {
		respondWithError(w, http.StatusBadRequest, "malformed request", err)
		return
	}

	if reqParams.FirstName == "" {
		respondWithError(w, http.StatusBadRequest, "Missing First Name", errors.New("malformed request"))
		return
	}

	if reqParams.LastName == "" {
		respondWithError(w, http.StatusBadRequest, "Missing Last Name", errors.New("malformed request"))
		return
	}

	if reqParams.Email == "" {
		respondWithError(w, http.StatusBadRequest, "Missing Email", errors.New("malformed request"))
		return
	}

	if len(reqParams.Password) < 10 {
		respondWithError(w, http.StatusBadRequest, "Password must be at least 10 characters", errors.New("malformed request"))
		return
	}

	hashedPassword, err := auth.HashPassword(reqParams.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
		return
	}

	middleName := sql.NullString{Valid: false}
	if reqParams.MiddleName != "" {
		middleName.Valid = true
		middleName.String = strings.ToLower(reqParams.MiddleName)
	}

	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		FirstName:      strings.ToLower(reqParams.FirstName),
		MiddleName:     middleName,
		LastName:       strings.ToLower(reqParams.LastName),
		Email:          strings.ToLower(reqParams.Email),
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating user in database", err)
		return
	}

	accessToken, err := auth.MakeJWT(user.ID, cfg.privateKey)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create access token", err)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create refresh token", err)
		return
	}

	err = cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().UTC().AddDate(0, 0, 60),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error storing refresh token", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, response{
		User: User{
			ID:         user.ID,
			FirstName:  user.FirstName,
			MiddleName: user.MiddleName.String,
			LastName:   user.LastName,
			Email:      user.Email,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func (cfg *apiConfig) loginUsersHandler(w http.ResponseWriter, r *http.Request) {
	type requestParams struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	reqParams := requestParams{}
	if err := parseJSON(r, &reqParams); err != nil {
		respondWithError(w, http.StatusBadRequest, "Malformed request", err)
		return
	}

	if reqParams.Email == "" {
		respondWithError(w, http.StatusBadRequest, "Missing Email", errors.New("malformed request"))
		return
	}

	if reqParams.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Missing Password", errors.New("malformed request"))
		return
	}

	user, err := cfg.db.GetUser(r.Context(), reqParams.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	err = auth.CheckPasswordHash(reqParams.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	accessToken, err := auth.MakeJWT(user.ID, cfg.privateKey)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create access token", err)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create refresh token", err)
		return
	}

	err = cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().UTC().AddDate(0, 0, 60),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error storing refresh token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func (cfg *apiConfig) updateUsersHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Split this handler into two separate handlers
	// updateUsersHandlerPassword
	// updateUsersHandlerInfo (handles all other user information)
	type requestParams struct {
		Email    *string `json:"email,omitempty"`
		Password *string `json:"password,omitempty"`
	}

	userID := r.Context().Value(userIDKey).(uuid.UUID)

	reqParams := requestParams{}
	if err := parseJSON(r, &reqParams); err != nil {
		respondWithError(w, http.StatusBadRequest, "malformed request", err)
		return
	}

	userParams := database.UpdateUserParams{
		ID: userID,
	}

	// TODO: refactor?
	if reqParams.Password != nil {
		hashedPassword, err := auth.HashPassword(*reqParams.Password)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error hashing password", err)
			return
		}
		userParams.HashedPassword.String = hashedPassword
		userParams.HashedPassword.Valid = true
	}

	if reqParams.Email != nil {
		userParams.Email.String = *reqParams.Email
		userParams.Email.Valid = true
	}

	updatedUser, err := cfg.db.UpdateUser(r.Context(), userParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error updating user", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{User: User{
		ID:        updatedUser.ID,
		FirstName: updatedUser.FirstName,
		LastName:  updatedUser.LastName,
		Email:     updatedUser.Email,
	}})
}

func (cfg *apiConfig) deleteUsersHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIDKey).(uuid.UUID)

	if err := cfg.db.DeleteUser(r.Context(), userID); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error deleting user profile", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, User{})
}
