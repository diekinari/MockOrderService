package monitoring

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	kafkaMessagesProduced = prometheus.NewCounterVec(
		prometheus.CounterOpts{Name: "kafka_messages_produced_total", Help: "Produced messages"},
		[]string{"topic", "success"},
	)
	kafkaMessagesConsumed = prometheus.NewCounterVec(
		prometheus.CounterOpts{Name: "kafka_messages_consumed_total", Help: "Consumed messages"},
		[]string{"topic", "success"},
	)
)

func init() {
	Registry.MustRegister(kafkaMessagesProduced, kafkaMessagesConsumed)
}

func RecordKafkaMessagesProduced(topic string, success bool) {
	kafkaMessagesProduced.WithLabelValues(topic, strconv.FormatBool(success)).Inc()
}

func RecordKafkaMessagesConsumed(topic string, success bool) {
	kafkaMessagesConsumed.WithLabelValues(topic, strconv.FormatBool(success)).Inc()
}
