package internalgrpc

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/dimonk33/kdv_hw_otus/hw12_13_14_15_calendar/api/gen"
	"github.com/dimonk33/kdv_hw_otus/hw12_13_14_15_calendar/internal/storage"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	gen.EventsServer
	addr       string
	logger     Logger
	worker     Application
	server     *http.Server
	grpcServer *grpc.Server
}

type Logger interface {
	Info(msg string)
	Error(msg string)
	Warning(msg string)
}

type Application interface {
	CreateEvent(ctx context.Context, e *storage.Event) (int64, error)
	UpdateEvent(ctx context.Context, e *storage.Event) error
	DeleteEvent(ctx context.Context, id int64) error
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
		worker:     app,
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

func (s *Server) CreateEvent(ctx context.Context, req *gen.CreateEventRequest) (*gen.CreateEventResult, error) {
	ev := req.GetData()
	data := storage.Event{
		Title:       ev.Title,
		Description: ev.Description,
		StartTime:   ev.StartTime.AsTime(),
		EndTime:     ev.EndTime.AsTime(),
		OwnUserID:   ev.OwnUserID,
	}
	id, err := s.worker.CreateEvent(ctx, &data)

	res := gen.CreateEventResult{ID: id, Err: &gen.Error{Description: err.Error()}}

	return &res, nil
}

func (s *Server) UpdateEvent(ctx context.Context, req *gen.UpdateEventRequest) (*gen.UpdateEventResult, error) {
	ev := req.GetData()
	data := storage.Event{
		ID:          ev.ID,
		Title:       ev.Title,
		Description: ev.Description,
		StartTime:   ev.StartTime.AsTime(),
		EndTime:     ev.EndTime.AsTime(),
		OwnUserID:   ev.OwnUserID,
	}
	err := s.worker.UpdateEvent(ctx, &data)

	res := gen.UpdateEventResult{Err: &gen.Error{Description: err.Error()}}
	return &res, nil
}

func (s *Server) DeleteEvent(ctx context.Context, req *gen.DeleteEventRequest) (*gen.DeleteEventResult, error) {
	err := s.worker.DeleteEvent(ctx, req.ID)
	res := gen.DeleteEventResult{Err: &gen.Error{Description: err.Error()}}

	return &res, nil
}

func (s *Server) ListEventOnDate(
	ctx context.Context,
	req *gen.ListEventOnDateRequest,
) (*gen.ListEventOnDateResult, error) {
	data, err := s.worker.ListEventOnDate(ctx, int(req.Year), int(req.Month), int(req.Day))
	res := gen.ListEventOnDateResult{Data: s.convertListData(data), Err: &gen.Error{Description: err.Error()}}

	return &res, nil
}

func (s *Server) ListEventOnWeek(
	ctx context.Context,
	req *gen.ListEventOnWeekRequest,
) (*gen.ListEventOnWeekResult, error) {
	data, err := s.worker.ListEventOnWeek(ctx, int(req.Year), int(req.Week))
	res := gen.ListEventOnWeekResult{Data: s.convertListData(data), Err: &gen.Error{Description: err.Error()}}

	return &res, nil
}

func (s *Server) ListEventOnMonth(
	ctx context.Context,
	req *gen.ListEventOnMonthRequest,
) (*gen.ListEventOnMonthResult, error) {
	data, err := s.worker.ListEventOnMonth(ctx, int(req.Year), int(req.Month))
	res := gen.ListEventOnMonthResult{Data: s.convertListData(data), Err: &gen.Error{Description: err.Error()}}

	return &res, nil
}

func (s *Server) convertListData(inData []storage.Event) []*gen.Event {
	outData := make([]*gen.Event, len(inData))
	for i, item := range inData {
		outData[i] = &gen.Event{
			ID:          item.ID,
			Title:       item.Title,
			Description: item.Description,
			StartTime:   timestamppb.New(item.StartTime),
			EndTime:     timestamppb.New(item.EndTime),
		}
	}
	return outData
}
