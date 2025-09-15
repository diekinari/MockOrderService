package kafka

import (
	"context"
	"errors"
	"fmt"

	"github.com/segmentio/kafka-go"
)

// Client represents a Kafka client
type Client struct {
	reader *kafka.Reader
	writer *kafka.Writer
}

func NewClient(broker string, groupID string, topic string) *Client {
	return &Client{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:        []string{broker},
			GroupID:        groupID,
			Topic:          topic,
			CommitInterval: 0,
		}),
		writer: &kafka.Writer{
			Addr:     kafka.TCP(broker),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (c *Client) ReadMessage(ctx context.Context) (kafka.Message, error) {
	return c.reader.ReadMessage(ctx)
}

func (c *Client) CommitMessages(ctx context.Context, messages ...kafka.Message) error {
	return c.reader.CommitMessages(ctx, messages...)
}

func (c *Client) WriteMessages(ctx context.Context, messages ...kafka.Message) error {
	return c.writer.WriteMessages(ctx, messages...)
}

func (c *Client) Close() error {
	rError := fmt.Errorf("reader: %w", c.reader.Close())
	wError := fmt.Errorf("writer: %w", c.writer.Close())
	return errors.Join(rError, wError)
}
