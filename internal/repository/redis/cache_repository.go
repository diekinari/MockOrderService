package redis

import (
	"MockOrderService/internal/domain/model"
	"MockOrderService/internal/monitoring"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheRepository struct {
	client *redis.Client
}

func NewCacheRepository(client *redis.Client) *CacheRepository {
	return &CacheRepository{client: client}
}

// SaveOrder saves order to cache with expiration time of 5 minutes
func (r *CacheRepository) SaveOrder(ctx context.Context, order *model.Order) (err error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		success := err == nil
		monitoring.RecordCacheOpDuration("save", success, duration)
	}()

	orderKey := fmt.Sprintf("order:%s", order.OrderUID)
	data, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("caching error – failed to marshal order: %w", err)
	}

	err = r.client.Set(ctx, orderKey, data, time.Minute*5).Err()
	if err != nil {
		return fmt.Errorf("caching error: %w", err)
	}

	return nil
}

// GetOrder returns order from cache if it exists, otherwise returns error
func (r *CacheRepository) GetOrder(ctx context.Context, orderUID string) (result *model.Order, err error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		// cache miss (redis.Nil) или успех - это нормально
		success := err == nil || errors.Is(err, redis.Nil)
		monitoring.RecordCacheOpDuration("get", success, duration)
	}()

	orderKey := fmt.Sprintf("order:%s", orderUID)
	// redis value is a json object, so we take bytes right away
	val, err := r.client.Get(ctx, orderKey).Bytes()
	if err != nil {
		// cache miss (redis.Nil) или ошибка
		if errors.Is(err, redis.Nil) {
			monitoring.RecordCacheMiss()
		}
		return nil, err
	}
	// cache hit
	monitoring.RecordCacheHit()

	var order model.Order
	err = json.Unmarshal(val, &order)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// IsCacheEmpty returns true if cache is empty, otherwise returns false
func (r *CacheRepository) IsCacheEmpty(ctx context.Context) (bool, error) {
	n, err := r.client.DBSize(ctx).Result()
	if err != nil {
		return false, err
	}
	return n == 0, nil
}
