package server

import (
	"context"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"defect-tracker/internal/pkg/config"
)

// HTTPServer wraps the standard http.Server with structured logging helpers.
type HTTPServer struct {
	srv *http.Server
	log *zap.Logger
}

func NewHTTPServer(cfg config.Config, handler http.Handler, log *zap.Logger) *HTTPServer {
	address := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)

	return &HTTPServer{
		log: log,
		srv: &http.Server{
			Addr:         address,
			Handler:      handler,
			ReadTimeout:  cfg.Server.ReadTimeout,
			WriteTimeout: cfg.Server.WriteTimeout,
			IdleTimeout:  cfg.Server.IdleTimeout,
		},
	}
}

func (s *HTTPServer) Start() error {
	s.log.Info("starting http server", zap.String("addr", s.srv.Addr))
	err := s.srv.ListenAndServe()
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}

func (s *HTTPServer) Shutdown(ctx context.Context) error {
	s.log.Info("shutting down http server")
	return s.srv.Shutdown(ctx)
}
