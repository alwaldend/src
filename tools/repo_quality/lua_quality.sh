#!/usr/bin/env bash

set -euo pipefail

f=bazel_tools/tools/bash/runfiles/runfiles.bash
# shellcheck disable=SC1090
source "${RUNFILES_DIR:-/dev/null}/$f" 2>/dev/null ||
    source "$(grep -sm1 "^$f " \
        "${RUNFILES_MANIFEST_FILE:-/dev/null}" | cut -f2- -d' ')" \
        2>/dev/null ||
    source "$0.runfiles/$f" 2>/dev/null ||
    source "$(grep -sm1 "^$f " "$0.runfiles_manifest" |
        cut -f2- -d' ')" 2>/dev/null || {
    echo >&2 "ERROR: cannot find $f"
    exit 1
}
runfiles_export_envvars

mode="${1:?usage: lua_quality.sh check|fix}"
case "$mode" in
check | fix) ;;
*)
    echo >&2 "unknown mode: $mode"
    exit 2
    ;;
esac

if [[ -n "${BUILD_WORKING_DIRECTORY:-}" ]]; then
    workspace="$BUILD_WORKING_DIRECTORY"
elif [[ -n "${BUILD_WORKSPACE_DIRECTORY:-}" ]]; then
    workspace="$BUILD_WORKSPACE_DIRECTORY"
elif [[ -n "${WORKSPACE_MARKER:-}" ]]; then
    workspace="$(dirname "$(realpath "${WORKSPACE_MARKER}")")"
else
    echo >&2 "unable to locate the repository workspace"
    exit 2
fi

stylua="$(rlocation "${STYLUA_RLOCATION:?}")"
stylua_config="$(rlocation "${STYLUA_CONFIG_RLOCATION:?}")"
selene="$(rlocation "${SELENE_RLOCATION:?}")"
selene_config="$(rlocation "${SELENE_CONFIG_RLOCATION:?}")"

files=()
while IFS= read -r -d '' path &&
    IFS= read -r -d '' _attribute &&
    IFS= read -r -d '' value; do
    if [[ "$value" != "set" && "$value" != "true" ]]; then
        files+=("${workspace}/${path}")
    fi
done < <(
    git -C "$workspace" ls-files -z '*.lua' |
        git -C "$workspace" check-attr -z --stdin rules-lint-ignored
)

if [[ ${#files[@]} -eq 0 ]]; then
    exit 0
fi

if [[ "$mode" == "fix" ]]; then
    exec "$stylua" --config-path "$stylua_config" "${files[@]}"
fi

"$stylua" --config-path "$stylua_config" --check "${files[@]}"
(
    cd "$(dirname "$selene_config")"
    "$selene" \
        --config "$(basename "$selene_config")" \
        --color Never \
        "${files[@]}"
)
