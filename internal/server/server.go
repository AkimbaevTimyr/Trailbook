package server

import (
	"context"
	"fmt"
	"net/http"
)

// Server wraps an http.Server with graceful shutdown support.
type Server struct {
	httpServer *http.Server
}

// New creates a new HTTP server on the given port with the provided handler.
func New(port string, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:    fmt.Sprintf(":%s", port),
			Handler: handler,
		},
	}
}

// Start begins serving HTTP requests.
func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully stops the server.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
