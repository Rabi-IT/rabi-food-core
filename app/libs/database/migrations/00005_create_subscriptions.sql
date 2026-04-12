-- +goose Up
CREATE SCHEMA IF NOT EXISTS subscription;

CREATE TABLE
    subscription.subscriptions (
        id UUID PRIMARY KEY,
        root_subscription_id UUID REFERENCES subscription.subscriptions (id),
        tenant_id UUID NOT NULL REFERENCES iam.tenants (id),
        user_id UUID NOT NULL REFERENCES iam.users (id),
        delivery_days JSONB,
        delivery_weekdays_mask SMALLINT NOT NULL DEFAULT 0,
        items JSONB,
        notes TEXT,
        total_cycles BIGINT NOT NULL,
        remaining_cycles BIGINT NOT NULL,
        cycle_discount BIGINT NOT NULL DEFAULT 0,
        cutoff_offset_minutes INTEGER NOT NULL DEFAULT 0,
        auto_renew BOOLEAN NOT NULL DEFAULT TRUE,
        max_attempts_per_order SMALLINT NOT NULL DEFAULT 1,
        items_total BIGINT NOT NULL DEFAULT 0,
        items_discount BIGINT NOT NULL DEFAULT 0,
        payment_amount BIGINT NOT NULL,
        payment_status VARCHAR(20) NOT NULL,
        external_payment_id VARCHAR(100) NOT NULL,
        status VARCHAR(20),
        created_at TIMESTAMPTZ,
        updated_at TIMESTAMPTZ,
        CONSTRAINT uq_subscriptions_external_payment_id UNIQUE (external_payment_id)
    );

CREATE INDEX idx_subscriptions_root_subscription_id ON subscription.subscriptions (root_subscription_id);

CREATE INDEX idx_subscriptions_tenant_id ON subscription.subscriptions (tenant_id);

CREATE INDEX idx_subscriptions_user_id ON subscription.subscriptions (user_id);

CREATE INDEX idx_subscriptions_delivery_weekdays_mask ON subscription.subscriptions (delivery_weekdays_mask);

CREATE INDEX idx_subscriptions_remaining_cycles ON subscription.subscriptions (remaining_cycles);

CREATE INDEX idx_subscriptions_payment_status ON subscription.subscriptions (payment_status);

CREATE INDEX idx_subscriptions_status ON subscription.subscriptions (status);

CREATE TABLE
    subscription.subscription_configs (
        tenant_id UUID PRIMARY KEY NOT NULL REFERENCES iam.tenants (id),
        is_open BOOLEAN NOT NULL,
        max_attempts_per_order SMALLINT,
        discount_rules JSONB,
        cutoff_offset_minutes INTEGER,
        updated_at TIMESTAMPTZ
    );

CREATE TABLE
    subscription.subscription_deliveries (
        id UUID PRIMARY KEY,
        subscription_id UUID NOT NULL REFERENCES subscription.subscriptions (id),
        scheduled_date TIMESTAMPTZ,
        start_hour SMALLINT NOT NULL,
        end_hour SMALLINT NOT NULL,
        cutoff_at TIMESTAMPTZ NOT NULL,
        status VARCHAR(20),
        delivery_attempts SMALLINT,
        max_delivery_attempts SMALLINT,
        created_at TIMESTAMPTZ,
        updated_at TIMESTAMPTZ,
        CONSTRAINT subscription_delivery_unique UNIQUE (subscription_id, scheduled_date)
    );

CREATE INDEX idx_subscription_deliveries_subscription_id ON subscription.subscription_deliveries (subscription_id);

CREATE INDEX idx_subscription_deliveries_scheduled_date ON subscription.subscription_deliveries (scheduled_date);

CREATE INDEX idx_subscription_deliveries_status ON subscription.subscription_deliveries (status);

-- +goose Down
DROP TABLE IF EXISTS subscription.subscription_deliveries;

DROP TABLE IF EXISTS subscription.subscription_configs;

DROP TABLE IF EXISTS subscription.subscriptions;

DROP SCHEMA IF EXISTS subscription;
