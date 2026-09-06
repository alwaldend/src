#!/usr/bin/env bash
# SessionEnd is synchronous and short-lived. Let systemd own the cleanup.
set -euo pipefail

list_workspace_markers() {
    # Restrict discovery to tracked files; never traverse build outputs.
    local markers=(
        MODULE.bazel
        WORKSPACE
        WORKSPACE.bazel
        ':(glob)**/MODULE.bazel'
        ':(glob)**/WORKSPACE'
        ':(glob)**/WORKSPACE.bazel'
    )
    git ls-files -z -- "${markers[@]}"
}

clean_workspace() {
    local workspace=$1
    local status=0

    # Do not start a fresh server for a workspace without build outputs.
    if [[ ! -L "$workspace/bazel-out" ]]; then
        return 0
    fi

    # The startup lock flag skips active builds instead of waiting for them.
    # bazel_agent cannot forward startup flags, so invoke Bazel directly.
    (
        cd "$workspace" || exit $?
        "$bazel_bin" --noblock_for_lock clean --config=agent --expunge_async
    ) || status=$?

    case "$status" in
    0) return 0 ;;
    9) printf 'Skipping busy Bazel workspace: %s\n' "$workspace" ;;
    *)
        printf 'Bazel cleanup failed for %s (exit %s)\n' "$workspace" "$status" >&2
        return 1
        ;;
    esac
}

clean_workspaces() {
    local marker workspace
    local failed=0
    local -A seen=()

    while IFS= read -r -d '' marker; do
        workspace=${marker%/*}
        if [[ "$workspace" == "$marker" ]]; then
            workspace=.
        fi

        # Several marker files can describe the same workspace.
        if [[ -n "${seen[$workspace]:-}" ]]; then
            continue
        fi
        seen[$workspace]=1

        clean_workspace "$workspace" || failed=1
    done
    return "$failed"
}

root=$(git rev-parse --show-toplevel)
git_dir=$(git rev-parse --absolute-git-dir)
common_dir=$(git rev-parse --path-format=absolute --git-common-dir)

# Only linked worktrees: never clean the primary checkout or a submodule.
if [[ "$git_dir" != "$common_dir"/worktrees/* ]]; then
    exit 0
fi

if [[ "${1:-}" != --worker ]]; then
    bash_bin=$(command -v bash)
    exec systemd-run --user --quiet --collect \
        --property=KillMode=process \
        --working-directory="$root" \
        --setenv="PATH=$PATH" \
        "$bash_bin" "$root/tools/bazel_session_cleanup/hooks/session_end.sh" --worker
fi

cd "$root"
bazel_bin=$(command -v bazel)
# pipefail reports discovery failures as well as cleanup failures.
list_workspace_markers | clean_workspaces
