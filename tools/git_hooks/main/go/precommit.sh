#!/usr/bin/env sh

set -eux

if ! command -v bazel_agent >/dev/null 2>&1; then
    echo "bazel_agent is required; bootstrap it with:" >&2
    echo "bazel --batch run --config=agent //projects/bazel_agent:install" >&2
    exit 127
fi

bazel_agent run //:buildifier
bazel_agent run //:tf.direct -- -chdir "${PWD}" fmt --recursive
bazel_agent run @rules_go//go -- fmt "${PWD}"
