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

mode="${1:?usage: go_module_quality.sh check|fix}"
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

go="$(rlocation "${GO_RLOCATION:?}")"
cd "$workspace"
status=0
while IFS= read -r -d '' path &&
    IFS= read -r -d '' _attribute &&
    IFS= read -r -d '' value; do
    if [[ "$value" == "set" || "$value" == "true" ]]; then
        continue
    fi

    file="${workspace}/${path}"
    if [[ "$path" == "go.work" ]]; then
        kind=work
    else
        kind=mod
    fi

    if [[ "$mode" == "fix" ]]; then
        "$go" "$kind" edit -fmt "$file"
    elif ! diff -u "$file" <("$go" "$kind" edit -fmt -print "$file"); then
        status=1
    fi
done < <(
    git -C "$workspace" ls-files -z 'go.mod' 'go.work' '**/go.mod' '**/go.work' |
        git -C "$workspace" check-attr -z --stdin rules-lint-ignored
)

exit "$status"
