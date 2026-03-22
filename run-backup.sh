#!/bin/bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

# load .env
if [ -f "$SCRIPT_DIR/.env" ]; then
    set -a
    source "$SCRIPT_DIR/.env"
    set +a
else
    echo ".env file could not be found in: $SCRIPT_DIR"
    exit 1
fi

# check backup path
if [ -z "${BACKUP_FOLDER_PATH:-}" ]; then
    echo "BACKUP_DIR not set in .env... use current directory ..."
    BACKUP_FOLDER_PATH="$SCRIPT_DIR"
fi

mkdir -p "$BACKUP_FOLDER_PATH"
echo "backup directory set to: $BACKUP_FOLDER_PATH"

echo "starting backup process..."
"$SCRIPT_DIR/pg-docker-backup" encrypt \
  -c "$CONTAINER_NAME" \
  -n "$DB_NAME" \
  -u "$DB_USER" \
  -p "$DB_PASSWORD"

echo "backup done."
