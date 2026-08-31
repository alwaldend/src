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
case "$(basename "$0")" in
buf)
    tool="${runfiles_dir}/rules_binary_toolchain++binary_toolchain_extension+com_alwaldend_src_tools_buf/buf_binary_/buf_native_binary.exe"
    ;;
taplo)
    tool="${runfiles_dir}/rules_binary_toolchain++binary_toolchain_extension+com_alwaldend_src_tools_taplo/taplo_binary_/taplo_native_binary.exe"
    ;;
*)
    printf 'unsupported adapter name: %s\n' "$(basename "$0")" >&2
    exit 2
    ;;
esac

exec "${tool}" "$@"
