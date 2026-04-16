-- +goose Up
CREATE TABLE
    catalog.products (
        id UUID PRIMARY KEY,
        tenant_id UUID NOT NULL REFERENCES iam.tenants (id),
        NAME TEXT NOT NULL,
        description TEXT,
        photo TEXT,
        category_id UUID REFERENCES catalog.categories (id),
        discount_rules JSONB NOT NULL DEFAULT '[]',
        unit TEXT NOT NULL,
        price BIGINT NOT NULL,
        is_active BOOLEAN NOT NULL DEFAULT FALSE,
        created_at TIMESTAMPTZ,
        updated_at TIMESTAMPTZ,
        deleted_at TIMESTAMPTZ
    );

CREATE INDEX idx_products_deleted_at ON catalog.products (deleted_at);

-- +goose Down
DROP TABLE IF EXISTS catalog.products;
