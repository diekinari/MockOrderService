package redis

import (
	"MockOrderService/config"
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Client represents a Redis client
type Client struct {
	Client *redis.Client
}

// NewClient creates a new Redis client
func NewClient(cfg *config.Config, ctx context.Context) (*Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisHost,
		Password: cfg.RedisPassword,
		DB:       0,
	})

	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := client.Ping(timeoutCtx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &Client{Client: client}, nil
}

// Close closes the connection with Redis
func (c *Client) Close() error {
	return c.Client.Close()
}

// Ping pings the Redis server
func (c *Client) Ping(ctx context.Context) error {
	return c.Client.Ping(ctx).Err()
}
