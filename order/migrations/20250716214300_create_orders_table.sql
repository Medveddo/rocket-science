-- +goose Up
CREATE TABLE orders (
    order_uuid UUID PRIMARY KEY,
    user_uuid UUID NOT NULL,
    part_uuids TEXT[] NOT NULL,
    total_price NUMERIC(15, 2) NOT NULL,
    transaction_uuid UUID,
    payment_method TEXT,
    status TEXT NOT NULL
);

-- +goose Down
DROP TABLE orders;
