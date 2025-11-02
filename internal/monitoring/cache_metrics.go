package monitoring

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	cacheOpDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "cache_op_duration_seconds",
			Help:    "Cache operation duration",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"op", "success"},
	)
	cacheHitsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "cache_hits_total",
			Help: "Total number of cache hits",
		},
	)
	cacheMissesTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "cache_misses_total",
			Help: "Total number of cache misses",
		},
	)
)

func init() {
	Registry.MustRegister(cacheOpDuration, cacheHitsTotal, cacheMissesTotal)
}

func RecordCacheOpDuration(op string, success bool, duration time.Duration) {
	cacheOpDuration.WithLabelValues(op, strconv.FormatBool(success)).Observe(duration.Seconds())
}

func RecordCacheHit() {
	cacheHitsTotal.Inc()
}

func RecordCacheMiss() {
	cacheMissesTotal.Inc()
}
