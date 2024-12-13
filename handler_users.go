package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/kairos4213/fithub/internal/database"
)

type User struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		FirstName      string `json:"first_name"`
		MiddleName     string `json:"middle_name"`
		LastName       string `json:"last_name"`
		Email          string `json:"email"`
		HashedPassword string `json:"hashed_password"`
	}

	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding request", err)
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		FirstName:      params.FirstName,
		MiddleName:     sql.NullString{String: params.MiddleName},
		LastName:       params.LastName,
		Email:          params.Email,
		HashedPassword: params.HashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating user in database", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, user)
}
