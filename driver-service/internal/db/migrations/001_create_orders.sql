-- Create orders

CREATE TABLE IF NOT EXISTS orders (
    id uuid PRIMARY KEY,
    pickup_location jsonb NOT NULL,
    dropoff_location jsonb NOT NULL,
    last_active_location jsonb NULL,
    completed_at timestamp NULL,
    user_id integer NOT NULL,
    driver_id integer NULL
);

CREATE INDEX idx_orders_driver_id ON orders (driver_id);
---- create above / drop below ----

drop table orders;
