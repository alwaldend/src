#!/usr/bin/env bash
# Smoke-check offline context rendering; this does not verify catalog freshness.
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

binary="$(rlocation "$AGENT_SYSTEM_RLOCATION")"
workspace="${BUILD_WORKSPACE_DIRECTORY:-}"
if [[ -z "$workspace" && -n "${WORKSPACE_MARKER:-}" ]]; then
    workspace="$(dirname "$(rlocation "$WORKSPACE_MARKER")")"
fi
"$binary" --workspace-root "$workspace" "$@" >/dev/null
