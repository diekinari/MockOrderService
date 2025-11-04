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
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type CacheRepository struct {
	client *redis.Client
}

func NewCacheRepository(client *redis.Client) *CacheRepository {
	return &CacheRepository{client: client}
}

// SaveOrder saves order to cache with expiration time of 5 minutes
func (r *CacheRepository) SaveOrder(ctx context.Context, order *model.Order) (err error) {
	// Создаем span для операции сохранения в кэш
	ctx, span := monitoring.Tracer.Start(ctx, "redis.SaveOrder")
	defer span.End()

	start := time.Now()
	defer func() {
		duration := time.Since(start)
		success := err == nil

		// Добавляем атрибуты к span
		span.SetAttributes(
			attribute.String("cache.operation", "set"),
			attribute.String("cache.key", fmt.Sprintf("order:%s", order.OrderUID)),
			attribute.String("cache.order_uid", order.OrderUID),
			attribute.Bool("cache.success", success),
			attribute.Int64("cache.duration_ms", duration.Milliseconds()),
		)

		monitoring.RecordCacheOpDuration("save", success, duration)

		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		} else {
			span.SetStatus(codes.Ok, "order cached successfully")
		}
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
	// Создаем span для операции получения из кэша
	ctx, span := monitoring.Tracer.Start(ctx, "redis.GetOrder")
	defer span.End()

	start := time.Now()
	defer func() {
		duration := time.Since(start)
		// cache miss (redis.Nil) или успех - это нормально
		success := err == nil || errors.Is(err, redis.Nil)
		cacheHit := err == nil
		cacheMiss := errors.Is(err, redis.Nil)

		// Добавляем атрибуты к span
		span.SetAttributes(
			attribute.String("cache.operation", "get"),
			attribute.String("cache.key", fmt.Sprintf("order:%s", orderUID)),
			attribute.String("cache.order_uid", orderUID),
			attribute.Bool("cache.success", success),
			attribute.Bool("cache.hit", cacheHit),
			attribute.Bool("cache.miss", cacheMiss),
			attribute.Int64("cache.duration_ms", duration.Milliseconds()),
		)

		monitoring.RecordCacheOpDuration("get", success, duration)

		if cacheMiss {
			monitoring.RecordCacheMiss()
			span.SetStatus(codes.Ok, "cache miss")
		} else if cacheHit {
			monitoring.RecordCacheHit()
			span.SetStatus(codes.Ok, "cache hit")
		} else {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
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
