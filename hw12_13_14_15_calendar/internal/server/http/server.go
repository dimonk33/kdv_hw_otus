package internalhttp

import (
	"context"
	"net/http"
	"time"
)

type Server struct {
	logger Logger
	server *http.Server
}

type Logger interface {
	Info(msg string)
}

type Application interface {
	HelloHandler(w http.ResponseWriter, r *http.Request)
}

func NewServer(addr string, logger Logger, app Application) *Server {
	router := http.NewServeMux()
	router.HandleFunc("/hello", app.HelloHandler)

	m := Middleware{
		logger: logger,
	}

	s := &Server{
		logger: logger,
	}
	s.server = &http.Server{
		Addr:              addr,
		Handler:           m.logging(router),
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
	}
	return s
}

func (s *Server) Start(ctx context.Context) error {
	if err := s.server.ListenAndServe(); err != nil {
		return err
	}
	s.logger.Info("Http сервер запущен")

	<-ctx.Done()

	return ctx.Err()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
