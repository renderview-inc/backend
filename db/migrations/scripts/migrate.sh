#!/bin/bash

set -e

echo "$(pwd)"
if [ -z "$MIGRATION_VERSION" ]; then
  echo "MIGRATION_VERSION is not set. Running all migrations..."
  liquibase --url=jdbc:${DB_DRIVER}://${DB_HOST}:${DB_PORT}/${POSTGRES_DB} \
            --username=${POSTGRES_USER} \
            --password=${POSTGRES_PASSWORD} \
            --changelog-file=changelogs/db.changelog-root.yaml \
            update
else
  echo "Running migrations up to version $MIGRATION_VERSION..."
  liquibase --url=jdbc:${DB_DRIVER}://${DB_HOST}:${DB_PORT}/${POSTGRES_DB} \
            --username=${POSTGRES_USER} \
            --password=${POSTGRES_PASSWORD} \
            --changelog-file=changelogs/db.changelog-root.yaml \
            update-to-tag --tag=$MIGRATION_VERSION
fi
