-- +goose Up
CREATE TABLE
    categories (
        id UUID PRIMARY KEY,
        tenant_id UUID NOT NULL,
        NAME TEXT NOT NULL,
        description TEXT
    );

-- +goose Down
DROP TABLE IF EXISTS categories;