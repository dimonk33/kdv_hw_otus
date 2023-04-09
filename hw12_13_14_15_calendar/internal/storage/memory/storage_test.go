package memorystorage

import (
	"context"
	"testing"
	"time"

	"github.com/dimonk33/kdv_hw_otus/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	s := New()

	t.Run("main logic", func(t *testing.T) {
		evStartTime := time.Now()
		evEndTime := evStartTime.Add(1 * time.Hour)
		event := storage.Event{
			ID:          5,
			Title:       "Тест1",
			StartTime:   evStartTime,
			EndTime:     evEndTime,
			Description: "Тестовое событие",
			OwnUserID:   1,
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		id, err := s.Create(ctx, event)
		require.Nil(t, err)
		require.Equal(t, id, int64(0))

		event.ID = id
		event.Title = "Тест2"
		err = s.Update(ctx, event)
		require.Nil(t, err)

		y, m, d := evStartTime.Date()
		events, errList := s.ListOnDate(context.Background(), y, int(m), d)
		require.Nil(t, errList)
		require.Equal(t, len(events), 1)
		require.Equal(t, events[0], event)

		err = s.Delete(ctx, id)
		require.Nil(t, err)

		cancel()
	})
}
