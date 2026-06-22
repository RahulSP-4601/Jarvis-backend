package api

import (
	"net/http"

	"jarvis-backend/internal/ai"
	"jarvis-backend/internal/auth"
	"jarvis-backend/internal/config"
)

type Server struct {
	config config.Config
	ai     ai.Service
	auth   auth.Client
	mux    *http.ServeMux
}

func NewServer(cfg config.Config) *Server {
	mux := http.NewServeMux()
	server := &Server{
		config: cfg,
		ai:     ai.NewService(cfg),
		auth:   auth.NewClient(cfg),
		mux:    mux,
	}
	server.registerRoutes()

	return server
}

func (s *Server) Handler() http.Handler {
	return s.withCORS(s.mux)
}

func (s *Server) registerRoutes() {
	s.mux.HandleFunc("GET /", s.handleRoot)
	s.mux.HandleFunc("GET /favicon.ico", s.handleNoContent)
	s.mux.HandleFunc("GET /favicon.png", s.handleNoContent)
	s.mux.HandleFunc("GET /health", s.handleHealth)
	s.mux.HandleFunc("POST /v1/bootstrap", s.handleBootstrap)
	s.mux.HandleFunc("POST /v1/research", s.handleResearch)
	s.mux.HandleFunc("POST /v1/voice/command", s.handleVoiceCommand)
}
