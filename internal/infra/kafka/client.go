package kafka

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"
)

// Client represents a Kafka client
type Client struct {
	reader *kafka.Reader
	writer *kafka.Writer
}

func NewClient(broker string, groupID string, topic string) (*Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := CreateTopicIfNotExists(ctx, broker, topic, 3, 1); err != nil {
		return nil, fmt.Errorf("ensure topic: %w", err)
	}

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
	}, nil
}

// ReadMessage reads a message as a kafka reader
func (c *Client) ReadMessage(ctx context.Context) (kafka.Message, error) {
	return c.reader.ReadMessage(ctx)
}

// CommitMessages commits messages as a kafka reader
func (c *Client) CommitMessages(ctx context.Context, messages ...kafka.Message) error {
	return c.reader.CommitMessages(ctx, messages...)
}

// WriteMessages writes several messages as a kafka writer
func (c *Client) WriteMessages(ctx context.Context, messages ...kafka.Message) error {
	return c.writer.WriteMessages(ctx, messages...)
}

// Close closes kafka client completely
func (c *Client) Close() error {
	rError := fmt.Errorf("reader: %w", c.reader.Close())
	wError := fmt.Errorf("writer: %w", c.writer.Close())
	return errors.Join(rError, wError)
}

// CreateTopicIfNotExists creates a topic.
func CreateTopicIfNotExists(ctx context.Context, broker string, topic string, partitions int, replicationFactor int) (resErr error) {
	// short helper to dial any broker
	dialAny := func(broker string) (*kafka.Conn, error) {
		// kafka.DialContext resolves host/port
		return kafka.DialContext(ctx, "tcp", broker)
	}

	// try to connect to broker to get controller
	var conn *kafka.Conn
	var err error
	conn, err = dialAny(broker)
	if conn == nil {
		return fmt.Errorf("unable to dial any broker: last err: %w", err)
	}
	defer func() {
		if cerr := conn.Close(); cerr != nil {
			if resErr == nil {
				resErr = fmt.Errorf("close broker conn: %w", cerr)
			} else {
				resErr = errors.Join(resErr, fmt.Errorf("close broker conn: %w", cerr))
			}
		}
	}()

	controller, err := conn.Controller()
	if err != nil {
		return fmt.Errorf("get controller: %w", err)
	}
	controllerAddr := net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port))

	ctrl, err := kafka.DialContext(ctx, "tcp", controllerAddr)
	if err != nil {
		return fmt.Errorf("dial controller %s: %w", controllerAddr, err)
	}
	defer func() {
		if cerr := ctrl.Close(); cerr != nil {
			if resErr == nil {
				resErr = fmt.Errorf("close controller conn: %w", cerr)
			} else {
				resErr = errors.Join(resErr, fmt.Errorf("close controller conn: %w", cerr))
			}
		}
	}()

	tc := kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     partitions,
		ReplicationFactor: replicationFactor,
	}
	if err := ctrl.CreateTopics(tc); err != nil {
		// CreateTopics can return an error even if the topic exists or was created concurrently.
		// Here we try to verify existence by reading partitions.
		parts, rerr := ctrl.ReadPartitions(topic)
		if rerr != nil {
			return fmt.Errorf("create topic error: %w; additionally ReadPartitions failed: %v", err, rerr)
		}
		if len(parts) == 0 {
			return fmt.Errorf("create topic error and no partitions found: %w", err)
		}
		// otherwise topic exists — treat as success
	}

	// wait until metadata is visible (some brokers may need a short moment)
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		// try broker to read partitions
		ok := false
		c, derr := kafka.DialContext(ctx, "tcp", broker)
		if derr != nil {
			continue
		}
		parts, rerr := c.ReadPartitions(topic)
		_ = c.Close()
		if rerr != nil {
			continue
		}
		if len(parts) > 0 {
			ok = true
		}
		if ok {
			return nil
		}
		time.Sleep(250 * time.Millisecond)
	}
	return errors.New("topic creation timed out waiting for metadata propagation")
}
