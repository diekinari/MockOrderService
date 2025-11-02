package monitoring

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	OrdersCreatedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "orders_created_total",
			Help: "Total number of orders created",
		},
	)
	OrdersFailedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "orders_failed_total",
			Help: "Total number of orders failed",
		},
	)
	OrdersProcessDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "orders_process_duration_seconds",
			Help:    "Duration of order processing",
			Buckets: prometheus.DefBuckets, // Стандартные buckets для секунд
		},
	)
)

func init() {
	Registry.MustRegister(OrdersCreatedTotal, OrdersFailedTotal, OrdersProcessDuration)
}

func RecordOrderCreated() {
	OrdersCreatedTotal.Inc()
}

func RecordOrderFailed() {
	OrdersFailedTotal.Inc()
}

func RecordOrderProcessDuration(duration time.Duration) {
	OrdersProcessDuration.Observe(duration.Seconds())
}
