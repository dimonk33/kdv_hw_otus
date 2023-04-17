package internalhttp

import (
	"context"
	"net/http"
	"time"

	gs "github.com/dimonk33/kdv_hw_otus/hw12_13_14_15_calendar/api/gen"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

type Server struct {
	addr       string
	logger     Logger
	grpcServer gs.EventsServer
	server     *http.Server
}

type Logger interface {
	Info(msg string)
}

func NewServer(addr string, logger Logger, grpcServer gs.EventsServer) *Server {
	s := &Server{
		addr:       addr,
		logger:     logger,
		grpcServer: grpcServer,
	}
	return s
}

func (s *Server) Start(ctx context.Context) error {
	mux := runtime.NewServeMux()
	err := gs.RegisterEventsHandlerServer(ctx, mux, s.grpcServer)
	if err != nil {
		return err
	}
	m := Middleware{
		logger: s.logger,
	}

	s.server = &http.Server{
		Addr:              s.addr,
		Handler:           m.Logging(mux),
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
	}

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
