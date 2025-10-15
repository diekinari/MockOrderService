package http

import (
	"MockOrderService/internal/domain/model"
	"MockOrderService/internal/mocks"
	"context"
	"encoding/json"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestServer(t *testing.T) (*ApiServer, *gomock.Controller, *mocks.MockOrderRepository, *mocks.MockCacheRepository) {
	ctrl := gomock.NewController(t)
	orderRepo := mocks.NewMockOrderRepository(ctrl)
	cacheRepo := mocks.NewMockCacheRepository(ctrl)
	sugar := zap.NewNop().Sugar()
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	return NewApiServer(sugar, ctx, orderRepo, cacheRepo), ctrl, orderRepo, cacheRepo
}

func TestHandleOrder_EmptyOrderUID(t *testing.T) {
	as, ctrl, _, _ := newTestServer(t)
	defer ctrl.Finish()

	req := httptest.NewRequest(http.MethodGet, "/api/order/", nil)
	// Inject empty path var
	req = mux.SetURLVars(req, map[string]string{"orderUID": ""})
	rr := httptest.NewRecorder()

	as.handleOrder(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
	var ae apiError
	if err := json.Unmarshal(rr.Body.Bytes(), &ae); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if ae.Error != "empty orderUID" {
		t.Fatalf("expected error 'empty orderUID', got %q", ae.Error)
	}
}

func TestHandleOrder_NotFound(t *testing.T) {
	as, ctrl, orderRepo, cacheRepo := newTestServer(t)
	defer ctrl.Finish()

	orderUID := "unknown"
	// Cache miss
	cacheRepo.EXPECT().GetOrder(gomock.Any(), orderUID).Return(nil, redis.Nil)
	// DB not found
	orderRepo.EXPECT().GetOrderByOrderUID(gomock.Any(), orderUID).Return(nil, pgx.ErrNoRows)

	req := httptest.NewRequest(http.MethodGet, "/api/order/"+orderUID, nil)
	// Use router-style vars so handler can read it
	req = mux.SetURLVars(req, map[string]string{"orderUID": orderUID})
	rr := httptest.NewRecorder()

	as.handleOrder(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rr.Code)
	}
	var ae apiError
	if err := json.Unmarshal(rr.Body.Bytes(), &ae); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if ae.Error != "no order found" {
		t.Fatalf("expected error 'no order found', got %q", ae.Error)
	}
}

func TestHandleOrder_InternalError(t *testing.T) {
	as, ctrl, _, cacheRepo := newTestServer(t)
	defer ctrl.Finish()

	orderUID := "id123"
	cacheRepo.EXPECT().GetOrder(gomock.Any(), orderUID).Return(nil, errors.New("redis down"))

	req := httptest.NewRequest(http.MethodGet, "/api/order/"+orderUID, nil)
	req = mux.SetURLVars(req, map[string]string{"orderUID": orderUID})
	rr := httptest.NewRecorder()

	as.handleOrder(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rr.Code)
	}
	var ae apiError
	if err := json.Unmarshal(rr.Body.Bytes(), &ae); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if ae.Error != "couldn't get order" {
		t.Fatalf("expected error 'couldn't get order', got %q", ae.Error)
	}
}

func TestHandleOrder_OK(t *testing.T) {
	as, ctrl, _, cacheRepo := newTestServer(t)
	defer ctrl.Finish()

	orderUID := "ok-1"
	ord := &model.Order{OrderUID: orderUID}
	cacheRepo.EXPECT().GetOrder(gomock.Any(), orderUID).Return(ord, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/order/"+orderUID, nil)
	req = mux.SetURLVars(req, map[string]string{"orderUID": orderUID})
	rr := httptest.NewRecorder()

	as.handleOrder(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
	var got model.Order
	if err := json.Unmarshal(rr.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to decode order: %v", err)
	}
	if got.OrderUID != orderUID {
		t.Fatalf("expected order_uid %q, got %q", orderUID, got.OrderUID)
	}
}

func TestGetOrder_CacheHit(t *testing.T) {
	as, ctrl, _, cacheRepo := newTestServer(t)
	defer ctrl.Finish()

	orderUID := "cache-hit"
	ord := &model.Order{OrderUID: orderUID}
	cacheRepo.EXPECT().GetOrder(gomock.Any(), orderUID).Return(ord, nil)

	res, err := as.getOrder(orderUID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res == nil || res.OrderUID != orderUID {
		t.Fatalf("unexpected result: %#v", res)
	}
}

func TestGetOrder_CacheMiss_DBHit(t *testing.T) {
	as, ctrl, orderRepo, cacheRepo := newTestServer(t)
	defer ctrl.Finish()

	orderUID := "cache-miss-db-hit"
	ord := &model.Order{OrderUID: orderUID}
	cacheRepo.EXPECT().GetOrder(gomock.Any(), orderUID).Return(nil, redis.Nil)
	orderRepo.EXPECT().GetOrderByOrderUID(gomock.Any(), orderUID).Return(ord, nil)

	res, err := as.getOrder(orderUID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res == nil || res.OrderUID != orderUID {
		t.Fatalf("unexpected result: %#v", res)
	}
}

func TestGetOrder_CacheMiss_DBNotFound(t *testing.T) {
	as, ctrl, orderRepo, cacheRepo := newTestServer(t)
	defer ctrl.Finish()

	orderUID := "cache-miss-db-404"
	cacheRepo.EXPECT().GetOrder(gomock.Any(), orderUID).Return(nil, redis.Nil)
	orderRepo.EXPECT().GetOrderByOrderUID(gomock.Any(), orderUID).Return(nil, pgx.ErrNoRows)

	res, err := as.getOrder(orderUID)
	if err == nil {
		t.Fatalf("expected error, got nil (res=%#v)", res)
	}
}

func TestGetOrder_CacheError(t *testing.T) {
	as, ctrl, _, cacheRepo := newTestServer(t)
	defer ctrl.Finish()

	orderUID := "cache-error"
	cacheRepo.EXPECT().GetOrder(gomock.Any(), orderUID).Return(nil, errors.New("boom"))

	res, err := as.getOrder(orderUID)
	if err == nil {
		t.Fatalf("expected error, got nil (res=%#v)", res)
	}
}
