#!/usr/bin/env sh

set -eux

bazel run //:format
bazel test //:repo_quality_test
bazel build --config=lint //...
