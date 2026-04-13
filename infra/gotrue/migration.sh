#!/bin/sh
set -eu

psql \
  -h postgres \
  -U "$TEST_DATABASE_USER" \
  -d "$TEST_DATABASE_NAME" \
  -f /docker-entrypoint-initdb.d/migration.sql