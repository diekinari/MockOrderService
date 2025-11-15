# Order Processing Microservice

Микросервис для обработки заказов с использованием Go, PostgreSQL, Redis и Kafka.

## Описание

Сервис предназначен для получения данных о заказах из Kafka, их сохранения в PostgreSQL, кэширования в Redis и предоставления API для доступа к данным заказов через HTTP интерфейс.

## Архитектура

- **Язык**: Go
- **База данных**: PostgreSQL для хранения данных заказов
- **Кэш**: Redis для быстрого доступа к данным
- **Очередь сообщений**: Kafka для получения данных о заказах
- **DLQ**: Dead Letter Queue для обработки проблемных сообщений
- **API**: HTTP/JSON API для доступа к данным
- **Веб-интерфейс**: Простой HTML/JS интерфейс для просмотра заказов
- **Мониторинг**: Prometheus для метрик, Grafana для визуализации
- **Трейсинг**: Jaeger для распределенной трассировки

## Функциональность

- Получение сообщений о заказах из Kafka
- Валидация и сохранение заказов в PostgreSQL
- Кэширование заказов в Redis для быстрого доступа
- HTTP API для получения информации о заказах по ID
- Веб-интерфейс для просмотра заказов
- **Retry механизм** с настраиваемым количеством попыток и задержкой
- **Dead Letter Queue (DLQ)** для сообщений, которые не удалось обработать после исчерпания попыток
- Мониторинг здоровья компонентов системы
- Метрики Prometheus для бизнес-логики, Kafka, БД и кэша
- Распределенная трассировка через OpenTelemetry и Jaeger
- Graceful shutdown при получении сигналов завершения

## Быстрый старт

### Предварительные требования

- Go 1.18+
- Docker и Docker Compose (рекомендуется)

### Запуск через Docker Compose

1. **Запустите инфраструктуру** (PostgreSQL, Kafka, Redis, Prometheus, Grafana, Jaeger):

```bash
cd deployments
docker-compose up -d
```

Это запустит все необходимые сервисы:
- PostgreSQL на порту 5433
- Kafka на порту 9092
- Redis на порту 6379
- Prometheus на порту 9090
- Grafana на порту 3000 (логин/пароль: `admin/admin`)
- Jaeger UI на порту 16686, OTLP endpoint на 4318

2. **Создайте файл `.env`** в корне проекта (см. раздел "Переменные окружения")

3. **Установите зависимости Go**:

```bash
go mod download
```

4. **Запустите приложение**:

```bash
go run cmd/app/main.go
```

Или соберите бинарник:

```bash
go build -o main cmd/app/main.go
./main
```

### Полезные команды

Остановить инфраструктуру:
```bash
cd deployments
docker-compose down
```

Просмотр логов:
```bash
cd deployments
docker-compose logs -f
```

Проверка статуса контейнеров:
```bash
cd deployments
docker-compose ps
```

### Доступные сервисы

После запуска доступны следующие интерфейсы:

- **API сервер**: http://localhost:8081
- **Веб-интерфейс**: http://localhost:8082
- **Prometheus**: http://localhost:9090
- **Grafana**: http://localhost:3000 (логин/пароль: `admin/admin`)
- **Jaeger UI**: http://localhost:16686

### Переменные окружения

Создайте файл `.env` в корне проекта со следующими переменными:

```env
# Database
DB_HOST=localhost
DB_PORT=5433
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=postgres
DB_SSL_MODE=disable

# Kafka
KAFKA_BROKER=localhost:9092
KAFKA_TOPIC=test-topic
KAFKA_GROUP_ID=demo-group

# Redis
REDIS_HOST=127.0.0.1:6379
REDIS_PASSWORD=

# Retry configuration (опционально, дефолты: MAX_RETRIES=3, RETRY_BACKOFF_MS=1000)
MAX_RETRIES=3
RETRY_BACKOFF_MS=1000

# Jaeger (опционально, дефолт: localhost:4318)
JAEGER_ENDPOINT=localhost:4318
```

**Примечание**: DLQ топик создается автоматически как `{KAFKA_TOPIC}-dlq` (например, `test-topic-dlq`).

## API Endpoints

