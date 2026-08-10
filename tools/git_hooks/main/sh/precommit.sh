#!/usr/bin/env sh

set -eux
bazel run //:buildifier
bazel run //:tf.direct -- -chdir "${PWD}" fmt --recursive
bazel run @rules_go//go -- fmt "${PWD}"
