// Package kafka implements a Kafka consumer for processing order-related messages.
// It provides functionality to consume messages from Kafka topics, process order data,
// and handle message commits and error scenarios with retry mechanisms.
package kafka

import (
	"MockOrderService/internal/domain/model"
	"MockOrderService/internal/monitoring"
	"MockOrderService/internal/service"
	"MockOrderService/internal/validation"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.uber.org/zap"
)

type consumerClient interface {
	ReadMessage(ctx context.Context) (kafka.Message, error)
	CommitMessages(ctx context.Context, messages ...kafka.Message) error
	Topic() string
	WriteMessagesToDLQ(ctx context.Context, errorInfo string, topic string, retryCount int, messages ...kafka.Message) error
}

// Consumer represents a Kafka consumer
type Consumer struct {
	client            consumerClient
	service           *service.OrderService
	sugar             *zap.SugaredLogger
	errorsCount       int
	maxRetriesDLQ     int
	retryBackoffMsDLQ int
}

// NewConsumer creates a new Kafka consumer with the given client, service, and logger.
func NewConsumer(client consumerClient, service *service.OrderService, maxRetriesDLQ int, retryBackoffMsDLQ int, sugar *zap.SugaredLogger) *Consumer {
	return &Consumer{client: client, service: service, maxRetriesDLQ: maxRetriesDLQ, retryBackoffMsDLQ: retryBackoffMsDLQ, sugar: sugar}
}

// Start functions starts a consumer. It reads the messages and process them accordingly with provided method.
func (c *Consumer) Start(ctx context.Context, stop context.CancelFunc) {
	for {
		select {
		case <-ctx.Done():
			c.sugar.Fatalw("context canceled", "error", ctx.Err())
			return
		default:
			msg, err := c.client.ReadMessage(ctx)
			if err != nil {
				c.errorsCount++
				c.sugar.Errorw("failed to read message", "error", err)
				// Записываем метрику с success=false при ошибке чтения
				monitoring.RecordKafkaMessagesConsumed(c.client.Topic(), false)
				if c.errorsCount > 3 {
					c.sugar.Fatal("consumer has reached maximum amount of errors, stopping the service")
					stop()
					return
				}
				continue
			}
			// Kafka работает нормально - сообщение прочитано успешно (техническая успешность)
			// success = true означает, что Kafka работает, а не то, что обработка удалась
			monitoring.RecordKafkaMessagesConsumed(c.client.Topic(), true)
			if err := c.processMessage(ctx, msg); err != nil {
				c.sugar.Errorw("failed to process message", "error", err)
				// Бизнес-ошибки отслеживаются через business_metrics.go (orders_failed_total)
			}

		}
	}
}

func (c *Consumer) processMessage(ctx context.Context, msg kafka.Message) error {
	// Извлекаем trace context из заголовков Kafka сообщения
	// Это позволяет связать обработку сообщения с трейсом, который создал сообщение
	propagator := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
	carrier := &KafkaMessageCarrier{msg: &msg}
	ctx = propagator.Extract(ctx, carrier)

	// Создаем span для обработки Kafka сообщения
	ctx, span := monitoring.Tracer.Start(ctx, "kafka.ProcessMessage")
	defer span.End()

	var order model.Order
	err := json.Unmarshal(msg.Value, &order)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to unmarshal order")

		errDLQ := c.client.WriteMessagesToDLQ(ctx, "failed to unmarshal order", msg.Topic, 0, msg)
		monitoring.RecordKafkaMessagesSentToDQL(msg.Topic, err.Error())
		if errDLQ != nil {
			span.RecordError(errDLQ)
			return errors.Join(err, errDLQ)
		}
		return fmt.Errorf("failed to unmarshal order: %w", err)
	}

	// Добавляем атрибуты к span
	span.SetAttributes(
		attribute.String("kafka.topic", msg.Topic),
		attribute.Int("kafka.partition", msg.Partition),
		attribute.Int64("kafka.offset", msg.Offset),
		attribute.String("kafka.order_uid", order.OrderUID),
	)

	c.sugar.Infow("order consumed", "orderUID", order.OrderUID)

	err = validation.ValidateOrder(&order)
	if err != nil {
		// edgy case: it's not an error actually, but we can't continue processing
		span.SetAttributes(attribute.Bool("validation.valid", false))
		span.SetStatus(codes.Ok, "invalid order - skipped")
		c.sugar.Warnw("invalid order", "orderUID", order.OrderUID)
		return nil
	}
	span.SetAttributes(attribute.Bool("validation.valid", true))
	c.sugar.Infow("order is validated", "orderUID", order.OrderUID)

	for i := 0; i < c.maxRetriesDLQ; i++ {
		if err := c.service.ProcessOrder(ctx, &order); err != nil {
			monitoring.RecordKafkaRetries(msg.Topic, i)
			if i == c.maxRetriesDLQ-1 {
				span.RecordError(err)
				span.SetStatus(codes.Error, "failed to process order")
				errDLQ := c.client.WriteMessagesToDLQ(ctx, "failed to process order", msg.Topic, i+1, msg)
				monitoring.RecordKafkaMessagesSentToDQL(msg.Topic, err.Error())
				if errDLQ != nil {
					span.RecordError(errDLQ)
					return errors.Join(err, errDLQ)
				}
				return err
			}
			time.Sleep(time.Duration(c.retryBackoffMsDLQ) * time.Millisecond)
			continue
		}
		break
	}

	for i := 0; i < c.maxRetriesDLQ; i++ {
		if err := c.client.CommitMessages(ctx, msg); err != nil {
			monitoring.RecordKafkaRetries(msg.Topic, i)
			if i == c.maxRetriesDLQ-1 {
				span.RecordError(err)
				span.SetStatus(codes.Error, "failed to commit message")
				errDLQ := c.client.WriteMessagesToDLQ(ctx, "failed to commit message", msg.Topic, i+1, msg)
				monitoring.RecordKafkaMessagesSentToDQL(msg.Topic, err.Error())
				if errDLQ != nil {
					span.RecordError(errDLQ)
					return errors.Join(err, errDLQ)
				}
				return err
			}
			time.Sleep(time.Duration(c.retryBackoffMsDLQ) * time.Millisecond)
			continue
		}
		break
	}
	span.SetStatus(codes.Ok, "order processed and committed successfully")
	c.sugar.Infow("order was committed", "orderUID", order.OrderUID)

	return nil
}

// KafkaMessageCarrier реализует TextMapCarrier для передачи trace context через Kafka заголовки
type KafkaMessageCarrier struct {
	msg *kafka.Message
}

// Get возвращает значение заголовка по ключу
func (c *KafkaMessageCarrier) Get(key string) string {
	for _, h := range c.msg.Headers {
		if h.Key == key {
			return string(h.Value)
		}
	}
	return ""
}

// Set устанавливает значение заголовка
func (c *KafkaMessageCarrier) Set(key, value string) {
	c.msg.Headers = append(c.msg.Headers, kafka.Header{
		Key:   key,
		Value: []byte(value),
	})
}

// Keys возвращает все ключи заголовков
func (c *KafkaMessageCarrier) Keys() []string {
	keys := make([]string, len(c.msg.Headers))
	for i, h := range c.msg.Headers {
		keys[i] = h.Key
	}
	return keys
}
