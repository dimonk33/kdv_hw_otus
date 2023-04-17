package internalgrpc

import (
	"context"
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
				require.NotNil(t, err)
			} else {
				require.Equal(t, int64(1), response.ID)
			}
		})
	}
}
