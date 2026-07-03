CREATE TABLE Orders (
    product_id UUID PRIMARY KEY,
    customer_id UUID not null  REFERENCES users(id) ON DELETE CASCADE,
    status VARCHAR(50) not null ,
    AMOUNT DECIMAL(10, 2) not null ,
    created_at TIMESTAMP NOT NULL DEFAULT NOW() ,
    updated_at TIMESTAMP NOT NULL DEFAULT NOW() ,
);


CREATE INDEX idx_orders_customer_id ON Orders(customer_id);





