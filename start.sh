#!/bin/sh

set -e

echo "running db migration"
source ./app.env 
/app/migrate -path /app/migrations -database "$DSN" -verbose up

echo "start the app"
exec "$@"