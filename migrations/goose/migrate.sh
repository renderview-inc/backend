#!/bin/sh
set -e

DBSTRING="clickhouse://${CLICKHOUSE_USER}:${CLICKHOUSE_PASSWORD}@${CLICKHOUSE_HOST}:${CLICKHOUSE_PORT}/${CLICKHOUSE_DB}"
VERSION="${MIGRATION_CLICKHOUSE_VERSION:-latest}"

echo "applying migrations up to version: ${VERSION}"

if [ "$VERSION" = "latest" ]; then
    goose -dir /migrations/goose/sql clickhouse "$DBSTRING" up
else
    goose -dir /migrations/goose/sql clickhouse "$DBSTRING" up-to "$VERSION"
fi