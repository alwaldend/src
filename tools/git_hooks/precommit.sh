#!/usr/bin/env sh

set -eux
# shellcheck disable=SC1090
. "$(bazel run //tools/bzlenv)"
bazel build //...
exec bazel test //...
