-- +goose Up
ALTER TABLE subscription.subscriptions
    ADD COLUMN IF NOT EXISTS paid_at TIMESTAMPTZ;

COMMENT ON COLUMN subscription.subscriptions.paid_at IS 'Timestamp when the payment was confirmed by the payment provider.';

-- +goose Down
ALTER TABLE subscription.subscriptions
    DROP COLUMN IF EXISTS paid_at;
