#!/usr/bin/env sh
set -eu

script_dir=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
workspace_root=$(CDPATH= cd -- "${script_dir}/../../../.." && pwd)
cd "${workspace_root}"

bootstrap_dir="${workspace_root}/out/versioning/bootstrap"
bootstrap_script="${bootstrap_dir}/versioning"
mkdir -p "${bootstrap_dir}"

bazel_agent bazel run \
    --script_path="${bootstrap_script}" \
    //tools/versioning/cmd/versioning
exec "${bootstrap_script}" --repo "${workspace_root}" "$@"
