package internalgrpc

import (
	"context"
	"fmt"
	"github.com/dimonk33/kdv_hw_otus/hw12_13_14_15_calendar/api/gen"
	"github.com/dimonk33/kdv_hw_otus/hw12_13_14_15_calendar/internal/storage"
	"google.golang.org/grpc"
	"net"
	"net/http"
)

type Server struct {
	gen.EventsServer
	addr       string
	logger     Logger
	server     *http.Server
	grpcServer *grpc.Server
}

type Logger interface {
	Info(msg string)
	Error(msg string)
	Warning(msg string)
}

type Application interface {
	CreateEvent(ctx context.Context, e storage.Event) (int64, error)
	UpdateEvent(ctx context.Context, e storage.Event) error
	DeleteEvent(ctx context.Context, ID int64) error
	ListEventOnDate(ctx context.Context, year int, month int, day int) ([]storage.Event, error)
	ListEventOnWeek(ctx context.Context, year int, week int) ([]storage.Event, error)
	ListEventOnMonth(ctx context.Context, year int, month int) ([]storage.Event, error)
}

func NewServer(addr string, logger Logger, app Application) *Server {
	s := grpc.NewServer()
	myInvoicerServer := &Server{
		logger:     logger,
		addr:       addr,
		grpcServer: s,
	}
	gen.RegisterEventsServer(s, myInvoicerServer)

	return myInvoicerServer
}

func (s *Server) Start(ctx context.Context) error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		s.logger.Error(fmt.Sprintf("failed to listen: %v", err))
		return err
	}
	s.logger.Info(fmt.Sprintf("gRPC сервер запущен на порту %v", lis.Addr()))
	if err = s.grpcServer.Serve(lis); err != nil {
		s.logger.Error(fmt.Sprintf("failed to serve: %v", err))
		return err
	}

	<-ctx.Done()

	return ctx.Err()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *Server) CreateEvent(context.Context, *gen.CreateEventRequest) (*gen.CreateEventResult, error) {
	res := gen.CreateEventResult{}

	return &res, nil
}

func (s *Server) UpdateEvent(context.Context, *gen.UpdateEventRequest) (*gen.UpdateEventResult, error) {
	res := gen.UpdateEventResult{}

	return &res, nil
}

func (s *Server) DeleteEvent(context.Context, *gen.DeleteEventRequest) (*gen.DeleteEventResult, error) {
	res := gen.DeleteEventResult{}

	return &res, nil
}

func (s *Server) ListEventOnDate(context.Context, *gen.ListEventOnDateRequest) (*gen.ListEventOnDateResult, error) {
	res := gen.ListEventOnDateResult{}

	return &res, nil
}

func (s *Server) ListEventOnWeek(context.Context, *gen.ListEventOnWeekRequest) (*gen.ListEventOnWeekResult, error) {
	res := gen.ListEventOnWeekResult{}

	return &res, nil
}

func (s *Server) ListEventOnMonth(context.Context, *gen.ListEventOnMonthRequest) (*gen.ListEventOnMonthResult, error) {
	res := gen.ListEventOnMonthResult{}

	return &res, nil
}
