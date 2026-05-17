-- +goose Up
CREATE SCHEMA IF NOT EXISTS payment;

CREATE TABLE
    payment.stripe_customers (
        user_id UUID PRIMARY KEY REFERENCES auth.users (id),
        stripe_customer_id VARCHAR(255) NOT NULL UNIQUE,
        stripe_payment_method_id VARCHAR(255),
        has_payment_method BOOLEAN NOT NULL DEFAULT FALSE,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
    );

-- +goose Down
DROP TABLE IF EXISTS payment.stripe_customers;

DROP SCHEMA IF EXISTS payment;
