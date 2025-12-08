#!/usr/bin/env sh

set -eu
bazel build //tools/bzlenv
# shellcheck disable=SC1090
. "$(bazel run //tools/bzlenv)"
set -x
bazel build //...
exec bazel test //...
