-- +goose Up
ALTER TABLE subscription.subscriptions
    ADD COLUMN IF NOT EXISTS checkout_key UUID;

COMMENT ON COLUMN subscription.subscriptions.checkout_key IS 'Client-generated idempotency key (UUIDv7). Used with a partial unique index on (user_id, checkout_key, product_id) to deduplicate concurrent checkout attempts for the same cart.';

CREATE UNIQUE INDEX IF NOT EXISTS idx_subscriptions_checkout_key
    ON subscription.subscriptions (user_id, checkout_key, product_id)
    WHERE checkout_key IS NOT NULL;

-- +goose Down
DROP INDEX IF EXISTS subscription.idx_subscriptions_checkout_key;

ALTER TABLE subscription.subscriptions
    DROP COLUMN IF EXISTS checkout_key;
