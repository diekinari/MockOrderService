package monitoring

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"
	"go.opentelemetry.io/otel/trace"
)

var (
	// TracerProvider - глобальный провайдер трейсов
	TracerProvider *sdktrace.TracerProvider
	// Tracer - глобальный трейсер для создания спанов
	Tracer trace.Tracer
)

// InitTracing инициализирует OpenTelemetry трейсинг с Jaeger.
func InitTracing(serviceName, jaegerEndpoint string, ctx context.Context) error {
	// Создаем ресурс с информацией о сервисе
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion("1.0.0"),
		),
	)
	if err != nil {
		return fmt.Errorf("failed to create resource: %w", err)
	}

	// Создаем OTLP HTTP экспортер для отправки трейсов в Jaeger
	// Jaeger поддерживает OTLP через HTTP endpoint
	exporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(jaegerEndpoint),
		otlptracehttp.WithInsecure(), // для dev окружения, в prod используйте TLS
	)
	if err != nil {
		return fmt.Errorf("failed to create OTLP exporter: %w", err)
	}

	// Создаем TracerProvider с батчингом
	// Батчинг позволяет отправлять несколько спанов одной пачкой для эффективности
	TracerProvider = sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()), // в prod можно использовать другие сэмплеры
	)

	// Устанавливаем глобальный TracerProvider
	otel.SetTracerProvider(TracerProvider)

	// Устанавливаем глобальный propagator для распространения trace context
	// Это позволяет передавать trace ID между сервисами через HTTP заголовки
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	// Создаем глобальный трейсер
	Tracer = otel.Tracer(serviceName)

	return nil
}

// ShutdownTracing корректно завершает работу трейсинга.
// Должен вызываться при завершении приложения для корректной отправки
// всех трейсов в Jaeger.
func ShutdownTracing(ctx context.Context) error {
	if TracerProvider == nil {
		return nil
	}

	// Создаем контекст с таймаутом для shutdown
	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Завершаем работу TracerProvider
	// Это отправит все оставшиеся трейсы в Jaeger
	if err := TracerProvider.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("failed to shutdown tracer provider: %w", err)
	}

	return nil
}

