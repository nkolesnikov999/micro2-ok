-- +goose Up
CREATE TABLE orders (
    order_uuid VARCHAR(36) PRIMARY KEY,
    user_uuid VARCHAR(36) NOT NULL,
    part_uuids TEXT[] NOT NULL DEFAULT '{}',
    total_price DECIMAL(10,2) NOT NULL,
    transaction_uuid VARCHAR(36),
    payment_method VARCHAR(50),
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING'
);

-- +goose Down
DROP TABLE orders;