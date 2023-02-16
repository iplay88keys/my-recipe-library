#!/usr/bin/env bash

set -e

root_dir="$(git rev-parse --show-toplevel)"

source "${root_dir}/scripts/dev_db_creds.sh"

pushd "${root_dir}/migrations"
    echo "Migrating the database"
    flyway migrate \
        -url="jdbc:mysql://${DATABASE_HOST}:${DATABASE_PORT}/${DATABASE_NAME}" \
        -user="${DATABASE_USERNAME}" \
        -password="${DATABASE_PASSWORD}" \
        -locations=filesystem:.
popd
