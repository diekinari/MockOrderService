package http

import (
	"MockOrderService/internal/domain/model"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type OrderRepository interface {
	GetOrderByOrderUID(ctx context.Context, orderUID string) (*model.Order, error)
}

type CacheRepository interface {
	GetOrder(ctx context.Context, orderUID string) (*model.Order, error)
}

type ApiServer struct {
	sugar     *zap.SugaredLogger
	ctx       context.Context
	orderRepo OrderRepository
	cacheRepo CacheRepository
	server    *http.Server
}

type apiError struct {
	Error string `json:"Error"`
}

func NewApiServer(sugar *zap.SugaredLogger, ctx context.Context, orderRepo OrderRepository, cacheRepo CacheRepository) *ApiServer {
	return &ApiServer{
		sugar:     sugar,
		ctx:       ctx,
		orderRepo: orderRepo,
		cacheRepo: cacheRepo,
	}
}

// StartApiServer starts api server in a separate goroutine.
// handles order requests and shuts down gracefully on context cancel.
// If server fails to start or shutdown, logs error.
func (as *ApiServer) StartApiServer() error {
	r := mux.NewRouter()
	r.HandleFunc("/api/order/{orderUID}", as.handleOrder)

	srv := &http.Server{
		Addr:    ":8081",
		Handler: r,
	}

	go func() {
		<-as.ctx.Done()
		as.sugar.Infow("shutting down api server...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			as.sugar.Errorw("server shutdown failed", "error", err)
		}
	}()
	as.server = srv
	as.sugar.Infow("started api server at :8081")
	err := srv.ListenAndServe()
	if err != nil {
		as.sugar.Errorw("server stopped", "error", err)
		return err
	}
	return nil

}

func (as *ApiServer) Shutdown(ctx context.Context) error {
	if as.server != nil {
		return as.server.Shutdown(ctx)
	}
	return nil
}

func (as *ApiServer) handleOrder(w http.ResponseWriter, r *http.Request) {
	orderUID := mux.Vars(r)["orderUID"]
	if orderUID == "" {
		as.sugar.Infow("empty orderUID")

		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(&apiError{
			Error: "empty orderUID",
		}); err != nil {
			as.sugar.Errorw("couldn't encode error", "orderUID", orderUID, "error", err)
		}
		return
	}
	order, err := as.getOrder(orderUID)
	if err != nil {
		as.sugar.Errorw("couldn't get order", "orderUID", orderUID, "error", err)
		if errors.Is(err, pgx.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			if err := json.NewEncoder(w).Encode(&apiError{
				Error: "no order found",
			}); err != nil {
				as.sugar.Errorw("couldn't encode error", "orderUID", orderUID, "error", err)
			}
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			if err := json.NewEncoder(w).Encode(&apiError{
				Error: "couldn't get order",
			}); err != nil {
				as.sugar.Errorw("couldn't encode error", "orderUID", orderUID, "error", err)
			}
			return
		}

	}
	if err = json.NewEncoder(w).Encode(order); err != nil {
		as.sugar.Errorw("couldn't encode order", "orderUID", orderUID, "error", err)
	}
}

func (as *ApiServer) getOrder(orderUID string) (*model.Order, error) {
	// try redis first
	val, err := as.cacheRepo.GetOrder(as.ctx, orderUID)
	if err != nil {
		// key no found
		if errors.Is(err, redis.Nil) {
			// try from db
			order, err := as.orderRepo.GetOrderByOrderUID(as.ctx, orderUID)
			if err != nil {
				return nil, err
			}
			return order, nil
		} else {
			return nil, err
		}
	}
	return val, nil

}
