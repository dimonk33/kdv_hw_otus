//go:build integration

package integration_test

import (
	"context"
	"github.com/dimonk33/kdv_hw_otus/hw12_13_14_15_calendar/internal/logger"
	kafkaapp "github.com/dimonk33/kdv_hw_otus/hw12_13_14_15_calendar/internal/queue/kafka"
	"github.com/dimonk33/kdv_hw_otus/hw12_13_14_15_calendar/internal/scheduler"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
	"time"
)

type SenderSuite struct {
	suite.Suite
	ctx      context.Context
	sender   *kafkaapp.Producer
	consumer *kafkaapp.Consumer
}

func TestSenderSuite(t *testing.T) {
	suite.Run(t, new(SenderSuite))
}

func (s *SenderSuite) SetupSuite() {
	brokerAddr := os.Getenv("KAFKA_BROKER_ADDR")
	notifyTopic := os.Getenv("KAFKA_NOTIFY_TOPIC")
	messageTopic := os.Getenv("KAFKA_MESSAGE_TOPIC")
	if brokerAddr == "" {
		brokerAddr = "localhost:9092"
	}
	if notifyTopic == "" {
		notifyTopic = "notify"
	}
	if messageTopic == "" {
		messageTopic = "message"
	}

	logg := logger.New("info")
	s.sender = kafkaapp.NewProducer(brokerAddr, notifyTopic, logg)
	s.consumer = kafkaapp.NewConsumer(brokerAddr, messageTopic, logg)
}

func (s *SenderSuite) TearDownSuite() {
	s.sender.Stop()
	s.consumer.Stop()
}

func (s *SenderSuite) SetupTest() {

}

func (s *SenderSuite) TestSendEventSuccess() {
	event := scheduler.Notify{
		ID:    1,
		Title: "тестовое событие 1",
		Date:  time.Now(),
		User:  "Админ",
	}

	ctx := context.Background()
	err1 := s.sender.Send(ctx, event)
	s.Require().NoError(err1)
	time.Sleep(100 * time.Millisecond)

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	msg, err2 := s.consumer.Read(ctx)
	cancel()
	s.Require().NoError(err2)
	s.Require().NotNil(msg)
}
