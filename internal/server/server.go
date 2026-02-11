// Package server provides server and route registering
package server

import (
	"log"
	"net/http"
	"time"

	"github.com/kairos4213/fithub/internal/config"
	"github.com/kairos4213/fithub/internal/handlers"
	"github.com/kairos4213/fithub/internal/middleware"
)

type Server struct {
	port    string
	fileDir string
	handler *handlers.Handler
	mw      *middleware.Middleware
}

func New(port, fileDir string, cfg *config.Config) *Server {
	return &Server{
		port:    port,
		fileDir: fileDir,
		handler: handlers.New(cfg),
		mw:      middleware.New(cfg),
	}
}

func (s *Server) Start() {
	mux := http.NewServeMux()
	s.registerRoutes(mux)

	srv := &http.Server{
		Addr:         ":" + s.port,
		Handler:      s.mw.Log(mux),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("Serving on port: %s\n", s.port)
	log.Fatal(srv.ListenAndServe())
}
