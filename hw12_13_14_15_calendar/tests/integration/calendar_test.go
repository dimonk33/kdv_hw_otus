//go:build integration

package integration_test

import (
	"context"
	"github.com/dimonk33/kdv_hw_otus/hw12_13_14_15_calendar/api/gen"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
	"os"
	"strconv"
	"testing"
	"time"
)

type CalendarSuite struct {
	suite.Suite
	ctx       context.Context
	apiConn   *grpc.ClientConn
	apiClient gen.EventsClient
}

func TestCalendarSuite(t *testing.T) {
	suite.Run(t, new(CalendarSuite))
}

func (s *CalendarSuite) SetupSuite() {
	apiHost := os.Getenv("API_SERVER_HOST")
	if apiHost == "" {
		apiHost = "localhost:9190"
	}
	apiConn, err := grpc.Dial(apiHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)

	s.apiClient = gen.NewEventsClient(apiConn)
	s.ctx = context.Background()
}

func (s *CalendarSuite) SetupTest() {

}

func (s *CalendarSuite) TestCreateEventSuccess() {
	startDate := time.Now().Add(1 * time.Hour)
	endDate := startDate.Add(1 * time.Hour)
	reqCreate := gen.CreateEventRequest{
		Data: &gen.Event{
			Title:     "Тестовое событие",
			StartTime: timestamppb.New(startDate),
			EndTime:   timestamppb.New(endDate),
			OwnUserID: 1,
		},
	}
	resp, err := s.apiClient.CreateEvent(s.ctx, &reqCreate)
	s.Require().NoError(err)
	s.Require().Greater(resp.GetID(), int64(0))

	reqDel := gen.DeleteEventRequest{
		ID: resp.GetID(),
	}
	respDel, errDel := s.apiClient.DeleteEvent(s.ctx, &reqDel)
	s.Require().NoError(errDel)
	s.Require().Nil(respDel.GetErr())
}

func (s *CalendarSuite) TestCreateEventFail() {
	startDate := time.Now().Add(1 * time.Hour)
	endDate := startDate.Add(-1 * time.Hour)
	req := gen.CreateEventRequest{
		Data: &gen.Event{
			Title:     "Тестовое событие",
			StartTime: timestamppb.New(startDate),
			EndTime:   timestamppb.New(endDate),
			OwnUserID: 1,
		},
	}
	resp, err := s.apiClient.CreateEvent(s.ctx, &req)
	s.Require().Error(err)
	s.Require().Equal(resp.GetID(), int64(0))

	req = gen.CreateEventRequest{
		Data: &gen.Event{
			Title:     "Тестовое событие",
			StartTime: timestamppb.New(startDate),
			EndTime:   timestamppb.New(endDate),
		},
	}
	resp, err = s.apiClient.CreateEvent(s.ctx, &req)
	s.Require().Error(err)
	s.Require().Equal(resp.GetID(), int64(0))

}

func (s *CalendarSuite) TestGetEvents() {
	const numEvents = 20

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
		resp, err := s.apiClient.CreateEvent(s.ctx, &req)
		s.Require().NoError(err)
		s.Require().Greater(resp.GetID(), int64(0))

		startDate = startDate.AddDate(0, 0, 1)
	}

	req := gen.ListEventOnDateRequest{
		Year:  int32(initDate.Year()),
		Month: int32(initDate.Month()),
		Day:   int32(initDate.Day() + 1),
	}
	resp, err := s.apiClient.ListEventOnDate(s.ctx, &req)
	s.Require().NoError(err)
	list := resp.GetData()
	s.Require().Equal(len(list), 1)

	y, w := initDate.ISOWeek()
	reqW := gen.ListEventOnWeekRequest{
		Year: int32(y),
		Week: int32(w + 1),
	}
	respW, errW := s.apiClient.ListEventOnWeek(s.ctx, &reqW)
	s.Require().NoError(errW)
	list = respW.GetData()
	s.Require().Equal(len(list), 7)

	reqM := gen.ListEventOnMonthRequest{
		Year:  int32(initDate.Year()),
		Month: int32(initDate.Month()),
	}
	respM, errM := s.apiClient.ListEventOnMonth(s.ctx, &reqM)
	s.Require().NoError(errM)
	list = respM.GetData()
	s.Require().Equal(len(list), numEvents)

	for _, ev := range list {
		reqDel := gen.DeleteEventRequest{
			ID: ev.GetID(),
		}
		respDel, errDel := s.apiClient.DeleteEvent(s.ctx, &reqDel)
		s.Require().NoError(errDel)
		s.Require().Nil(respDel.GetErr())
	}
}
