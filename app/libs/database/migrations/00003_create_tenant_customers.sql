-- +goose Up
CREATE TABLE
    iam.tenant_customers (
        tenant_id     UUID        NOT NULL REFERENCES iam.tenants (id),
        user_id       UUID        NOT NULL REFERENCES auth.users (id),
        registered_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        PRIMARY KEY (tenant_id, user_id)
    );

COMMENT ON TABLE iam.tenant_customers IS 'Business relationship: customers registered to a tenant to place orders. Distinct from iam.tenant_members, which controls access authorization.';

-- +goose Down
DROP TABLE IF EXISTS iam.tenant_customers;
