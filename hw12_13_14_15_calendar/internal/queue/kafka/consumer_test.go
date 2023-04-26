package kafkaapp

import (
	"context"
	"github.com/dimonk33/kdv_hw_otus/hw12_13_14_15_calendar/internal/logger"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

//type testData struct {
//	Text string    `json:"text"`
//	Date time.Time `json:"date"`
//}

func TestConsumer(t *testing.T) {

	t.Run("read message", func(t *testing.T) {
		//p := NewProducer("localhost:9092", "test", logger.New(logger.LevelDebug))
		//require.NotNil(t, p)
		//p.start()
		//defer p.stop()
		//data := testData{
		//	Text: "Hello",
		//	Date: time.Now(),
		//}
		ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)
		//err := p.send(ctx, data)
		//require.Nil(t, err)

		c := NewConsumer("localhost:9092", "test", logger.New(logger.LevelDebug))
		c.Start()
		defer c.Stop()
		rdata, err := c.Read(ctx)
		require.Nil(t, err)
		require.NotNil(t, rdata)
	})
}
