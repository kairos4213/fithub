package config

import "github.com/kairos4213/fithub/internal/database"

type Config struct {
	DB          *database.Queries
	TokenSecret string
}

func New(db *database.Queries, tokenSecret string) *Config {
	return &Config{DB: db, TokenSecret: tokenSecret}
}
