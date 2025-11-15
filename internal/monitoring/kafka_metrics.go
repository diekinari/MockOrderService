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
	kafkaMessagesSentToDQL = prometheus.NewCounterVec(
		prometheus.CounterOpts{Name: "kafka_messages_sent_to_dql_total", Help: "Messages sent to DQL"},
		[]string{"topic", "error_info"},
	)
	kafkaRetries = prometheus.NewCounterVec(
		prometheus.CounterOpts{Name: "kafka_retries_total", Help: "Retries"},
		[]string{"topic", "attempt"},
	)
)

func init() {
	Registry.MustRegister(kafkaMessagesProduced, kafkaMessagesConsumed, kafkaMessagesSentToDQL, kafkaRetries)
}

func RecordKafkaMessagesProduced(topic string, success bool) {
	kafkaMessagesProduced.WithLabelValues(topic, strconv.FormatBool(success)).Inc()
}

func RecordKafkaMessagesConsumed(topic string, success bool) {
	kafkaMessagesConsumed.WithLabelValues(topic, strconv.FormatBool(success)).Inc()
}

func RecordKafkaMessagesSentToDQL(topic string, error_info string) {
	kafkaMessagesSentToDQL.WithLabelValues(topic, error_info).Inc()
}

func RecordKafkaRetries(topic string, attempt int) {
	kafkaRetries.WithLabelValues(topic, strconv.Itoa(attempt)).Inc()
}