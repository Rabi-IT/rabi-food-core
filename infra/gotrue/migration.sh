#!/bin/sh
set -eu

psql \
  -h "${DATABASE_HOST:-postgres}" \
  -p "${DATABASE_PORT:-5432}" \
  -U "$DATABASE_USER" \
  -d "$DATABASE_NAME" \
  -c "CREATE SCHEMA IF NOT EXISTS auth;"