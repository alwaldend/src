#!/usr/bin/env sh
set -eu

workspace_root="${BUILD_WORKSPACE_DIRECTORY:-$PWD}"
cd "${workspace_root}"

bootstrap_dir="${workspace_root}/out/al/bootstrap"
bootstrap_script="${bootstrap_dir}/al"
mkdir -p "${bootstrap_dir}"

bazel run --config=agent \
    --script_path="${bootstrap_script}" \
    @rules_al//:al
exec "${bootstrap_script}" "$@"
