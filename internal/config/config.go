package config

import (
	"database/sql"
	"log/slog"

	"github.com/kairos4213/fithub/internal/database"
)

type Config struct {
	DB          *database.Queries
	RawDB       *sql.DB
	Logger      *slog.Logger
	TokenSecret string
}

func New(db *database.Queries, rawDB *sql.DB, logger *slog.Logger, tokenSecret string) *Config {
	return &Config{DB: db, RawDB: rawDB, Logger: logger, TokenSecret: tokenSecret}
}
