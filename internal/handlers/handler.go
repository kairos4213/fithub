// Package handlers creates http handlers for server to use
package handlers

import (
	"github.com/kairos4213/fithub/internal/config"
)

type Handler struct {
	cfg *config.Config
}

func New(cfg *config.Config) *Handler {
	return &Handler{cfg: cfg}
}
