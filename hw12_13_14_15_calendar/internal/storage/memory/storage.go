package memorystorage

import (
	"context"
	"fmt"
	"github.com/dimonk33/kdv_hw_otus/hw12_13_14_15_calendar/internal/storage"
	"sync"
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

func (s *Storage) Create(data storage.Event) (int64, error) {
	s.mu.Lock()
	curId := s.id
	s.id++
	s.mu.Unlock()
	data.ID = curId
	s.db[curId] = data

	return curId, nil
}

func (s *Storage) Update(data storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.db[data.ID]
	if !ok {
		return fmt.Errorf("отсутствует запись с id = %d", data.ID)
	}
	s.db[data.ID] = data
	return nil
}

func (s *Storage) Delete(id int64) error {
	delete(s.db, id)
	return nil
}

func (s *Storage) ListOnDate(ctx context.Context, year int, month int, day int) ([]storage.Event, error) {
	var out []storage.Event
	for _, item := range s.db {
		select {
		case <-ctx.Done():
			return out, ctx.Err()
		default:
			yS, mS, dS := item.StartTime.Date()
			yE, mE, dE := item.EndTime.Date()
			if (yS == year && int(mS) == month && dS == day) || (yE == year && int(mE) == month && dE == day) {
				out = append(out, item)
			}
		}
	}

	return out, nil
}

func (s *Storage) ListOnWeek(ctx context.Context, year int, week int) ([]storage.Event, error) {
	var out []storage.Event
	for _, item := range s.db {
		select {
		case <-ctx.Done():
			return out, ctx.Err()
		default:
			yS, wS := item.StartTime.ISOWeek()
			yE, wE := item.EndTime.ISOWeek()
			if (yS == year && wS == week) || (yE == year && wE == week) {
				out = append(out, item)
			}
		}
	}

	return out, nil
}

func (s *Storage) ListOnMonth(ctx context.Context, year int, month int) ([]storage.Event, error) {
	var out []storage.Event
	for _, item := range s.db {
		select {
		case <-ctx.Done():
			return out, ctx.Err()
		default:
			yS, mS, _ := item.StartTime.Date()
			yE, mE, _ := item.EndTime.Date()
			if (yS == year && int(mS) == month) || (yE == year && int(mE) == month) {
				out = append(out, item)
			}
		}
	}

	return out, nil
}
