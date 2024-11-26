CREATE TABLE IF NOT EXISTS orders
(
    order_uid          VARCHAR PRIMARY KEY,
    track_number       TEXT,
    entry              TEXT,
    locale             TEXT,
    internal_signature TEXT,
    customer_id        TEXT,
    delivery_service   TEXT,
    shardkey           TEXT,
    sm_id              BIGINT,
    date_created       TIMESTAMP,
    oof_shard          TEXT
);

CREATE TABLE IF NOT EXISTS delivery
(
    id        SERIAL PRIMARY KEY,
    order_uid TEXT NOT NULL,
    name      TEXT NOT NULL,
    phone     TEXT NOT NULL,
    zip       TEXT NOT NULL,
    city      TEXT NOT NULL,
    address   TEXT NOT NULL,
    region    TEXT NOT NULL,
    email     TEXT NOT NULL,
    FOREIGN KEY (order_uid) REFERENCES orders (order_uid) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS payment
(   
    id            SERIAL PRIMARY KEY,
    order_uid     TEXT   NOT NULL,
    transaction   TEXT   NOT NULL,
    request_id    TEXT,
    currency      TEXT   NOT NULL,
    provider      TEXT   NOT NULL,
    amount        BIGINT NOT NULL,
    payment_dt    BIGINT NOT NULL,
    bank          TEXT   NOT NULL,
    delivery_cost BIGINT NOT NULL,
    goods_total   BIGINT NOT NULL,
    custom_fee    BIGINT,
    FOREIGN KEY (order_uid) REFERENCES orders (order_uid) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS item
(
    id           SERIAL PRIMARY KEY,
    order_uid    TEXT   NOT NULL,
    chrt_id      BIGINT NOT NULL,
    track_number TEXT   NOT NULL,
    price        BIGINT NOT NULL,
    rid          TEXT   NOT NULL,
    name         TEXT   NOT NULL,
    sale         BIGINT NOT NULL,
    size         TEXT   NOT NULL,
    total_price  BIGINT NOT NULL,
    nm_id        BIGINT NOT NULL,
    brand        TEXT   NOT NULL,
    status       BIGINT NOT NULL,
    FOREIGN KEY (order_uid) REFERENCES orders (order_uid) ON DELETE CASCADE
);

CREATE INDEX idx_order_uid ON orders(order_uid);
CREATE INDEX idx_delivery_order_uid ON delivery(order_uid);
CREATE INDEX idx_payment_order_uid ON payment(order_uid);
CREATE INDEX idx_item_order_uid ON item(order_uid);