#!/bin/sh
set -e

echo "start the app"
export DSN="$(cat /run/secrets/db_password)"
exec "$@"