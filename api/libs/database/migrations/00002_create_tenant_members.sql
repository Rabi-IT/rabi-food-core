-- +goose Up
CREATE TABLE
    iam.tenant_members (
        tenant_id UUID NOT NULL REFERENCES iam.tenants (id),
        user_id UUID NOT NULL REFERENCES auth.users (id),
        ROLE VARCHAR(20) NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        PRIMARY KEY (tenant_id, user_id)
    );

COMMENT ON TABLE iam.tenant_members IS 'Authorization control: defines which users have access to each tenant and with which role. Distinct from iam.tenant_customers, which records the business relationship (customer of the tenant).';

COMMENT ON COLUMN iam.tenant_members.role IS 'Role of the member within the tenant, used for authorization.';

-- +goose Down
DROP TABLE IF EXISTS iam.tenant_members;