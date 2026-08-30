#!/usr/bin/env bash

set -o errexit -o nounset -o pipefail

project_directory="$(
    cd -- "$(dirname -- "${BASH_SOURCE[0]}")"
    pwd -P
)"
workspace_root="$(
    git -C "${project_directory}" rev-parse --show-toplevel
)"
scratch_directory="${workspace_root}/out/mcp_cordis"

mkdir -p -- "${scratch_directory}"
run_script="$(mktemp "${scratch_directory}/launch.XXXXXXXX")"

cd -- "${workspace_root}"
bazel_agent run "--script_path=${run_script}" \
    //projects/mcp_cordis:mcp_cordis -- \
    --workspace-root "${workspace_root}" >&2
exec "${run_script}"
