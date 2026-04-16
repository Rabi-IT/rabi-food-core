-- +goose Up
CREATE TABLE
    iam.tenant_members (
        tenant_id  UUID        NOT NULL REFERENCES iam.tenants (id),
        user_id    UUID        NOT NULL REFERENCES auth.users (id),
        role       VARCHAR(20) NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        PRIMARY KEY (tenant_id, user_id)
    );

-- +goose Down
DROP TABLE IF EXISTS iam.tenant_members;
