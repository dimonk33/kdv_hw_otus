package memorystorage

import (
	"context"
	"fmt"
	"sync"

	"github.com/dimonk33/kdv_hw_otus/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	db map[int64]storage.Event
	id int64
	mu sync.RWMutex
}

func New() *Storage {
	db := make(map[int64]storage.Event, 1)
	return &Storage{db: db}
}

type ValidateDate func(item storage.Event) bool

func (s *Storage) Create(ctx context.Context, data *storage.Event) (int64, error) {
	s.mu.Lock()
	curID := s.id
	s.id++
	s.mu.Unlock()
	data.ID = curID
	s.db[curID] = *data

	return curID, nil
}

func (s *Storage) Update(ctx context.Context, data *storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.db[data.ID]; !ok {
		return fmt.Errorf("отсутствует запись с id = %d", data.ID)
	}
	s.db[data.ID] = *data
	return nil
}

func (s *Storage) Delete(ctx context.Context, id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.db, id)
	return nil
}

func (s *Storage) ListOnDate(ctx context.Context, year, month, day int) ([]storage.Event, error) {
	return s.listItems(ctx, func(item storage.Event) bool {
		yS, mS, dS := item.StartTime.Date()
		yE, mE, dE := item.EndTime.Date()
		return (yS == year && int(mS) == month && dS == day) || (yE == year && int(mE) == month && dE == day)
	})
}

func (s *Storage) ListOnWeek(ctx context.Context, year, week int) ([]storage.Event, error) {
	return s.listItems(ctx, func(item storage.Event) bool {
		yS, wS := item.StartTime.ISOWeek()
		yE, wE := item.EndTime.ISOWeek()
		return (yS == year && wS == week) || (yE == year && wE == week)
	})
}

func (s *Storage) ListOnMonth(ctx context.Context, year, month int) ([]storage.Event, error) {
	return s.listItems(ctx, func(item storage.Event) bool {
		yS, mS, _ := item.StartTime.Date()
		yE, mE, _ := item.EndTime.Date()
		return (yS == year && int(mS) == month) || (yE == year && int(mE) == month)
	})
}

func (s *Storage) listItems(ctx context.Context, validate ValidateDate) ([]storage.Event, error) {
	var out []storage.Event
	var i int64
	var ok bool
	var item storage.Event

	for {
		s.mu.RLock()
		item, ok = s.db[i]
		s.mu.RUnlock()
		if !ok {
			break
		}
		select {
		case <-ctx.Done():
			return out, ctx.Err()
		default:
			if validate(item) {
				out = append(out, item)
			}
		}
		i++
	}

	return out, nil
}
