#!/usr/bin/env sh

set -eu
bazel build //tools/bzlenv
# shellcheck source=../bzlenv/activate.sh
. "$(bazel run //tools/bzlenv)"
set -x
bazel build //...
exec bazel test //...
