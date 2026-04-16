-- +goose Up
CREATE SCHEMA IF NOT EXISTS commerce;

CREATE TABLE
    commerce.orders (
        id UUID PRIMARY KEY,
        tenant_id UUID NOT NULL REFERENCES iam.tenants (id),
        user_id UUID NOT NULL REFERENCES auth.users (id),
        code TEXT NOT NULL,
        delivery_status VARCHAR(20) NOT NULL,
        fulfillment_status VARCHAR(20) NOT NULL,
        notes TEXT,
        total_price BIGINT NOT NULL,
        payment_status VARCHAR(20) NOT NULL,
        external_payment_id VARCHAR(100),
        paid_at TIMESTAMPTZ,
        created_at TIMESTAMPTZ,
        updated_at TIMESTAMPTZ,
        deleted_at TIMESTAMPTZ,
        items JSONB NOT NULL,
        CONSTRAINT uq_orders_code UNIQUE (code),
        CONSTRAINT uq_orders_external_payment_id UNIQUE (external_payment_id)
    );

COMMENT ON COLUMN commerce.orders.total_price IS 'Total order value in cents.';

COMMENT ON COLUMN commerce.orders.fulfillment_status IS 'Preparation lifecycle.';

COMMENT ON COLUMN commerce.orders.delivery_status IS 'Last-mile delivery lifecycle.';

COMMENT ON COLUMN commerce.orders.payment_status IS 'Financial lifecycle.';

COMMENT ON COLUMN commerce.orders.external_payment_id IS 'Transaction ID from the external payment gateway.';

COMMENT ON COLUMN commerce.orders.paid_at IS 'Timestamp of payment confirmation.';

COMMENT ON COLUMN commerce.orders.items IS 'JSON array of order items with quantity and unit price at the time of purchase.';

COMMENT ON COLUMN commerce.orders.notes IS 'Customer notes for preparation or delivery.';

CREATE INDEX idx_orders_deleted_at ON commerce.orders (deleted_at);

-- +goose Down
DROP TABLE IF EXISTS commerce.orders;

DROP SCHEMA IF EXISTS commerce;