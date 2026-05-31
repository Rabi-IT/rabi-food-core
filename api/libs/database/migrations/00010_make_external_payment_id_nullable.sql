-- +goose Up
-- external_payment_id is now set after a successful charge (via webhook), not during subscription creation.
ALTER TABLE subscription.subscriptions
    ALTER COLUMN external_payment_id DROP NOT NULL;

-- +goose Down
-- Restoring NOT NULL may fail if NULLs exist; handle manually if needed.
ALTER TABLE subscription.subscriptions
    ALTER COLUMN external_payment_id SET NOT NULL;
