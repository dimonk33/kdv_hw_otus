package scheduler

import (
	"context"
	"strconv"
	"time"

	"github.com/dimonk33/kdv_hw_otus/hw12_13_14_15_calendar/internal/storage"
)

type Scheduler struct {
	notifyTime NotifyTime
	storage    Storage
	sender     Sender
	logger     Logger
}

type Logger interface {
	Info(msg string)
	Error(msg string)
	Warning(msg string)
	Debug(msg string)
}

type Storage interface {
	Delete(ctx context.Context, id int64) error
	ListOnDate(ctx context.Context, year int, month int, day int) ([]storage.Event, error)
	ListLessDate(ctx context.Context, year, month, day int) ([]storage.Event, error)
}

type Sender interface {
	Send(ctx context.Context, data interface{}) error
}

type Notify struct {
	ID    int64     `json:"id"`
	Title string    `json:"title"`
	Date  time.Time `json:"date"`
	User  string    `json:"user"`
}

type NotifyTime struct {
	H int
	M int
}

func NewScheduler(_storage Storage, _sender Sender, _notifyTime NotifyTime, _logger Logger) *Scheduler {
	s := &Scheduler{
		storage:    _storage,
		sender:     _sender,
		notifyTime: _notifyTime,
		logger:     _logger,
	}

	return s
}

func (s *Scheduler) Start(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(time.Minute)
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				h, m, _ := time.Now().Clock()
				if m == s.notifyTime.M && h == s.notifyTime.H {
					s.logger.Info("отправка событий на день")
					err := s.sendEvents(ctx)
					if err != nil {
						s.logger.Error("ошибка при отправке уведомления: " + err.Error())
					}
					err = s.clearEvents(ctx)
					if err != nil {
						s.logger.Error("ошибка при удалении уведомления: " + err.Error())
					}
				}
			}
		}
	}()
}

func (s *Scheduler) sendEvents(ctx context.Context) error {
	eventDate := time.Now()
	y, m, d := eventDate.Date()
	items, err := s.storage.ListOnDate(ctx, y, int(m), d)
	if err != nil {
		return err
	}
	for _, item := range items {
		err = s.sendEventToQueue(ctx, item)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Scheduler) clearEvents(ctx context.Context) error {
	clearDate := time.Now().AddDate(-1, 0, 0)
	y, m, d := clearDate.Date()
	items, err := s.storage.ListLessDate(ctx, y, int(m), d)
	if err != nil {
		return err
	}
	for _, item := range items {
		err = s.storage.Delete(ctx, item.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Scheduler) sendEventToQueue(ctx context.Context, event storage.Event) error {
	notify := Notify{
		ID:    event.ID,
		Title: event.Title,
		Date:  event.StartTime,
		User:  strconv.Itoa(int(event.OwnUserID)),
	}
	return s.sender.Send(ctx, notify)
}
