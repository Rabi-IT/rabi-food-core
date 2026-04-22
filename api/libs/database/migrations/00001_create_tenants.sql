-- +goose Up
CREATE SCHEMA IF NOT EXISTS iam;

CREATE TABLE
    iam.tenants (
        id       UUID PRIMARY KEY,
        name     TEXT NOT NULL,
        slug     TEXT NOT NULL UNIQUE,
        language TEXT NOT NULL DEFAULT 'pt-BR'
    );

-- +goose Down
DROP TABLE IF EXISTS iam.tenants;

DROP SCHEMA IF EXISTS iam;
