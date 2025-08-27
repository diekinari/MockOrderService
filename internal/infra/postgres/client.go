package postgres

import (
	"MockOrderService/config"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

// Client represents a PostgreSQL client
type Client struct {
	Pool *pgxpool.Pool
}

// NewClient creates a new PostgreSQL client
func NewClient(cfg *config.Config, ctx context.Context) (*Client, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode)
	dbCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(dbCtx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %w", err)

	}
	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	return &Client{Pool: pool}, nil
}

// Close closes the connection pool
func (c *Client) Close() {
	c.Pool.Close()
}

// Ping pings the database
func (c *Client) Ping(ctx context.Context) error {
	return c.Pool.Ping(ctx)
}
