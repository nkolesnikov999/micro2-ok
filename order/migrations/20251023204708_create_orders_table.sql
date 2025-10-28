-- +goose Up
CREATE TABLE orders (
    order_uuid UUID PRIMARY KEY,
    user_uuid UUID NOT NULL,
    total_price DECIMAL(10,2) NOT NULL,
    transaction_uuid UUID,
    payment_method TEXT,
    status TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE TABLE order_parts (
    order_uuid UUID NOT NULL,
    part_uuid UUID NOT NULL,
    quantity INTEGER NOT NULL DEFAULT 1,
    PRIMARY KEY (order_uuid, part_uuid),
    FOREIGN KEY (order_uuid) REFERENCES orders(order_uuid) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE order_parts;
DROP TABLE orders;