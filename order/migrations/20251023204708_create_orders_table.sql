-- +goose Up
CREATE TABLE orders (
    order_uuid UUID PRIMARY KEY,
    user_uuid UUID NOT NULL,
    part_uuids TEXT[] NOT NULL DEFAULT '{}',
    total_price DECIMAL(10,2) NOT NULL,
    transaction_uuid UUID,
    payment_method TEXT,
    status TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE orders;