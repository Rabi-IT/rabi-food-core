-- +goose Up
CREATE SCHEMA IF NOT EXISTS catalog;

CREATE TABLE
    catalog.categories (
        id UUID PRIMARY KEY,
        tenant_id UUID NOT NULL REFERENCES iam.tenants (id),
        NAME TEXT NOT NULL,
        description TEXT
    );

-- +goose Down
DROP TABLE IF EXISTS catalog.categories;

DROP SCHEMA IF EXISTS catalog;
