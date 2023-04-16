package internalhttp

import (
	"context"
	"net/http"

	gs "github.com/dimonk33/kdv_hw_otus/hw12_13_14_15_calendar/api/gen"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

type Server struct {
	addr       string
	logger     Logger
	grpcServer gs.EventsServer
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
		Logger: s.logger,
	}

	if err := http.ListenAndServe(s.addr, m.Logging(mux)); err != nil {
		return err
	}
	s.logger.Info("Http сервер запущен")

	<-ctx.Done()

	return ctx.Err()
}
