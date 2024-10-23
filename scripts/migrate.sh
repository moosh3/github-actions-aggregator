#!/bin/bash

set -e

# Configuration
DB_HOST=${DB_HOST:-"localhost"}
DB_PORT=${DB_PORT:-5432}
DB_USER=${DB_USER:-"postgres"}
DB_PASSWORD=${DB_PASSWORD:-"password"}
DB_NAME=${DB_NAME:-"github_actions_aggregator"}
MIGRATIONS_DIR=${MIGRATIONS_DIR:-"./migrations"}

# Command-line arguments
COMMAND=$1

# Check if the migrate tool is installed
if ! [ -x "$(command -v migrate)" ]; then
  echo 'Error: migrate is not installed.' >&2
  echo 'You can install it by running:'
  echo '  curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz'
  echo '  sudo mv migrate.linux-amd64 /usr/local/bin/migrate'
  exit 1
fi

# Build the database URL
DB_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"

case $COMMAND in
  "up")
    echo "Applying all up migrations..."
    migrate -path "${MIGRATIONS_DIR}" -database "${DB_URL}" up
    ;;
  "down")
    echo "Reverting the last migration..."
    migrate -path "${MIGRATIONS_DIR}" -database "${DB_URL}" down 1
    ;;
  "force")
    VERSION=$2
    if [ -z "$VERSION" ]; then
      echo "Please specify the version to force."
      exit 1
    fi
    echo "Forcing migration version to ${VERSION}..."
    migrate -path "${MIGRATIONS_DIR}" -database "${DB_URL}" force "${VERSION}"
    ;;
  "version")
    echo "Current migration version:"
    migrate -path "${MIGRATIONS_DIR}" -database "${DB_URL}" version
    ;;
  *)
    echo "Usage: $0 {up|down|force <version>|version}"
    exit 1
    ;;
esac

echo "Migration command '${COMMAND}' completed successfully."
