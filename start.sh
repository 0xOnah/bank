#!/bin/sh

set -e

echo "running db migration"
/app/migrate -path /app/migrations -database "$DSN" -verbose up

echo "start the app"
exec "$@"