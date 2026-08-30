#!/usr/bin/env bash

set -euo pipefail

runtime_file="$(find -L "${RUNFILES_DIR:?}" -type f \
    -name filegroup-runtime.txt -print -quit)"
test -n "${runtime_file}"
test "$(cat "${runtime_file}")" = "filegroup-runtime"
