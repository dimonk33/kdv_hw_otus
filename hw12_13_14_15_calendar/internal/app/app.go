package app

import (
	"context"
	"io"
	"net/http"

	"github.com/dimonk33/kdv_hw_otus/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	logger  Logger
	storage Storage
}

type Logger interface {
	Info(msg string)
	Error(msg string)
	Warning(msg string)
	Debug(msg string)
}

type Storage interface {
	Create(data storage.Event) (int64, error)
	Update(data storage.Event) error
	Delete(id int64) error
	ListOnDate(ctx context.Context, year int, month int, day int) ([]storage.Event, error)
	ListOnWeek(ctx context.Context, year int, week int) ([]storage.Event, error)
	ListOnMonth(ctx context.Context, year int, month int) ([]storage.Event, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	// TODO
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

func (a *App) HelloHandler(w http.ResponseWriter, r *http.Request) {
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
