CREATE TABLE IF NOT EXISTS Orders (
    id UUID PRIMARY KEY,
    customer_id UUID not null,
    status VARCHAR(50) not null ,
    AMOUNT DECIMAL(10, 2) not null ,
    created_at TIMESTAMP NOT NULL DEFAULT NOW() ,
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);


CREATE INDEX idx_orders_customer_id ON Orders(customer_id);



CREATE TABLE IF NOT EXISTS order_status_history (
    id UUID PRIMARY KEY,
    order_id UUID NOT NULL REFERENCES Orders(id) ON DELETE CASCADE,
    status VARCHAR(50) NOT NULL,
    comment VARCHAR(512),
    created_at TIMESTAMP NOT NULL DEFAULT NOW
);

CREATE INDEX idx_order_status_history_order_id ON order_status_history(order_id);
