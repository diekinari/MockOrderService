CREATE USER app_user WITH password 'app_password';
CREATE DATABASE app_db OWNER app_user;

-- Переключаемся на целевую базу
\connect app_db app_user

-- Основная таблица заказов
CREATE TABLE IF NOT EXISTS orders (
                                      order_uid       TEXT PRIMARY KEY,
                                      track_number    TEXT,
                                      entry           TEXT,
                                      locale          TEXT,
                                      internal_signature TEXT,
                                      customer_id     TEXT,
                                      delivery_service TEXT,
                                      shardkey        TEXT,
                                      sm_id           INTEGER,
                                      date_created    TIMESTAMPTZ,  -- ISO8601 like "2021-11-26T06:22:19Z"
                                      oof_shard       TEXT,
                                      created_at      TIMESTAMPTZ DEFAULT now()
    );

-- Таблица доставки (one-to-one с за
-- казом)
CREATE TABLE IF NOT EXISTS deliveries (
                                          id        SERIAL PRIMARY KEY,
                                          order_uid TEXT NOT NULL REFERENCES orders(order_uid) ON DELETE CASCADE,
    name      TEXT,
    phone     TEXT,
    zip       TEXT,
    city      TEXT,
    address   TEXT,
    region    TEXT,
    email     TEXT,
    created_at TIMESTAMPTZ DEFAULT now()
    );

-- Таблица платежа (one-to-one с заказом)
CREATE TABLE IF NOT EXISTS payments (
                                        id           SERIAL PRIMARY KEY,
                                        order_uid    TEXT NOT NULL REFERENCES orders(order_uid) ON DELETE CASCADE,
    transaction_id TEXT,
    request_id   TEXT,
    currency     TEXT,
    provider     TEXT,
    amount       BIGINT,      -- cents / integer amount as in example (1817)
    payment_dt   BIGINT,      -- epoch seconds as in example (1637907727)
    bank         TEXT,
    delivery_cost BIGINT,
    goods_total  BIGINT,
    custom_fee   BIGINT,
    created_at   TIMESTAMPTZ DEFAULT now()
    );

-- Таблица items (каждый товар в отдельной строке; order_uid связывает с orders)
CREATE TABLE IF NOT EXISTS items (
                                     id          SERIAL PRIMARY KEY,
                                     order_uid   TEXT NOT NULL REFERENCES orders(order_uid) ON DELETE CASCADE,
    chrt_id     BIGINT,
    track_number TEXT,
    price       BIGINT,
    rid         TEXT NOT NULL,
    name        TEXT,
    sale        INTEGER,
    size        TEXT,
    total_price BIGINT,
    nm_id       BIGINT,
    brand       TEXT,
    status      INTEGER,
    created_at  TIMESTAMPTZ DEFAULT now()
    );

-- Индексы для ускорения поиска
CREATE INDEX IF NOT EXISTS idx_orders_customer_id ON orders(customer_id);
CREATE INDEX IF NOT EXISTS idx_orders_date_created ON orders(date_created);
CREATE INDEX IF NOT EXISTS idx_orders_track_number ON orders(track_number);
CREATE INDEX IF NOT EXISTS idx_items_nm_id ON items(nm_id);
CREATE INDEX IF NOT EXISTS idx_items_chrt_id ON items(chrt_id);
CREATE INDEX IF NOT EXISTS idx_items_order_uid ON items(order_uid);


ALTER TABLE deliveries
    ADD CONSTRAINT deliveries_order_uid_key UNIQUE (order_uid);

ALTER TABLE payments
    ADD CONSTRAINT payments_order_uid_key UNIQUE (order_uid);

ALTER TABLE items
    ADD CONSTRAINT items_order_rid_key UNIQUE (order_uid, rid);

-- Доступы для пользователя приложения (app_user)
-- GRANT USAGE ON SCHEMA public TO app_user;
-- GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO app_user;
-- GRANT USAGE, SELECT ON SEQUENCE deliveries_id_seq TO app_user;
-- GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO app_user;
--
-- -- Чтобы будущие таблицы автоматически получали права app_user
-- ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO app_user;


-- === Конец блока ===

