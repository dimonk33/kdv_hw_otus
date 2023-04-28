package kafkaapp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/dimonk33/kdv_hw_otus/hw12_13_14_15_calendar/internal/logger"
	kafka "github.com/segmentio/kafka-go"
)

type Producer struct {
	logger *logger.Logger
	writer *kafka.Writer
}

func NewProducer(brokerAddr string, topic string, logger *logger.Logger) *Producer {
	p := Producer{logger: logger}
	p.writer = &kafka.Writer{
		Addr:                   kafka.TCP(brokerAddr),
		Topic:                  topic,
		AllowAutoTopicCreation: true,
		Balancer:               &kafka.LeastBytes{},
	}

	return &p
}

func (p *Producer) Stop() {
	if err := p.writer.Close(); err != nil {
		p.logger.Warning(fmt.Sprintf("%v", err))
	}
}

func (p *Producer) Send(ctx context.Context, data interface{}) error {
	sendData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	const retries = 3
	var cancel context.CancelFunc
	for i := 0; i < retries; i++ {
		ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
		err = p.writer.WriteMessages(ctx, kafka.Message{Value: sendData})
		if errors.Is(err, kafka.UnknownTopicOrPartition) || errors.Is(err, context.DeadlineExceeded) {
			time.Sleep(time.Millisecond * 250)
			continue
		}
		if err != nil {
			cancel()
			return err
		}
		break
	}

	cancel()
	return nil
}
