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

### Установка

1. Клонируйте репозиторий:
```bash
git clone <repository-url>
cd order-service
```

2. Установите зависимости:
```bash
go mod download
```

3. Настройте переменные окружения:
```bash
cp .env.example .env
# Отредактируйте .env файл под вашу конфигурацию
```

4. Запустите зависимости (через Docker Compose):
```bash
docker-compose up -d
```

5. Запустите приложение:
```bash
go run cmd/app/main.go
```

### Переменные окружения

Создайте файл `.env` в корне проекта со следующими переменными:

```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=app_user
DB_PASSWORD=app_password
DB_NAME=app_db
DB_SSL_MODE=disable

# Redis
REDIS_HOST=localhost:6379
REDIS_PASSWORD=my_very_secure_password
REDIS_DB=0

# Kafka
KAFKA_BROKER=localhost:9092
KAFKA_TOPIC=orders
KAFKA_GROUP_ID=order-service-group

# HTTP Servers
API_SERVER_ADDR=:8081
WEB_SERVER_ADDR=:8080

# Logging
LOG_LEVEL=info
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

## Развертывание с Docker

Для развертывания с помощью Docker:

1. Соберите образ:
```bash
docker build -t order-service .
```

2. Запустите контейнер:
```bash
docker run -p 8080:8080 -p 8081:8081 --env-file .env order-service
```

Или используйте Docker Compose для полного развертывания со всеми зависимостями:

```bash
docker-compose up
```

## Мониторинг здоровья

Сервис включает эндпоинты для проверки здоровья:

- Health check: `/health`
- Готовность: `/ready`
- Статус: `/status`

## Логирование

Приложение использует структурированное логирование с различными уровнями детализации (debug, info, warn, error). Уровень логирования можно настроить через переменную окружения `LOG_LEVEL`.

## Разработка

### Добавление новых зависимостей

```bash
go get <package-name>
```

### Тестирование

```bash
go test ./...
```

### Форматирование кода

```bash
go fmt ./...
```

### Проверка стиля кода

```bash
go vet ./...
```
