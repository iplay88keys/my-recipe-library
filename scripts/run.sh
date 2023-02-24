#!/usr/bin/env bash

set -e

. ./scripts/dev_db_creds.sh

#buildUI=true
while test $# -gt 0
do
    case "$1" in
#        --skip-ui-build)
#            buildUI=false
#            ;;
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

#if [[ "${buildUI}" = "true" ]]; then
#    pushd ui
#        yarn build
#    popd
#fi

echo "Restarting the database"
./scripts/start_database.sh

function finish {
    echo "Stopping the database"
    ./scripts/stop_database.sh
}
trap finish EXIT

echo "Migrating the database"
./scripts/migrate_database.sh

echo "Importing example data"
./scripts/import_example_database_data.sh

echo "Exporting env vars"
export MYSQL_CREDS="{\"url\": \"mysql://${DATABASE_USERNAME}:${DATABASE_PASSWORD}@tcp(${DATABASE_HOST}:${DATABASE_PORT})/${DATABASE_NAME}\"}"
export REDIS_URL="redis://:@127.0.0.1:6379"
export ACCESS_SECRET="access_secret"
export REFRESH_SECRET="refresh_secret"

echo "Starting the apps"
watchtower -config-file ./scripts/reload.yaml
pid=$!
