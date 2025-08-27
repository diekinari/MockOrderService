# Order Processing Microservice

Микросервис для обработки заказов с использованием Go, PostgreSQL, Redis и Kafka.

## Описание

Сервис предназначен для получения данных о заказах из Kafka, их сохранения в PostgreSQL, кэширования в Redis и предоставления API для доступа к данным заказов через HTTP интерфейс.

## Архитектура

- **Язык**: Go
- **База данных**: PostgreSQL для хранения данных заказов
- **Кэш**: Redis для быстрого доступа к данным
- **Очередь сообщений**: Kafka для получения данных о заказах
- **API**: HTTP/JSON API для доступа к данным
- **Веб-интерфейс**: Простой HTML/JS интерфейс для просмотра заказов

## Функциональность

- Получение сообщений о заказах из Kafka
- Валидация и сохранение заказов в PostgreSQL
- Кэширование заказов в Redis для быстрого доступа
- HTTP API для получения информации о заказах по ID
- Веб-интерфейс для просмотра заказов
- Мониторинг здоровья компонентов системы
- Graceful shutdown при получении сигналов завершения

## Быстрый старт

### Предварительные требования

- Go 1.18+
- PostgreSQL 13+
- Redis 6+
- Kafka 2.8+
- Docker и Docker Compose (опционально)

### Переменные окружения

Создайте файл `.env` в корне проекта со следующими переменными:

```env
DB_HOST:localhost
DB_PORT:5433
DB_USER:app_user
DB_PASSWORD:app_password
DB_NAME:app_db
DB_SSL_MODE:disable

KAFKA_BROKER:localhost:9092
KAFKA_TOPIC:test-topic
KAFKA_GROUP_ID:demo-group

REDIS_HOST:127.0.0.1:6379
REDIS_PASSWORD:my_very_secure_password
```

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


## Логирование

Приложение использует структурированное логирование с различными уровнями детализации (debug, info, warn, error). 

