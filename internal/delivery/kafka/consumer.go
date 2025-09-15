package kafka

import (
	"MockOrderService/internal/domain/model"
	"MockOrderService/internal/service"
	"MockOrderService/internal/validation"
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type consumerClient interface {
	ReadMessage(ctx context.Context) (kafka.Message, error)
	CommitMessages(ctx context.Context, messages ...kafka.Message) error
}

// Consumer represents a Kafka consumer
type Consumer struct {
	client      consumerClient
	service     *service.OrderService
	sugar       *zap.SugaredLogger
	errorsCount int
}

func NewConsumer(client consumerClient, service *service.OrderService, sugar *zap.SugaredLogger) *Consumer {
	return &Consumer{client: client, service: service, sugar: sugar}
}

func (c *Consumer) Start(ctx context.Context, stop context.CancelFunc) {
	for {
		select {
		case <-ctx.Done():
			c.sugar.Fatalw("context canceled", "error", ctx.Err())
			return
		default:
			msg, err := c.client.ReadMessage(ctx)
			if err != nil {
				c.errorsCount += 1
				c.sugar.Errorw("failed to read message", "error", err)
				if c.errorsCount > 3 {
					c.sugar.Fatal("consumer has reached maximum amount of errors, stopping the service")
					stop()
					return
				}
				continue
			}
			if err := c.processMessage(ctx, msg); err != nil {
				c.sugar.Errorw("failed to process message", "error", err)
			}

		}
	}
}

func (c *Consumer) processMessage(ctx context.Context, msg kafka.Message) error {
	var order model.Order
	err := json.Unmarshal(msg.Value, &order)
	if err != nil {
		return fmt.Errorf("failed to unmarshal order: %w", err)
	}
	c.sugar.Infow("order consumed", "orderUID", order.OrderUID)

	err = validation.ValidateOrder(&order)
	if err != nil {
		// edgy case: it's not an error actually, but we can't continue processing
		c.sugar.Warnw("invalid order", "orderUID", order.OrderUID)
		return nil
	}
	c.sugar.Infow("order is validated", "orderUID", order.OrderUID)

	if err := c.service.ProcessOrder(ctx, &order); err != nil {
		return err
	}

	err = c.client.CommitMessages(ctx, msg)
	if err != nil {
		return fmt.Errorf("failed to commit message: %w", err)
	}
	c.sugar.Infow("order was committed", "orderUID", order.OrderUID)

	return nil
}
