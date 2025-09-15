package service

import (
	"MockOrderService/internal/domain/model"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
)

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
			s.sugar.Fatalw("CACHE HEAT-UP: failed to get fresh data from db", "error", err)
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
	err := s.orderRepo.SaveOrder(ctx, order)
	if err != nil {
		return fmt.Errorf("failed to save message to db – orderUID: %v – err: %w", order.OrderUID, err)
	}
	s.sugar.Infow("order was saved to db", "orderUID", order.OrderUID)

	err = s.cacheRepo.SaveOrder(ctx, order)
	if err != nil {
		s.sugar.Errorw("failed to cache order", "orderUID", order.OrderUID, "error", err)
		return nil
	}
	s.sugar.Infow("order was cached", "orderUID", order.OrderUID)
	return nil
}
