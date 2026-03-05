#!/usr/bin/env sh

set -eux
set "${PWD}" "${PWD}/projects/rules_template"
for dir in "${@}"; do
    cd "${dir}"
    bazel build "//..."
    bazel test "//..."
done
