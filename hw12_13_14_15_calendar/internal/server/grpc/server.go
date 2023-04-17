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
	res := &gen.CreateEventResult{}
	ev := req.GetData()

	if ev == nil {
		return res, errors.New("пустой запрос")
	}

	if ev.Title == "" {
		return res, errors.New("пустой заголовок")
	}

	if !ev.StartTime.IsValid() {
		return res, errors.New("неверный формат даты старта")
	}

	if !ev.EndTime.IsValid() {
		return res, errors.New("неверный формат даты окончания")
	}

	if ev.OwnUserID == 0 {
		return res, errors.New("неверный пользователь")
	}

	data := storage.Event{
		Title:       ev.Title,
		Description: ev.Description,
		StartTime:   ev.StartTime.AsTime(),
		EndTime:     ev.EndTime.AsTime(),
		OwnUserID:   ev.OwnUserID,
	}

	if data.StartTime.Before(time.Now()) {
		return res, errors.New("неверная дата старта")
	}

	if data.EndTime.Before(data.StartTime) {
		return res, errors.New("неверная дата окончания")
	}

	id, err := s.worker.CreateEvent(ctx, &data)

	res.ID = id
	if err != nil {
		res.Err = &gen.Error{Description: err.Error()}
	}

	return res, nil
}

func (s *Server) UpdateEvent(ctx context.Context, req *gen.UpdateEventRequest) (*gen.UpdateEventResult, error) {
	res := &gen.UpdateEventResult{}
	ev := req.GetData()

	if ev == nil {
		return res, errors.New("пустой запрос")
	}

	if ev.ID == 0 {
		return res, errors.New("неверный ID события")
	}

	if ev.Title == "" {
		return res, errors.New("пустой заголовок")
	}

	if !ev.StartTime.IsValid() {
		return res, errors.New("неверный формат даты старта")
	}

	if !ev.EndTime.IsValid() {
		return res, errors.New("неверный формат даты окончания")
	}

	if ev.OwnUserID == 0 {
		return res, errors.New("неверный пользователь")
	}

	data := storage.Event{
		ID:          ev.ID,
		Title:       ev.Title,
		Description: ev.Description,
		StartTime:   ev.StartTime.AsTime(),
		EndTime:     ev.EndTime.AsTime(),
		OwnUserID:   ev.OwnUserID,
	}

	if data.StartTime.Before(time.Now()) {
		return res, errors.New("неверная дата старта")
	}

	if data.EndTime.Before(data.StartTime) {
		return res, errors.New("неверная дата окончания")
	}

	err := s.worker.UpdateEvent(ctx, &data)
	if err != nil {
		res.Err = &gen.Error{Description: err.Error()}
	}

	return res, nil
}

func (s *Server) DeleteEvent(ctx context.Context, req *gen.DeleteEventRequest) (*gen.DeleteEventResult, error) {
	res := &gen.DeleteEventResult{}

	if req == nil {
		return res, errors.New("пустой запрос")
	}

	if req.ID == 0 {
		return res, errors.New("неверный ID события")
	}

	err := s.worker.DeleteEvent(ctx, req.ID)
	if err != nil {
		res.Err = &gen.Error{Description: err.Error()}
	}
	return res, nil
}

func (s *Server) ListEventOnDate(
	ctx context.Context,
	req *gen.ListEventOnDateRequest,
) (*gen.ListEventOnDateResult, error) {
	res := &gen.ListEventOnDateResult{}
	if req == nil {
		return res, errors.New("пустой запрос")
	}

	data, err := s.worker.ListEventOnDate(ctx, int(req.Year), int(req.Month), int(req.Day))

	res.Data = s.convertListData(data)
	if err != nil {
		res.Err = &gen.Error{Description: err.Error()}
	}
	return res, nil
}

func (s *Server) ListEventOnWeek(
	ctx context.Context,
	req *gen.ListEventOnWeekRequest,
) (*gen.ListEventOnWeekResult, error) {
	res := &gen.ListEventOnWeekResult{}
	if req == nil {
		return res, errors.New("пустой запрос")
	}
	data, err := s.worker.ListEventOnWeek(ctx, int(req.Year), int(req.Week))
	res.Data = s.convertListData(data)
	if err != nil {
		res.Err = &gen.Error{Description: err.Error()}
	}

	return res, nil
}

func (s *Server) ListEventOnMonth(
	ctx context.Context,
	req *gen.ListEventOnMonthRequest,
) (*gen.ListEventOnMonthResult, error) {
	res := &gen.ListEventOnMonthResult{}
	if req == nil {
		return res, errors.New("пустой запрос")
	}
	data, err := s.worker.ListEventOnMonth(ctx, int(req.Year), int(req.Month))
	res.Data = s.convertListData(data)
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
