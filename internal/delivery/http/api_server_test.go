package http

import (
	"MockOrderService/internal/mocks"
	"context"
	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestApiServer_StartApiServer(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	orderRepo := mocks.NewMockOrderRepository(ctrl)

	cacheRepo := mocks.NewMockCacheRepository(ctrl)

	sugar := zap.NewExample().Sugar()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	apiServer := NewApiServer(sugar, ctx, orderRepo, cacheRepo)
	go func() {
		err := apiServer.StartApiServer()
		if err != nil {
			t.Error(err)
		}
	}()

	request := httptest.NewRequest(http.MethodGet, "/", nil)
}
