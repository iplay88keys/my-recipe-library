#!/usr/bin/env bash

root_dir="$(git rev-parse --show-toplevel)"

set -e

source "${root_dir}/scripts/dev_db_creds.sh"

echo "Stopping the database if currently running"
"${root_dir}/scripts/stop_database.sh"

pushd "${root_dir}/db"
    echo "Bringing up the new database"
    podman compose up -d
popd

exit_code=1
set +e
while [[ "${exit_code}" -eq 1 ]]; do
    echo "Waiting for mysql to be available..."
    # Use podman exec to check from inside the container to avoid host authentication issues
    podman exec db_db_1 mysqladmin -u "${DATABASE_USERNAME}" \
        -p"${DATABASE_PASSWORD}" ping  > /dev/null 2>&1

    exit_code=$?

    sleep 2
done
set -e
