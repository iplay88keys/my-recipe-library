#!/usr/bin/env bash

set -e

root_dir="$(git rev-parse --show-toplevel)"

source "${root_dir}/scripts/dev_db_creds.sh"

pushd "${root_dir}/db"
    echo "Importing example data into the database"
    podman exec -i db_db_1 mysql -u "${DATABASE_USERNAME}" \
        -p"${DATABASE_PASSWORD}" \
        -D "${DATABASE_NAME}" < examples/example.sql
popd
