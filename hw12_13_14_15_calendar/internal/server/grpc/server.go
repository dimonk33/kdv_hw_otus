package internalgrpc

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

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
	m := Middleware{Logger: logger}
	s := grpc.NewServer(grpc.UnaryInterceptor(m.unaryInterceptor))
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

	if !ev.GetStartTime().IsValid() {
		return nil, errors.New("неверный формат даты старта")
	}

	if !ev.GetEndTime().IsValid() {
		return nil, errors.New("неверный формат даты окончания")
	}

	data := storage.Event{
		Title:       ev.GetTitle(),
		Description: ev.GetDescription(),
		StartTime:   ev.GetStartTime().AsTime(),
		EndTime:     ev.GetEndTime().AsTime(),
		OwnUserID:   ev.GetOwnUserID(),
	}

	if data.Title == "" {
		return nil, errors.New("пустой заголовок")
	}

	if data.OwnUserID <= 0 {
		return nil, errors.New("неверный пользователь")
	}

	if data.StartTime.Before(time.Now()) {
		return nil, errors.New("неверная дата старта")
	}

	if data.EndTime.Before(data.StartTime) {
		return nil, errors.New("неверная дата окончания")
	}

	id, err := s.worker.CreateEvent(ctx, &data)

	res := &gen.CreateEventResult{ID: id}
	if err != nil {
		res.Err = &gen.Error{Description: err.Error()}
	}

	return res, nil
}

func (s *Server) UpdateEvent(ctx context.Context, req *gen.UpdateEventRequest) (*gen.UpdateEventResult, error) {
	ev := req.GetData()

	if !ev.StartTime.IsValid() {
		return nil, errors.New("неверный формат даты старта")
	}

	if !ev.EndTime.IsValid() {
		return nil, errors.New("неверный формат даты окончания")
	}

	data := storage.Event{
		ID:          ev.GetID(),
		Title:       ev.GetTitle(),
		Description: ev.GetDescription(),
		StartTime:   ev.GetStartTime().AsTime(),
		EndTime:     ev.GetEndTime().AsTime(),
		OwnUserID:   ev.GetOwnUserID(),
	}

	if data.OwnUserID <= 0 {
		return nil, errors.New("неверный ID события")
	}

	if data.Title == "" {
		return nil, errors.New("пустой заголовок")
	}

	if data.OwnUserID <= 0 {
		return nil, errors.New("неверный пользователь")
	}

	if data.StartTime.Before(time.Now()) {
		return nil, errors.New("неверная дата старта")
	}

	if data.EndTime.Before(data.StartTime) {
		return nil, errors.New("неверная дата окончания")
	}

	err := s.worker.UpdateEvent(ctx, &data)

	var res *gen.UpdateEventResult
	if err != nil {
		res = &gen.UpdateEventResult{Err: &gen.Error{Description: err.Error()}}
	}

	return res, nil
}

func (s *Server) DeleteEvent(ctx context.Context, req *gen.DeleteEventRequest) (*gen.DeleteEventResult, error) {
	if req.GetID() <= 0 {
		return nil, errors.New("неверный ID события")
	}

	err := s.worker.DeleteEvent(ctx, req.ID)

	var res *gen.DeleteEventResult
	if err != nil {
		res = &gen.DeleteEventResult{Err: &gen.Error{Description: err.Error()}}
	}
	return res, nil
}

func (s *Server) ListEventOnDate(
	ctx context.Context,
	req *gen.ListEventOnDateRequest,
) (*gen.ListEventOnDateResult, error) {
	y := int(req.GetYear())
	m := int(req.GetMonth())
	d := int(req.GetDay())
	if y <= 0 || m <= 0 || d <= 0 {
		return nil, errors.New("неверная дата")
	}
	data, err := s.worker.ListEventOnDate(ctx, y, m, d)

	res := &gen.ListEventOnDateResult{Data: s.convertListData(data)}
	if err != nil {
		res.Err = &gen.Error{Description: err.Error()}
	}
	return res, nil
}

func (s *Server) ListEventOnWeek(
	ctx context.Context,
	req *gen.ListEventOnWeekRequest,
) (*gen.ListEventOnWeekResult, error) {
	y := int(req.GetYear())
	w := int(req.GetWeek())
	if y <= 0 || w <= 0 {
		return nil, errors.New("неверная дата")
	}
	data, err := s.worker.ListEventOnWeek(ctx, y, w)

	res := &gen.ListEventOnWeekResult{Data: s.convertListData(data)}
	if err != nil {
		res.Err = &gen.Error{Description: err.Error()}
	}

	return res, nil
}

func (s *Server) ListEventOnMonth(
	ctx context.Context,
	req *gen.ListEventOnMonthRequest,
) (*gen.ListEventOnMonthResult, error) {
	y := int(req.GetYear())
	m := int(req.GetMonth())
	if y <= 0 || m <= 0 {
		return nil, errors.New("неверная дата")
	}
	data, err := s.worker.ListEventOnMonth(ctx, y, m)

	res := &gen.ListEventOnMonthResult{Data: s.convertListData(data)}
	if err != nil {
		res.Err = &gen.Error{Description: err.Error()}
	}

	return res, nil
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
