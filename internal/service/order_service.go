package service

import (
	"MockOrderService/internal/domain/model"
	"MockOrderService/internal/monitoring"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
)

//go:generate mockgen -source=order_service.go -destination=../mocks/mock_order_service.go -package=mocks -mock_names OrderRepository=MockOrderRepositoryService,CacheRepository=MockCacheRepositoryService
type OrderRepository interface {
	GetRecentOrders(ctx context.Context, limit int) ([]*model.Order, error)
	SaveOrder(ctx context.Context, order *model.Order) error
}

type CacheRepository interface {
	SaveOrder(ctx context.Context, order *model.Order) error
	IsCacheEmpty(ctx context.Context) (bool, error)
}

type OrderService struct {
	sugar     *zap.SugaredLogger
	orderRepo OrderRepository
	cacheRepo CacheRepository
}

func NewOrderService(sugar *zap.SugaredLogger, orderRepo OrderRepository, cacheRepo CacheRepository) *OrderService {
	return &OrderService{
		sugar:     sugar,
		orderRepo: orderRepo,
		cacheRepo: cacheRepo,
	}
}
func (s *OrderService) HeatUpCache(ctx context.Context) {
	s.sugar.Info("starting heating up cache...")
	timeoutCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	cacheIsEmpty, err := s.cacheRepo.IsCacheEmpty(timeoutCtx)
	if err != nil {
		s.sugar.Errorw("CACHE HEAT-UP: failed to check if cache is empty", "error", err)
		return
	}
	if cacheIsEmpty {
		orders, err := s.orderRepo.GetRecentOrders(timeoutCtx, 5)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				s.sugar.Infow("CACHE HEAT-UP: database is empty, cache remains empty too", "error", err)
				return
			}
			s.sugar.Errorw("CACHE HEAT-UP: failed to get fresh data from db", "error", err)
			return
		}
		for _, order := range orders {
			if err = s.cacheRepo.SaveOrder(ctx, order); err != nil {
				s.sugar.Errorw("CACHE HEAT-UP: failed to save order", "error", err)
				continue
			}
			s.sugar.Infow("CACHE HEAT-UP: order was cached", "orderUID", order.OrderUID)
		}
		s.sugar.Infow("cache heat-up completed", "total_orders", len(orders))
	} else {
		s.sugar.Info("CACHE HEAT-UP: cache is not empty")
	}
}

// ProcessOrder processes an order.
// Keep in mind: any returning error will result in skipping commiting the message.
// This is a design choice to ensure message integrity.
// Business rules imply that we should commit the message after it being saved to db, regardless of caching.
// So error is returned in case of failure to save the order to db.
// But there is no returning error in case of failure to save the order to cache.
func (s *OrderService) ProcessOrder(ctx context.Context, order *model.Order) error {
	// Создаем span для обработки заказа
	ctx, span := monitoring.Tracer.Start(ctx, "service.ProcessOrder")
	defer span.End()

	start := time.Now()

	// Добавляем атрибуты к span
	span.SetAttributes(
		attribute.String("order.uid", order.OrderUID),
		attribute.String("order.track_number", order.TrackNumber),
	)

	defer func() {
		duration := time.Since(start)
		span.SetAttributes(attribute.Int64("order.process_duration_ms", duration.Milliseconds()))
		monitoring.RecordOrderProcessDuration(duration)
	}()

	err := s.orderRepo.SaveOrder(ctx, order)
	if err != nil {
		// Ошибка сохранения в БД - заказ не обработан
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to save order to database")
		span.SetAttributes(attribute.Bool("order.db_saved", false))
		monitoring.RecordOrderFailed()
		return fmt.Errorf("failed to save message to db – orderUID: %v – err: %w", order.OrderUID, err)
	}
	span.SetAttributes(attribute.Bool("order.db_saved", true))
	s.sugar.Infow("order was saved to db", "orderUID", order.OrderUID)

	err = s.cacheRepo.SaveOrder(ctx, order)
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.Bool("order.cache_saved", false))
		s.sugar.Errorw("failed to cache order", "orderUID", order.OrderUID, "error", err)
		// Ошибка кэша не критична - заказ сохранен в БД, считаем успешной обработкой
		span.SetStatus(codes.Ok, "order saved to DB, cache failed (non-critical)")
		monitoring.RecordOrderCreated()
		return nil
	}
	span.SetAttributes(attribute.Bool("order.cache_saved", true))
	s.sugar.Infow("order was cached", "orderUID", order.OrderUID)

	// Успешная обработка - заказ сохранен в БД и кэш
	span.SetStatus(codes.Ok, "order processed successfully")
	monitoring.RecordOrderCreated()
	return nil
}
