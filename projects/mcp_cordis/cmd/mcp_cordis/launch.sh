#!/usr/bin/env bash

set -o errexit -o nounset -o pipefail

project_directory="$(
    cd -- "$(dirname -- "${BASH_SOURCE[0]}")"
    pwd -P
)"
workspace_root="$(
    git -C "${project_directory}" rev-parse --show-toplevel
)"
task_id="${AGENT_TASK_ID:-mcp-cordis-${PPID}}"
run_id="${AGENT_RUN_ID:-run-${BASHPID}}"
worker_id="${AGENT_WORKER_ID:-worker-${BASHPID}}"
ownership_pattern='^[a-z][a-z0-9]*([._-][a-z0-9]+)*$'
for value in "${task_id}" "${run_id}" "${worker_id}"; do
    if [[ ! "${value}" =~ ${ownership_pattern} ]]; then
        echo "mcp_cordis: invalid task/run/worker identity" >&2
        exit 2
    fi
done
scratch_directory="${workspace_root}/out/${task_id}/mcp_cordis/runs/${run_id}"

mkdir -p -- "${scratch_directory}"
run_script="$(mktemp "${scratch_directory}/launch.XXXXXXXX")"

cd -- "${workspace_root}"
bazel_agent run "--script_path=${run_script}" \
    //projects/mcp_cordis:mcp_cordis -- \
    --workspace-root "${workspace_root}" \
    --task-id "${task_id}" \
    --run-id "${run_id}" \
    --worker-id "${worker_id}" >&2
exec "${run_script}"
