package kafka

import (
	faker "MockOrderService/internal/faker"
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type producerClient interface {
	WriteMessages(ctx context.Context, messages ...kafka.Message) error
}

// Producer represents a Kafka producer
type Producer struct {
	client      producerClient
	sugar       *zap.SugaredLogger
	errorsCount int
}

// NewProducer creates a new Kafka producer with the given client and logger.
func NewProducer(client producerClient, sugar *zap.SugaredLogger) *Producer {
	return &Producer{
		client:      client,
		sugar:       sugar,
		errorsCount: 0,
	}
}

// Start launches the producer.
// Producer emulates limited independent random order generator.
// Perhaps something inside should be in the main service, but it's too short.
func (p *Producer) Start(stop context.CancelFunc) {
	for i := 0; i < 15; i++ {
		order := faker.NewFakeOrder()
		value, err := json.Marshal(order)
		if err != nil {
			p.sugar.Errorw("Failed to marshal order", "orderUID", order.OrderUID, "error", err)
			continue
		}

		err = p.client.WriteMessages(context.Background(), kafka.Message{
			Key:   []byte(order.OrderUID),
			Value: value,
		})
		if err != nil {
			p.errorsCount++
			p.sugar.Errorw("failed to write messages", "orderUID", order.OrderUID, "error", err)
			if p.errorsCount > 3 {
				p.sugar.Fatal("producer has reached maximum amount of writing errors, stopping the service")
				stop()
				return
			}
			continue
		}
		p.sugar.Infow("order produced", "orderUID", order.OrderUID)
	}
	p.sugar.Infow("producer has finished")
}
