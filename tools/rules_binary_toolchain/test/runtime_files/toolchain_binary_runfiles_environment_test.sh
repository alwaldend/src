#!/usr/bin/env bash

set -euo pipefail

wrapper="$(find -L "${RUNFILES_DIR:?}" -type f \
    -path '*/fixture_native_binary_/*' -print -quit)"
test -n "${wrapper}"

env -u RUNFILES_DIR -u RUNFILES_MANIFEST_FILE RUNFILES="${RUNFILES_DIR}" \
    "${wrapper}"
