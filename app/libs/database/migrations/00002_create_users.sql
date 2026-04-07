-- +goose Up
CREATE TABLE
    users (
        id UUID PRIMARY KEY,
        social_id TEXT,
        street TEXT,
        complement TEXT,
        NAME TEXT NOT NULL,
        email TEXT NOT NULL,
        photo TEXT,
        tax_id TEXT NOT NULL,
        phone TEXT NOT NULL,
        city TEXT,
        state TEXT,
        zip TEXT,
        neighborhood TEXT,
        ROLE TEXT,
        tenant_id UUID
    );

-- +goose Down
DROP TABLE IF EXISTS users;