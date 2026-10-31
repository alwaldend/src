#!/usr/bin/env sh

set -eux

if ! command -v al >/dev/null 2>&1; then
    echo "al is required; bootstrap it with:" >&2
    echo "bazel run --config=agent //tools/al:install" >&2
    exit 127
fi

al bazel run //:format
al bazel test //:repo_quality_test
al bazel build --config=lint //...
