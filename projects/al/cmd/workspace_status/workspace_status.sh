#!/usr/bin/env sh
set -eu

if ! command -v al >/dev/null 2>&1; then
    echo "release stamping requires AL on PATH" >&2
    echo "install it with: bazel run //tools/al:install" >&2
    exit 1
fi

exec al version --repository . bazel-status
