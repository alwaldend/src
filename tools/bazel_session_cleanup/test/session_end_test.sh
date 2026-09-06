#!/usr/bin/env bash
set -euo pipefail

hook="$TEST_SRCDIR/$TEST_WORKSPACE/tools/bazel_session_cleanup/hooks/session_end.sh"
export FIXTURE_ROOT="$TEST_TMPDIR/worktree with spaces"
export FIXTURE_LOG="$TEST_TMPDIR/calls"
mkdir -p "$FIXTURE_ROOT/nested project" "$FIXTURE_ROOT/unbuilt" "$TEST_TMPDIR/bin"
ln -s "$(command -v bash)" "$TEST_TMPDIR/bin/bash"
ln -s /nonexistent "$FIXTURE_ROOT/bazel-out"
ln -s /nonexistent "$FIXTURE_ROOT/nested project/bazel-out"
export PATH="$TEST_TMPDIR/bin:$PATH"

cat >"$TEST_TMPDIR/bin/git" <<'EOF'
#!/usr/bin/env bash
case "$*" in
  'rev-parse --show-toplevel') printf '%s\n' "$FIXTURE_ROOT" ;;
  'rev-parse --absolute-git-dir') printf '%s\n' "${FIXTURE_GIT_DIR:-/repo/.git/worktrees/task}" ;;
  'rev-parse --path-format=absolute --git-common-dir') printf '/repo/.git\n' ;;
  'ls-files '*)
    if [[ "${FIXTURE_DISCOVERY_STATUS:-0}" != 0 ]]; then
      exit "$FIXTURE_DISCOVERY_STATUS"
    fi
    printf '%s\0' MODULE.bazel WORKSPACE 'nested project/MODULE.bazel' unbuilt/MODULE.bazel ;;
  *) exit 1 ;;
esac
EOF
cat >"$TEST_TMPDIR/bin/systemd-run" <<'EOF'
#!/usr/bin/env bash
printf '%s\n' "$@" > "$FIXTURE_LOG"
exit "${FIXTURE_SYSTEMD_STATUS:-0}"
EOF
cat >"$TEST_TMPDIR/bin/bazel" <<'EOF'
#!/usr/bin/env bash
printf '%s\n' "$PWD" "$*" >> "$FIXTURE_LOG"
if [[ "$PWD" == "$FIXTURE_ROOT" ]]; then
  exit "${FIXTURE_BAZEL_STATUS:-0}"
fi
EOF
chmod +x "$TEST_TMPDIR/bin/git" "$TEST_TMPDIR/bin/systemd-run" "$TEST_TMPDIR/bin/bazel"

# Enqueue from a nested cwd without running Bazel in the hook process.
cd "$FIXTURE_ROOT/nested project"
bash "$hook"
grep -Fx -- '--property=KillMode=process' "$FIXTURE_LOG"
grep -Fx -- "--working-directory=$FIXTURE_ROOT" "$FIXTURE_LOG"
grep -Fx -- "$TEST_TMPDIR/bin/bash" "$FIXTURE_LOG"
grep -Fx -- "$FIXTURE_ROOT/tools/bazel_session_cleanup/hooks/session_end.sh" "$FIXTURE_LOG"
grep -Fx -- '--worker' "$FIXTURE_LOG"

# Primary checkout and submodules never queue cleanup.
for git_dir in /repo/.git /repo/.git/modules/submodule; do
    : >"$FIXTURE_LOG"
    FIXTURE_GIT_DIR="$git_dir" bash "$hook"
    [[ ! -s "$FIXTURE_LOG" ]]
done

# Worker deduplicates markers and skips workspaces without outputs.
: >"$FIXTURE_LOG"
bash "$hook" --worker
[[ $(wc -l <"$FIXTURE_LOG") == 4 ]]
grep -Fx -- "$FIXTURE_ROOT" "$FIXTURE_LOG"
grep -Fx -- "$FIXTURE_ROOT/nested project" "$FIXTURE_LOG"
[[ $(grep -Fxc -- '--noblock_for_lock clean --config=agent --expunge_async' "$FIXTURE_LOG") == 2 ]]

# A busy root is skipped, while the nested workspace still gets cleaned.
: >"$FIXTURE_LOG"
FIXTURE_BAZEL_STATUS=9 bash "$hook" --worker
[[ $(wc -l <"$FIXTURE_LOG") == 4 ]]

# Real cleanup and enqueue failures remain visible.
if FIXTURE_BAZEL_STATUS=1 bash "$hook" --worker; then
    exit 1
fi
if FIXTURE_SYSTEMD_STATUS=1 bash "$hook"; then
    exit 1
fi
if FIXTURE_DISCOVERY_STATUS=1 bash "$hook" --worker; then
    exit 1
fi
