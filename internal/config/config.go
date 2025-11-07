package config

import (
	"log/slog"

	"github.com/kairos4213/fithub/internal/database"
)

type Config struct {
	DB          *database.Queries
	Logger      *slog.Logger
	TokenSecret string
}

func New(db *database.Queries, logger *slog.Logger, tokenSecret string) *Config {
	return &Config{DB: db, Logger: logger, TokenSecret: tokenSecret}
}
