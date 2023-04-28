package kafkaapp

import (
	"context"
	"fmt"

	"github.com/dimonk33/kdv_hw_otus/hw12_13_14_15_calendar/internal/logger"
	kafka "github.com/segmentio/kafka-go"
)

type Consumer struct {
	logger *logger.Logger
	reader *kafka.Reader
}

func NewConsumer(brokerAddr string, topic string, logger *logger.Logger) *Consumer {
	c := Consumer{logger: logger}
	c.reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{brokerAddr},
		Topic:     topic,
		Partition: 0,
		MinBytes:  10e3, // 10KB
		MaxBytes:  10e6, // 10MB
	})

	return &c
}

func (c *Consumer) Stop() {
	if err := c.reader.Close(); err != nil {
		c.logger.Warning(fmt.Sprintf("%v", err))
	}
}

func (c *Consumer) Read(ctx context.Context) ([]byte, error) {
	msg, err := c.reader.ReadMessage(ctx)
	if err != nil {
		return nil, err
	}

	return msg.Value, nil
}
