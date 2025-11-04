package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	KafkaBroker  string
	KafkaTopic   string
	KafkaGroupId string

	RedisHost     string
	RedisPassword string

	JaegerEndpoint string
}

func Load() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	dbHost, err := getEnv("DB_HOST")
	if err != nil {
		return nil, err
	}
	dbPort, err := getEnv("DB_PORT")
	if err != nil {
		return nil, err
	}
	dbUser, err := getEnv("DB_USER")
	if err != nil {
		return nil, err
	}
	dbPass, err := getEnv("DB_PASSWORD")
	if err != nil {
		return nil, err
	}
	dbName, err := getEnv("DB_NAME")
	if err != nil {
		return nil, err
	}
	dbSSLMode, err := getEnv("DB_SSL_MODE")
	if err != nil {
		return nil, err
	}
	kafkaBroker, err := getEnv("KAFKA_BROKER")
	if err != nil {
		return nil, err
	}
	kafkaTopic, err := getEnv("KAFKA_TOPIC")
	if err != nil {
		return nil, err
	}
	kafkaGroupId, err := getEnv("KAFKA_GROUP_ID")
	if err != nil {
		return nil, err
	}
	redisHost, err := getEnv("REDIS_HOST")
	if err != nil {
		return nil, err
	}
	redisPass, err := getEnv("REDIS_PASSWORD")
	if err != nil {
		return nil, err
	}

	// Jaeger endpoint опциональный, если не указан - используем дефолтный
	jaegerEndpoint := os.Getenv("JAEGER_ENDPOINT")
	if jaegerEndpoint == "" {
		jaegerEndpoint = "localhost:4318" // дефолтный OTLP HTTP endpoint для Jaeger
	}

	config := &Config{
		DBHost:         dbHost,
		DBPort:         dbPort,
		DBUser:         dbUser,
		DBPassword:     dbPass,
		DBName:         dbName,
		DBSSLMode:      dbSSLMode,
		KafkaBroker:    kafkaBroker,
		KafkaTopic:     kafkaTopic,
		KafkaGroupId:   kafkaGroupId,
		RedisHost:      redisHost,
		RedisPassword:  redisPass,
		JaegerEndpoint: jaegerEndpoint,
	}

	return config, nil

}

func getEnv(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", errors.New("no variable found: " + key)
	}
	return value, nil
}
