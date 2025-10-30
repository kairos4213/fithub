package handlers

import "github.com/kairos4213/fithub/internal/database"

type Handler struct {
	DB          *database.Queries
	TokenSecret string
}
