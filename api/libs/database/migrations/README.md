# Migrations

Migrations are managed by [Goose](https://github.com/pressly/goose).

## Schemas

| Schema         | Owner         | Description                              |
|----------------|---------------|------------------------------------------|
| `iam`          | This app      | Tenants, tenant-customer relationships   |
| `auth`         | GoTrue        | Users, sessions, refresh tokens          |
| `catalog`      | This app      | Products and categories                  |
| `commerce`     | This app      | Orders                                   |
| `subscription` | This app      | Subscriptions                            |

## User management

Users are managed by **GoTrue** (Supabase Auth). The `auth` schema is created and
owned by GoTrue on startup — do not add migrations targeting it.

User identity (email, password, OAuth, sessions) lives in `auth.users`.
All other domain data that references a user does so by the GoTrue user UUID.
