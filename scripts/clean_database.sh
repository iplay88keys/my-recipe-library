#!/usr/bin/env bash

set -e

root_dir="$(git rev-parse --show-toplevel)"

source "${root_dir}/scripts/dev_db_creds.sh"

echo "Clearing Database"
podman exec db_db_1 mysql -u "${DATABASE_USERNAME}" \
    -p"${DATABASE_PASSWORD}" \
    -e "DROP DATABASE IF EXISTS ${DATABASE_NAME}; CREATE DATABASE ${DATABASE_NAME};"
