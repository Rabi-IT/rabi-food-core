-- +goose Up
CREATE TABLE
    CATALOG.products (
        id UUID PRIMARY KEY,
        tenant_id UUID NOT NULL REFERENCES iam.tenants (id),
        NAME TEXT NOT NULL,
        description TEXT,
        photo TEXT,
        category_id UUID REFERENCES CATALOG.categories (id),
        discount_rules JSONB NOT NULL DEFAULT '[]',
        unit TEXT NOT NULL,
        price BIGINT NOT NULL,
        is_active BOOLEAN NOT NULL DEFAULT FALSE,
        created_at TIMESTAMPTZ,
        updated_at TIMESTAMPTZ,
        deleted_at TIMESTAMPTZ
    );

COMMENT ON COLUMN CATALOG.products.price IS 'Price in cents.';

COMMENT ON COLUMN CATALOG.products.unit IS 'Unit of measurement (e.g. un, kg, L, g).';

COMMENT ON COLUMN CATALOG.products.photo IS 'URL of the product image.';

COMMENT ON COLUMN CATALOG.products.discount_rules IS 'JSON array of discount rules applicable to this product.';

COMMENT ON COLUMN CATALOG.products.deleted_at IS 'Soft delete — record is logically removed when set. Non-null records are excluded from listings.';

CREATE INDEX idx_products_deleted_at ON CATALOG.products (deleted_at);

-- +goose Down
DROP TABLE IF EXISTS CATALOG.products;