package kafkaapp

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dimonk33/kdv_hw_otus/hw12_13_14_15_calendar/internal/logger"
	"github.com/segmentio/kafka-go"
)

type Producer struct {
	logger *logger.Logger
	topic  string
	writer *kafka.Writer
}

func NewProducer(brokerAddr string, topic string, _logger *logger.Logger) *Producer {
	p := Producer{logger: _logger}
	p.writer = &kafka.Writer{
		Addr:                   kafka.TCP(brokerAddr),
		Topic:                  topic,
		AllowAutoTopicCreation: true,
		Balancer:               &kafka.LeastBytes{},
	}

	return &p
}

func (p *Producer) Start() {

}

func (p *Producer) Stop() {
	err := p.writer.Close()
	if err != nil {
		p.logger.Warning(fmt.Sprintf("%v", err))
	}
}

func (p *Producer) Send(ctx context.Context, data any) error {
	sendData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = p.writer.WriteMessages(ctx, kafka.Message{Value: sendData})
	if err != nil {
		return err
	}

	return nil
}
