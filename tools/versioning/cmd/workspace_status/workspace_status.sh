#!/usr/bin/env sh
set -eu

echo "release stamping requires the current versioning bootstrap" >&2
echo "run: tools/versioning/cmd/versioning/versioning.sh bazel -- <command>" >&2
exit 1
