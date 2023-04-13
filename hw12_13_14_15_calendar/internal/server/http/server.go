package internalhttp

import (
	"context"
	"github.com/dimonk33/kdv_hw_otus/hw12_13_14_15_calendar/internal/storage"
	"io"
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
	CreateEvent(ctx context.Context, e storage.Event) (int64, error)
	UpdateEvent(ctx context.Context, e storage.Event) error
	DeleteEvent(ctx context.Context, ID int64) error
	ListEventOnDate(ctx context.Context, year int, month int, day int) ([]storage.Event, error)
	ListEventOnWeek(ctx context.Context, year int, week int) ([]storage.Event, error)
	ListEventOnMonth(ctx context.Context, year int, month int) ([]storage.Event, error)
}

func NewServer(addr string, logger Logger, app Application) *Server {
	s := &Server{
		logger: logger,
	}

	router := http.NewServeMux()
	router.HandleFunc("/hello", s.HelloHandler)
	router.HandleFunc("/create-event", s.CreateEventHandler)
	router.HandleFunc("/update-event", s.UpdateEventHandler)
	router.HandleFunc("/delete-event", s.DeleteEventHandler)
	router.HandleFunc("/list-event-on-date", s.ListOnDateHandler)
	router.HandleFunc("/list-event-on-week", s.ListOnWeekHandler)
	router.HandleFunc("/list-event-on-month", s.ListOnMonthHandler)

	m := Middleware{
		Logger: logger,
	}

	s.server = &http.Server{
		Addr:              addr,
		Handler:           m.Logging(router),
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

func (s *Server) HelloHandler(w http.ResponseWriter, r *http.Request) {
	writeBytes, err := io.WriteString(w, "Hello, HTTP!\n")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("500 - " + err.Error()))
	}
	if writeBytes == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("500 - Данные не отправлены"))
	}
}

func (s *Server) CreateEventHandler(w http.ResponseWriter, r *http.Request) {
	writeBytes, err := io.WriteString(w, "Hello, HTTP!\n")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("500 - " + err.Error()))
	}
	if writeBytes == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("500 - Данные не отправлены"))
	}
}

func (s *Server) UpdateEventHandler(w http.ResponseWriter, r *http.Request) {
	writeBytes, err := io.WriteString(w, "Hello, HTTP!\n")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("500 - " + err.Error()))
	}
	if writeBytes == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("500 - Данные не отправлены"))
	}
}

func (s *Server) DeleteEventHandler(w http.ResponseWriter, r *http.Request) {
	writeBytes, err := io.WriteString(w, "Hello, HTTP!\n")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("500 - " + err.Error()))
	}
	if writeBytes == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("500 - Данные не отправлены"))
	}
}

func (s *Server) ListOnDateHandler(w http.ResponseWriter, r *http.Request) {
	writeBytes, err := io.WriteString(w, "Hello, HTTP!\n")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("500 - " + err.Error()))
	}
	if writeBytes == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("500 - Данные не отправлены"))
	}
}

func (s *Server) ListOnWeekHandler(w http.ResponseWriter, r *http.Request) {
	writeBytes, err := io.WriteString(w, "Hello, HTTP!\n")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("500 - " + err.Error()))
	}
	if writeBytes == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("500 - Данные не отправлены"))
	}
}

func (s *Server) ListOnMonthHandler(w http.ResponseWriter, r *http.Request) {
	writeBytes, err := io.WriteString(w, "Hello, HTTP!\n")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("500 - " + err.Error()))
	}
	if writeBytes == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("500 - Данные не отправлены"))
	}
}
