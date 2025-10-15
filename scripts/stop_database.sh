#!/usr/bin/env bash

set -e

root_dir="$(git rev-parse --show-toplevel)"

pushd "${root_dir}/db"
  podman compose down
popd
