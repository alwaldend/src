#!/usr/bin/env sh

set -eu

exec "$(dirname "$0")/deploy_project.sh" "$@"
