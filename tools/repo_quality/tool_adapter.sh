#!/usr/bin/env bash

set -euo pipefail

workspace="${BUILD_WORKSPACE_DIRECTORY:-${PWD}}"
quality_state="${workspace}/out/repo_quality"

mkdir -p \
    "${quality_state}/buf" \
    "${quality_state}/cache" \
    "${quality_state}/config" \
    "${quality_state}/tmp"
export BUF_CACHE_DIR="${quality_state}/buf"
export TMPDIR="${quality_state}/tmp"
export XDG_CACHE_HOME="${quality_state}/cache"
export XDG_CONFIG_HOME="${quality_state}/config"

runfiles_dir="${RUNFILES_DIR:-${0}.runfiles}"
exec "${runfiles_dir}/repo_quality_tools/$(basename "$0")" "$@"
