#!/usr/bin/env bash

set -e

root_dir="$(git rev-parse --show-toplevel)"

source "${root_dir}/scripts/dev_db_creds.sh"

buildUI=true
while test $# -gt 0
do
    case "$1" in
        --skip-ui-build)
            buildUI=false
            ;;
        --*)
            echo "bad option $1"
            exit 0
            ;;
        *)
            echo "bad argument $1"
            exit 0
            ;;
    esac
    shift
done

if [[ "${buildUI}" = "true" ]]; then
    pushd "${root_dir}/ui"
        yarn build
    popd
fi

echo "Checking database container status..."
container_running=false
set +e
if podman ps --format "{{.Names}}" | grep -q "db_db_1"; then
    container_running=true
    echo "Database container is running, checking MySQL availability..."
    podman exec db_db_1 mysqladmin -u "${DATABASE_USERNAME}" \
        -p"${DATABASE_PASSWORD}" ping  > /dev/null 2>&1
    mysql_available=$?
else
    echo "Database container is not running"
    mysql_available=1
fi
set -e

if [[ "${mysql_available}" -eq 1 ]]; then
    echo "MySQL is not available, starting database..."

    "${root_dir}/scripts/stop_database.sh" > /dev/null 2>&1
    "${root_dir}/scripts/start_database.sh" > /dev/null 2>&1
else
    echo "MySQL is available, cleaning database..."
    "${root_dir}/scripts/clean_database.sh"
fi

function finish {
    "${root_dir}/scripts/stop_database.sh" > /dev/null 2>&1
}
trap finish EXIT

echo "Migrating the database"
"${root_dir}/scripts/migrate_database.sh"

echo "Importing example data"
"${root_dir}/scripts/import_example_database_data.sh"

echo "Exporting env vars"
mysql_url="mysql://${DATABASE_USERNAME}:${DATABASE_PASSWORD}@tcp(${DATABASE_HOST}:${DATABASE_PORT})/${DATABASE_NAME}"
export MYSQL_CREDS="{\"url\": \"${mysql_url}\"}"
export REDIS_URL="redis://:@127.0.0.1:6379"
export ACCESS_SECRET="access_secret"
export REFRESH_SECRET="refresh_secret"

echo "----------------------------------------------------"
echo "MySQL url:"
echo "${mysql_url}"
echo "----------------------------------------------------"

echo "Starting the app"
pushd "${root_dir}"
  go run main.go
popd
