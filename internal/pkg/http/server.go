package http

import (
	"context"
	"net/http"
	"time"

	"github.com/yrss1/workout/internal/pkg/http/config"
	"github.com/yrss1/workout/internal/pkg/log"
)

const (
	readHeaderTimeout = 20 * time.Second // Prevent Slowloris attacks
	readTimeout       = 30 * time.Second
	writeTimeout      = 30 * time.Second
	idleTimeout       = 60 * time.Second
)

type Server struct {
	server  *http.Server
	address string
	handler http.Handler
	logger  *log.Log
}

func NewServer(config config.HTTPConfig, log *log.Log, handler http.Handler) *Server {
	return &Server{
		address: config.BindAddr,
		handler: handler,
		logger:  log,
	}
}

func (s *Server) Start(ctx context.Context) error {
	s.server = &http.Server{
		Addr:              s.address,
		Handler:           s.handler,
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
	}
	s.logger.Info("Starting the server.")
	err := s.server.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("Stopping the server.")
	err := s.server.Shutdown(ctx)
	if err != nil {
		s.logger.Error("Error while shutting down the server.")
		return err
	}
	return nil
}
