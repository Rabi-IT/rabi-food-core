# GoTrue migration bootstrap

This directory exists to handle a GoTrue bootstrap detail when it runs against a vanilla Postgres instance.

## Problem

GoTrue runs its own migrations automatically during startup.

In this project, the schema name is configured through the GoTrue environment as `GOTRUE_DB_NAMESPACE=auth`.

However, even with that configuration in place, GoTrue still expects the schema to already exist. Its first migration creates tables such as `auth.users`, `auth.refresh_tokens`, and others, but it does not run `CREATE SCHEMA auth` first.

In the official Supabase stack this is already prepared. In this project, since the database is a plain Postgres instance, that schema must be created explicitly before the GoTrue container starts.

Without that, GoTrue fails during startup with an error equivalent to:

```text
ERROR: schema "auth" does not exist
```

## Solution

The `gotrue-migration` service runs before `gotrue` in [infra/docker-compose.db.yaml](../docker-compose.db.yaml).

Files:

- [infra/gotrue/migration.sh](./migration.sh): runs `psql` against the `postgres` container
- [infra/gotrue/migration.sql](./migration.sql): creates the `auth` schema

Flow:

1. `postgres` starts
2. `gotrue-migration` runs `CREATE SCHEMA IF NOT EXISTS auth;`
3. `gotrue` starts
4. GoTrue runs its native migrations normally

## Note

This directory does not replace GoTrue migrations. It only prepares the minimum condition required for GoTrue's native migrations to run successfully.