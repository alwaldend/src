#!/usr/bin/env sh

set -eux
# shellcheck source=../bzlenv/activate.sh
. "$(bazel run //tools/bzlenv)"
bazel build //...
exec bazel test //...
