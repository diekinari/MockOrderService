package main

import (
	"MockOrderService/config"
	httpdelivery "MockOrderService/internal/delivery/http"
	"MockOrderService/internal/delivery/kafka"
	kafkaInfra "MockOrderService/internal/infra/kafka"
	"MockOrderService/internal/infra/postgres"
	"MockOrderService/internal/infra/redis"
	"MockOrderService/internal/logger"
	"MockOrderService/internal/monitoring"
	postgresRepo "MockOrderService/internal/repository/postgres"
	redisRepo "MockOrderService/internal/repository/redis"
	"MockOrderService/internal/service"
	"context"
	"errors"
	"fmt"
	"log"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	// Initiating logger
	sugar, err := logger.NewLogger()
	if err != nil {
		log.Fatalf("cannot init logger %v", err)
		return
	}
	//defer sugar.Sync()

	cfg, err := config.Load()
	if err != nil {
		sugar.Fatalw("cannot load config", "error", err)
		return
	}

	sugar.Infow("starting application", "version", "1.0.0")
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	pgClient, err := postgres.NewClient(cfg, ctx)
	if err != nil {
		sugar.Fatalw("failed to initialize PostgreSQL client", "error", err)
		return
	}
	defer pgClient.Close()
	sugar.Infow("database connection is established")

	redisClient, err := redis.NewClient(cfg, ctx)
	if err != nil {
		sugar.Fatalw("failed to initialize Redis client", "error", err)
		return
	}
	//defer redisClient.Close()

	kafkaClient, err := kafkaInfra.NewClient(cfg.KafkaBroker, cfg.KafkaGroupId, cfg.KafkaTopic)
	if err != nil {
		sugar.Fatalw("failed to initialize Kafka client", "error", err)
		return
	}
	//defer kafkaClient.Close()

	orderRepo := postgresRepo.NewOrderRepository(pgClient.Pool)
	cacheRepo := redisRepo.NewCacheRepository(redisClient.Client)

	orderService := service.NewOrderService(sugar, orderRepo, cacheRepo)
	go orderService.HeatUpCache(ctx)

	kafkaProducer := kafka.NewProducer(kafkaClient, sugar)
	go kafkaProducer.Start(stop)

	kafkaConsumer := kafka.NewConsumer(kafkaClient, orderService, sugar)
	go kafkaConsumer.Start(ctx, stop)

	healthChecker := monitoring.NewHealthChecker(pgClient, redisClient, 10*time.Second, sugar, stop)
	go healthChecker.Start(ctx)

	serverErrors := make(chan error, 2)

	// api for frontend
	apiServer := httpdelivery.NewApiServer(sugar, ctx, orderRepo, cacheRepo)

	webServer := &httpdelivery.WebServer{}

	go func() {
		if err := apiServer.StartApiServer(); err != nil {
			serverErrors <- fmt.Errorf("api server error: %w", err)
		}
	}()

	// web-server with dashboard
	go func() {
		if err := webServer.StartWebServer(sugar); err != nil {
			serverErrors <- fmt.Errorf("web server error: %w", err)
		}
	}()

	select {
	case <-ctx.Done():
		sugar.Info("Received shutdown signal, initiating graceful shutdown")
	case err := <-serverErrors:
		sugar.Errorw("Server error, initiating shutdown", "error", err)
		stop()
	}

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var errs []error

	// Останавливаем серверы
	if err := apiServer.Shutdown(shutdownCtx); err != nil {
		//sugar.Errorw("Failed to shutdown API server gracefully", "error", err)
		errs = append(errs, fmt.Errorf("API_Server: %w", err))
	}

	if err := webServer.Shutdown(shutdownCtx); err != nil {
		//sugar.Errorw("Failed to shutdown Web server gracefully", "error", err)
		errs = append(errs, fmt.Errorf("Web_Server: %w", err))
	}

	// Закрываем подключения
	if err := redisClient.Close(); err != nil {
		errs = append(errs, fmt.Errorf("redis: %w", err))
	}

	if err := kafkaClient.Close(); err != nil {
		errs = append(errs, fmt.Errorf("kafka: %w", err))
	}

	if err := sugar.Sync(); err != nil {
		errs = append(errs, fmt.Errorf("logger: %w", err))
	}

	if len(errs) > 0 {
		sugar.Errorw("errors during shutdown", "error", errors.Join(errs...))
	}

	sugar.Info("Application shutdown complete")

}
