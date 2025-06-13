package http

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

type Server struct {
	router chi.Router
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Handler() http.Handler {
	return s.router
}

func (s *Server) setupMiddleware() {
	// TODO: Implement me.
}

func (s *Server) setupRoutes() {
	// TODO: Implement me.
}

func (s *Server) setupSwagger() {
	// TODO: Implement me.
}
