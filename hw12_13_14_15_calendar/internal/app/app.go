package app

import (
	"context"

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
	Create(ctx context.Context, data *storage.Event) (int64, error)
	Update(ctx context.Context, data *storage.Event) error
	Delete(ctx context.Context, id int64) error
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

func (a *App) CreateEvent(ctx context.Context, e *storage.Event) (int64, error) {
	return a.storage.Create(ctx, e)
}

func (a *App) UpdateEvent(ctx context.Context, e *storage.Event) error {
	return a.storage.Update(ctx, e)
}

func (a *App) DeleteEvent(ctx context.Context, id int64) error {
	return a.storage.Delete(ctx, id)
}

func (a *App) ListEventOnDate(ctx context.Context, year int, month int, day int) ([]storage.Event, error) {
	return a.storage.ListOnDate(ctx, year, month, day)
}

func (a *App) ListEventOnWeek(ctx context.Context, year int, week int) ([]storage.Event, error) {
	return a.storage.ListOnWeek(ctx, year, week)
}

func (a *App) ListEventOnMonth(ctx context.Context, year int, month int) ([]storage.Event, error) {
	return a.storage.ListOnMonth(ctx, year, month)
}
