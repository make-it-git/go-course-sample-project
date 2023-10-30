-- Create orders

CREATE TABLE IF NOT EXISTS orders (
    id uuid PRIMARY KEY,
    created_at timestamp NOT NULL,
    completed_at timestamp NULL,
    pickup_location jsonb NOT NULL,
    dropoff_location jsonb NOT NULL,
    total_price integer NOT NULL,
    user_id integer NOT NULL,
    idempotency_key uuid NOT NULL UNIQUE
);

CREATE INDEX idx_orders_user_id ON orders (user_id);
---- create above / drop below ----

drop table orders;
