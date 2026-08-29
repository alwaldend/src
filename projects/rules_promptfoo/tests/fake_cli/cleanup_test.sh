#!/usr/bin/env bash
set -euo pipefail
umask 077

f=bazel_tools/tools/bash/runfiles/runfiles.bash
# shellcheck disable=SC1090
source "${RUNFILES_DIR:-/dev/null}/$f" 2>/dev/null || \
  source "$(grep -sm1 "^$f " \
    "${RUNFILES_MANIFEST_FILE:-/dev/null}" | cut -f2- -d' ')" \
    2>/dev/null || \
  source "$0.runfiles/$f" 2>/dev/null || \
  source "$(grep -sm1 "^$f " "$0.runfiles_manifest" | \
    cut -f2- -d' ')" 2>/dev/null || \
  source "$(grep -sm1 "^$f " "$0.exe.runfiles_manifest" | \
    cut -f2- -d' ')" 2>/dev/null || \
  { echo >&2 "ERROR: cannot find $f"; exit 1; }
runfiles_export_envvars
root="$(runfiles_current_repository || true)"
if [[ -z "$root" ]]; then
  root="_main"
fi
launcher="$(rlocation "$root/tests/fake_cli/launcher_test.sh")"
login_launcher="$(rlocation \
  "$root/tests/fake_cli/external_codex_login_launcher.sh")"
test -x "$launcher"
test -x "$login_launcher"

success_record="${TEST_TMPDIR:?}/success.state"
RECORD_STATE_FILE="$success_record" "$launcher"
success_state="$(sed -n '1p' "$success_record")"
test -n "$success_state"
test ! -e "$success_state"

failure_record="${TEST_TMPDIR}/failure.state"
set +e
FAKE_PROMPTFOO_EXIT_CODE=42 \
  RECORD_STATE_FILE="$failure_record" \
  "$launcher"
status=$?
set -e
test "$status" = 42
failure_state="$(sed -n '1p' "$failure_record")"
test -n "$failure_state"
test ! -e "$failure_state"

unsafe_parent="${TEST_TMPDIR}/unsafe"
unsafe_scratch="${unsafe_parent}/scratch"
mkdir -p "${unsafe_parent}/.agents" "$unsafe_scratch"
set +e
TEST_TMPDIR="$unsafe_scratch" "$launcher"
status=$?
set -e
test "$status" = 2

set +e
env -u TEST_TMPDIR "$launcher"
status=$?
set -e
test "$status" != 0

set +e
TEST_TMPDIR=relative "$launcher"
status=$?
set -e
test "$status" = 2

host_codex_home="${TEST_TMPDIR}/host_codex"
mkdir -p "$host_codex_home"
printf '%s\n' '{"fixture":true}' >"$host_codex_home/auth.json"
chmod 600 "$host_codex_home/auth.json"
login_record="${TEST_TMPDIR}/login.state"
CODEX_HOME="$host_codex_home" \
  EXPECT_REUSED_CODEX_LOGIN=1 \
  RECORD_STATE_FILE="$login_record" \
  "$login_launcher"
login_state="$(sed -n '1p' "$login_record")"
test -n "$login_state"
test ! -e "$login_state"
test "$(sed -n '1p' "$host_codex_home/auth.json")" = \
  '{"fixture":"refreshed"}'
