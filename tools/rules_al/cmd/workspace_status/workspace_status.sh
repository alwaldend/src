#!/usr/bin/env sh
set -eu

echo "release stamping requires the current versioning bootstrap" >&2
echo "run: tools/rules_al/cmd/al/al.sh version bazel -- <command>" >&2
exit 1
