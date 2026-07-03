CREATE TABLE rates (
    id SERIAL PRIMARY KEY,
    base_currency VARCHAR(10) NOT NULL,
    target_currency VARCHAR(10) NOT NULL,
    rate DECIMAL(10,6) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);