### Получить информацию о заказе

```
GET /order/{order_uid}
```

Пример ответа:
```json
{
  "order_uid": "b563feb7b2b84b6test",
  "track_number": "WBILMTESTTRACK",
  "entry": "WBIL",
  "delivery": {
    "name": "Test Testov",
    "phone": "+9720000000",
    "zip": "2639809",
    "city": "Kiryat Mozkin",
    "address": "Ploshad Mira 15",
    "region": "Kraiot",
    "email": "test@gmail.com"
  },
  "payment": {
    "transaction": "b563feb7b2b84b6test",
    "request_id": "",
    "currency": "USD",
    "provider": "wbpay",
    "amount": 1817,
    "payment_dt": 1637907727,
    "bank": "alpha",
    "delivery_cost": 1500,
    "goods_total": 317,
    "custom_fee": 0
  },
  "items": [
    {
      "chrt_id": 9934930,
      "track_number": "WBILMTESTTRACK",
      "price": 453,
      "rid": "ab4219087a764ae0btest",
      "name": "Mascaras",
      "sale": 30,
      "size": "0",
      "total_price": 317,
      "nm_id": 2389212,
      "brand": "Vivienne Sabo",
      "status": 202
    }
  ],
  "locale": "en",
  "internal_signature": "",
  "customer_id": "test",
  "delivery_service": "meest",
  "shardkey": "9",
  "sm_id": 99,
  "date_created": "2021-11-26T06:22:19Z",
  "oof_shard": "1"
}
```

## Веб-интерфейс

Веб-интерфейс доступен по адресу http://localhost:8082 после запуска приложения. Он позволяет:

1. Ввести ID заказа
2. Получить информацию о заказе
3. Просмотреть детали заказа в удобном формате


## Dead Letter Queue (DLQ)

Сервис поддерживает механизм Dead Letter Queue для обработки сообщений, которые не удалось обработать после исчерпания попыток повтора.

### Механизм работы

1. **Retry механизм**: При ошибке обработки сообщение повторяется до `MAX_RETRIES` раз с задержкой `RETRY_BACKOFF_MS` миллисекунд между попытками.

2. **Отправка в DLQ**: Если после всех попыток обработка не удалась, сообщение отправляется в DLQ топик `{KAFKA_TOPIC}-dlq`.

3. **Метаданные в DLQ**: Каждое сообщение в DLQ содержит заголовки:
   - `x-error` - описание ошибки
   - `x-retry-count` - количество выполненных попыток
   - `x-original-topic` - исходный топик
   - `x-original-partition` - исходная партиция
   - `x-original-offset` - исходный offset
   - `x-original-key` - исходный ключ сообщения
   - `x-timestamp` - время отправки в DLQ

### Типы ошибок

- **Unmarshal ошибки**: Отправляются в DLQ сразу (постоянная ошибка формата)
- **Service ошибки**: Повторы → DLQ при исчерпании попыток
- **Commit ошибки**: Повторы → DLQ при исчерпании попыток
- **Validation ошибки**: Пропускаются (не отправляются в DLQ)

## Мониторинг

### Метрики Prometheus

Сервис экспортирует следующие метрики:

- `kafka_messages_produced_total` - количество отправленных сообщений в Kafka
- `kafka_messages_consumed_total` - количество полученных сообщений из Kafka
- `kafka_messages_sent_to_dql_total` - количество сообщений, отправленных в DLQ
- `kafka_retries_total` - количество попыток повтора
- `orders_created_total` - количество успешно созданных заказов
- `orders_failed_total` - количество неудачных попыток обработки заказов
- `order_process_duration_seconds` - время обработки заказа
- Метрики БД и кэша

Метрики доступны по адресу: http://localhost:8081/metrics

### Grafana Dashboard

В Grafana доступен готовый дашборд с визуализацией метрик сервиса. Дашборд автоматически подключается при запуске через docker-compose.

### Трейсинг

Сервис использует OpenTelemetry для распределенной трассировки. Трейсы отправляются в Jaeger и доступны через Jaeger UI.

## Логирование

Приложение использует структурированное логирование с различными уровнями детализации (debug, info, warn, error).
