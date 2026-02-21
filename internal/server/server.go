// Package server provides server and route registering
package server

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/kairos4213/fithub/internal/config"
	"github.com/kairos4213/fithub/internal/handlers"
	"github.com/kairos4213/fithub/internal/middleware"
)

type Server struct {
	port    string
	fileDir string
	db      *sql.DB
	handler *handlers.Handler
	mw      *middleware.Middleware
}

func New(port, fileDir string, cfg *config.Config, db *sql.DB) *Server {
	return &Server{
		port:    port,
		fileDir: fileDir,
		db:      db,
		handler: handlers.New(cfg),
		mw:      middleware.New(cfg),
	}
}

func (s *Server) Start() {
	mux := http.NewServeMux()
	s.registerRoutes(mux)

	// 100 requests per minute per IP as a global safety net
	globalLimit := s.mw.RateLimit(100, time.Minute)

	srv := &http.Server{
		Addr:         ":" + s.port,
		Handler:      s.mw.Log(globalLimit(s.mw.SecureHeaders(s.mw.Cop(mux)))),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Create a context that listens for SIGINT/SIGTERM
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Start server in a goroutine
	go func() {
		log.Printf("Serving on port: %s\n", s.port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Block until signal received
	<-ctx.Done()
	stop()
	log.Println("Shutting down server...")

	// Allow 30 seconds for in-flight requests to complete
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	if err := s.db.Close(); err != nil {
		log.Printf("Error closing database: %v", err)
	}

	log.Println("Server exited gracefully")
}
