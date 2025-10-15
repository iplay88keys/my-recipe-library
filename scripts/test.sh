#!/usr/bin/env bash

set -e

root_dir="$(git rev-parse --show-toplevel)"

source "${root_dir}/scripts/dev_db_creds.sh"

skipUI=false
skipIntegration=false
skipBackend=false

while test $# -gt 0; do
    case "$1" in
        --skip-ui)
            skipUI=true
            ;;
        --skip-backend)
            skipBackend=true
            ;;
        --skip-integration)
            skipIntegration=true
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

if [[ "${skipUI}" = "false" ]]; then
    pushd "${root_dir}/ui"
        echo "Running 'yarn test'"
        yarn test --watchAll=false

        echo "Compiling the UI"
        yarn build
    popd
fi

if [[ "${skipIntegration}" = "false" ]]; then
    echo "Setting up database for testing..."

    # Always stop any existing containers first to ensure clean state
    "${root_dir}/scripts/stop_database.sh" > /dev/null 2>&1

    # Start fresh database
    echo "Starting database for testing"
    "${root_dir}/scripts/start_database.sh" > /dev/null 2>&1

    function finish {
      echo "Cleaning up database after tests"
      "${root_dir}/scripts/stop_database.sh" > /dev/null 2>&1
    }
    trap finish EXIT

    "${root_dir}/scripts/migrate_database.sh"
fi


echo "Setting required env vars"
export REDIS_URL="redis://:@127.0.0.1:6379"
export ACCESS_SECRET="access_secret"
export REFRESH_SECRET="refresh_secret"

if [[ "${skipBackend}" = "false" ]]; then
    echo "Running ginkgo for everything except integration"
    pushd ${root_dir}
      ginkgo -r -p -skip-package pkg/integration --output-interceptor-mode=none
    popd
fi

if [[ "${skipIntegration}" = "false" ]]; then
    echo "Running ginkgo for integration"
    pushd ${root_dir}
      ginkgo -r pkg/integration --output-interceptor-mode=none
    popd
fi
