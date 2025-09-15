package monitoring

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
)

type HealthCheckable interface {
	Ping(ctx context.Context) error
}

// HealthChecker checks the health of the database and cache
type HealthChecker struct {
	dbClient    HealthCheckable
	cacheClient HealthCheckable
	interval    time.Duration
	sugar       *zap.SugaredLogger
	stop        context.CancelFunc
}

// NewHealthChecker creates a new health checker
func NewHealthChecker(dbClient HealthCheckable, cacheClient HealthCheckable, interval time.Duration, logger *zap.SugaredLogger, cancelFunc context.CancelFunc) *HealthChecker {
	return &HealthChecker{
		dbClient:    dbClient,
		cacheClient: cacheClient,
		interval:    interval,
		sugar:       logger,
		stop:        cancelFunc,
	}
}

// Start starts health checking
func (h *HealthChecker) Start(ctx context.Context) {
	ticker := time.NewTicker(h.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			h.sugar.Info("checking health")
			if err := h.checkHealth(ctx); err != nil {
				h.sugar.Fatalw("healthcheck failed", "err", err)
				h.stop()
				return
			}
			h.sugar.Info("healthcheck pased")
		}
	}
}

// checkHealth checks the health of the database and cache
func (h *HealthChecker) checkHealth(ctx context.Context) error {
	// Проверка базы данных
	if err := h.dbClient.Ping(ctx); err != nil {
		return fmt.Errorf("db healthcheck failed: %w", err)
	}

	// Проверка кэша
	if err := h.cacheClient.Ping(ctx); err != nil {
		return fmt.Errorf("cache healthcheck failed: %w", err)
	}

	return nil
}
