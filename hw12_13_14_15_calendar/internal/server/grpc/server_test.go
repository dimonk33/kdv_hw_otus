package internalgrpc

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/dimonk33/kdv_hw_otus/hw12_13_14_15_calendar/api/gen"
	"github.com/dimonk33/kdv_hw_otus/hw12_13_14_15_calendar/internal/app"
	"github.com/dimonk33/kdv_hw_otus/hw12_13_14_15_calendar/internal/logger"
	memorystorage "github.com/dimonk33/kdv_hw_otus/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestGrpcHandler(t *testing.T) {
	testCases := []struct {
		name        string
		req         *gen.CreateEventRequest
		ID          int64
		expectedErr bool
	}{
		{
			name: "req ok",
			req: &gen.CreateEventRequest{
				Data: &gen.Event{
					ID:          0,
					Title:       "тест",
					Description: "тестовое событие",
					OwnUserID:   1,
					StartTime:   timestamppb.New(time.Now().Add(1 * time.Hour)),
					EndTime:     timestamppb.New(time.Now().Add(2 * time.Hour)),
				},
			},
			ID:          1,
			expectedErr: false,
		},
		{
			name:        "req with empty event",
			req:         &gen.CreateEventRequest{},
			expectedErr: true,
		},
		{
			name:        "nil request",
			req:         nil,
			expectedErr: true,
		},
	}
	for _, tc := range testCases {
		testCase := tc
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			// call
			l := logger.Logger{}
			s := memorystorage.New()
			a := app.New(l, s)
			srv := NewServer(":5555", l, a)
			response, err := srv.CreateEvent(ctx, testCase.req)

			// assert results expectations
			if testCase.expectedErr {
				require.Error(t, err)
			} else {
				require.Equal(t, int64(1), response.ID)
			}
		})
	}

	t.Run("list events", func(t *testing.T) {
		const numEvents = 20

		ctx := context.Background()
		l := logger.Logger{}
		s := memorystorage.New()
		a := app.New(l, s)
		srv := NewServer(":5555", l, a)

		initDate := time.Now()
		initDate = time.Date(initDate.Year(), initDate.Month()+1, 1, 9, 0, 0, 0, time.Local)

		startDate := initDate
		for i := 0; i < numEvents; i++ {
			endDate := startDate.Add(1 * time.Hour)
			req := gen.CreateEventRequest{
				Data: &gen.Event{
					Title:     "Тестовое событие " + strconv.Itoa(i+1),
					StartTime: timestamppb.New(startDate),
					EndTime:   timestamppb.New(endDate),
					OwnUserID: 1,
				},
			}
			resp, err := srv.CreateEvent(ctx, &req)
			require.NoError(t, err)
			require.Greater(t, resp.GetID(), int64(0))

			startDate = startDate.AddDate(0, 0, 1)
		}

		req := gen.ListEventOnDateRequest{
			Year:  int32(initDate.Year()),
			Month: int32(initDate.Month()),
			Day:   int32(initDate.Day() + 1),
		}
		resp, err := srv.ListEventOnDate(ctx, &req)
		require.NoError(t, err)
		list := resp.GetData()
		require.Equal(t, len(list), 1)

		y, w := initDate.ISOWeek()
		reqW := gen.ListEventOnWeekRequest{
			Year: int32(y),
			Week: int32(w + 1),
		}
		respW, errW := srv.ListEventOnWeek(ctx, &reqW)
		require.NoError(t, errW)
		list = respW.GetData()
		require.Equal(t, len(list), 7)

		reqM := gen.ListEventOnMonthRequest{
			Year:  int32(initDate.Year()),
			Month: int32(initDate.Month()),
		}
		respM, errM := srv.ListEventOnMonth(ctx, &reqM)
		require.NoError(t, errM)
		list = respM.GetData()
		require.Equal(t, len(list), numEvents)

		for _, ev := range list {
			reqDel := gen.DeleteEventRequest{
				ID: ev.GetID(),
			}
			respDel, errDel := srv.DeleteEvent(ctx, &reqDel)
			require.NoError(t, errDel)
			require.Nil(t, respDel.GetErr())
		}
	})
}
