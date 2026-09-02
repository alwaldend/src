#!/usr/bin/env bash
# Regenerate the checked agent system index artifacts into the workspace.
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

if [[ -n "${BUILD_WORKING_DIRECTORY:-}" ]]; then
    workspace="$BUILD_WORKING_DIRECTORY"
elif [[ -n "${BUILD_WORKSPACE_DIRECTORY:-}" ]]; then
    workspace="$BUILD_WORKSPACE_DIRECTORY"
else
    echo >&2 "unable to locate the repository workspace"
    exit 2
fi

binary="$(rlocation "${INDEX_CHECK_RLOCATION:?}")"
exec "${binary}" \
    --workspace-root "${workspace}" \
    --output tools/agents/catalogs/index.json \
    --markdown tools/agents/catalogs/index.md
