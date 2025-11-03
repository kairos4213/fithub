package middleware

import "github.com/kairos4213/fithub/internal/database"

type Middleware struct {
	DB          *database.Queries
	TokenSecret string
}
