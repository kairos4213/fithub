package config

import (
	"database/sql"
	"log/slog"

	"github.com/kairos4213/fithub/internal/database"
)

type OAuthProvider struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

type Config struct {
	DB          *database.Queries
	RawDB       *sql.DB
	Logger      *slog.Logger
	TokenSecret string
	OAuth       map[string]OAuthProvider
}

func New(db *database.Queries, rawDB *sql.DB, logger *slog.Logger, tokenSecret string, oauth map[string]OAuthProvider) *Config {
	return &Config{DB: db, RawDB: rawDB, Logger: logger, TokenSecret: tokenSecret, OAuth: oauth}
}
