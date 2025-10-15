package service

import (
	"MockOrderService/internal/domain/model"
	"MockOrderService/internal/mocks"
	"context"
	"database/sql"
	"errors"
	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
	"testing"
)

func TestProcessOrder_DBError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sugar := zap.NewNop().Sugar()
	ord := &model.Order{OrderUID: "x"}

	orderRepo := mocks.NewMockOrderRepositoryService(ctrl)
	cacheRepo := mocks.NewMockCacheRepositoryService(ctrl)

	orderRepo.EXPECT().SaveOrder(gomock.Any(), ord).Return(errors.New("db down"))
	// cache should not be called

	svc := NewOrderService(sugar, orderRepo, cacheRepo)
	err := svc.ProcessOrder(context.Background(), ord)
	if err == nil {
		t.Fatalf("expected error from db save, got nil")
	}
}

func TestProcessOrder_CacheError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sugar := zap.NewNop().Sugar()
	ord := &model.Order{OrderUID: "ok"}

	orderRepo := mocks.NewMockOrderRepositoryService(ctrl)
	cacheRepo := mocks.NewMockCacheRepositoryService(ctrl)

	orderRepo.EXPECT().SaveOrder(gomock.Any(), ord).Return(nil)
	cacheRepo.EXPECT().SaveOrder(gomock.Any(), ord).Return(errors.New("redis err"))

	svc := NewOrderService(sugar, orderRepo, cacheRepo)
	err := svc.ProcessOrder(context.Background(), ord)
	if err != nil {
		t.Fatalf("expected nil when cache fails, got %v", err)
	}
}

func TestProcessOrder_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sugar := zap.NewNop().Sugar()
	ord := &model.Order{OrderUID: "ok"}

	orderRepo := mocks.NewMockOrderRepositoryService(ctrl)
	cacheRepo := mocks.NewMockCacheRepositoryService(ctrl)

	orderRepo.EXPECT().SaveOrder(gomock.Any(), ord).Return(nil)
	cacheRepo.EXPECT().SaveOrder(gomock.Any(), ord).Return(nil)

	svc := NewOrderService(sugar, orderRepo, cacheRepo)
	err := svc.ProcessOrder(context.Background(), ord)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestHeatUpCache_NotEmpty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sugar := zap.NewNop().Sugar()
	orderRepo := mocks.NewMockOrderRepositoryService(ctrl)
	cacheRepo := mocks.NewMockCacheRepositoryService(ctrl)

	cacheRepo.EXPECT().IsCacheEmpty(gomock.Any()).Return(false, nil)
	// Ensure DB is not called
	orderRepo.EXPECT().GetRecentOrders(gomock.Any(), gomock.Any()).Times(0)

	svc := NewOrderService(sugar, orderRepo, cacheRepo)
	svc.HeatUpCache(context.Background())
}

func TestHeatUpCache_DBEmpty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sugar := zap.NewNop().Sugar()
	orderRepo := mocks.NewMockOrderRepositoryService(ctrl)
	cacheRepo := mocks.NewMockCacheRepositoryService(ctrl)

	cacheRepo.EXPECT().IsCacheEmpty(gomock.Any()).Return(true, nil)
	orderRepo.EXPECT().GetRecentOrders(gomock.Any(), 5).Return(nil, sql.ErrNoRows)
	// No cache saves expected
	cacheRepo.EXPECT().SaveOrder(gomock.Any(), gomock.Any()).Times(0)

	svc := NewOrderService(sugar, orderRepo, cacheRepo)
	svc.HeatUpCache(context.Background())
}

func TestHeatUpCache_DBError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sugar := zap.NewNop().Sugar()
	orderRepo := mocks.NewMockOrderRepositoryService(ctrl)
	cacheRepo := mocks.NewMockCacheRepositoryService(ctrl)

	cacheRepo.EXPECT().IsCacheEmpty(gomock.Any()).Return(true, nil)
	orderRepo.EXPECT().GetRecentOrders(gomock.Any(), 5).Return(nil, errors.New("db fail"))

	svc := NewOrderService(sugar, orderRepo, cacheRepo)
	// Function logs fatalw in code path; with Nop logger this won't terminate.
	svc.HeatUpCache(context.Background())
}

func TestHeatUpCache_CacheSaves_PartialFailures(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sugar := zap.NewNop().Sugar()
	orderRepo := mocks.NewMockOrderRepositoryService(ctrl)
	cacheRepo := mocks.NewMockCacheRepositoryService(ctrl)

	orders := []*model.Order{{OrderUID: "o1"}, {OrderUID: "o2"}, {OrderUID: "o3"}}

	cacheRepo.EXPECT().IsCacheEmpty(gomock.Any()).Return(true, nil)
	orderRepo.EXPECT().GetRecentOrders(gomock.Any(), 5).Return(orders, nil)
	// Save o1 ok, o2 fails, o3 ok; handler should continue on errors
	cacheRepo.EXPECT().SaveOrder(gomock.Any(), orders[0]).Return(nil)
	cacheRepo.EXPECT().SaveOrder(gomock.Any(), orders[1]).Return(errors.New("cache fail"))
	cacheRepo.EXPECT().SaveOrder(gomock.Any(), orders[2]).Return(nil)

	svc := NewOrderService(sugar, orderRepo, cacheRepo)
	svc.HeatUpCache(context.Background())
}
