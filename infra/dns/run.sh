#!/usr/bin/env bash

function main() {
    # https://github.com/bazelbuild/rules_shell/blob/main/shell/private/sh_executable.bzl
    # --- begin runfiles.bash initialization v3 ---
    # Copy-pasted from the Bazel Bash runfiles library v3.
    set -uo pipefail; set +e; f=bazel_tools/tools/bash/runfiles/runfiles.bash
    # shellcheck disable=SC1090
    source "${RUNFILES_DIR:-/dev/null}/$f" 2>/dev/null || \
        source "$(grep -sm1 "^$f " "${RUNFILES_MANIFEST_FILE:-/dev/null}" | cut -f2- -d' ')" 2>/dev/null || \
        source "$0.runfiles/$f" 2>/dev/null || \
        source "$(grep -sm1 "^$f " "$0.runfiles_manifest" | cut -f2- -d' ')" 2>/dev/null || \
        source "$(grep -sm1 "^$f " "$0.exe.runfiles_manifest" | cut -f2- -d' ')" 2>/dev/null || \
        { echo>&2 "ERROR: cannot find $f"; exit 1; }; f=; set -e
    # --- end runfiles.bash initialization v3 ---

    set -eu

    runfiles_export_envvars

    root="$(runfiles_current_repository)"
    if [ -z "${root}" ]; then
        root="_main"
    fi
    vault="$(rlocation "${root}/tools/vault/vault.script.sh")"
    dnscontrol="$(rlocation "${root}/tools/dnscontrol/dnscontrol.script.sh")"

    "${vault}" kv get sre/data/infra/dns
}

main "${@}"
