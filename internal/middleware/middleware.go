package middleware

import (
	"github.com/kairos4213/fithub/internal/config"
)

type Middleware struct {
	cfg *config.Config
}

func New(cfg *config.Config) *Middleware {
	return &Middleware{cfg: cfg}
}
