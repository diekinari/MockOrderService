package monitoring

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	DBQueryDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "db_query_duration_seconds",
			Help:    "DB query duration",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"op", "table", "success"},
	)
)

func init() {
	Registry.MustRegister(DBQueryDuration)
}

func RecordDBQueryDuration(operation, table string, success bool, duration time.Duration) {
	DBQueryDuration.WithLabelValues(operation, table, strconv.FormatBool(success)).Observe(duration.Seconds())
}
