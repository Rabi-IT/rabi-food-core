-- +goose Up
CREATE TABLE
    tenants (id UUID PRIMARY KEY, NAME TEXT NOT NULL);

-- +goose Down
DROP TABLE IF EXISTS tenants;