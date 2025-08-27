package redis

import (
	"MockOrderService/internal/domain/model"
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type CacheRepository struct {
	client *redis.Client
}

func NewCacheRepository(client *redis.Client) *CacheRepository {
	return &CacheRepository{client: client}
}

// SaveOrder saves order to cache with expiration time of 5 minutes
func (r *CacheRepository) SaveOrder(ctx context.Context, order *model.Order) error {
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
func (r *CacheRepository) GetOrder(ctx context.Context, orderUID string) (*model.Order, error) {
	orderKey := fmt.Sprintf("order:%s", orderUID)
	// redis value is a json object, so we take bytes right away
	val, err := r.client.Get(ctx, orderKey).Bytes()
	if err != nil {
		// cache miss
		return nil, err
	}
	// cache hit
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
