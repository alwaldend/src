#!/usr/bin/env bash
set -euo pipefail

command="${1:?missing Promptfoo command}"
shift

config=""
output=""
no_cache=0
no_write=0
no_share=0
max_concurrency=""
while (($#)); do
  case "$1" in
    --config)
      config="${2:?missing config path}"
      shift 2
      ;;
    --output)
      output="${2:?missing output path}"
      shift 2
      ;;
    --no-cache)
      no_cache=1
      shift
      ;;
    --no-write)
      no_write=1
      shift
      ;;
    --no-share)
      no_share=1
      shift
      ;;
    --max-concurrency)
      max_concurrency="${2:?missing max concurrency}"
      shift 2
      ;;
    *)
      shift
      ;;
  esac
done

test -f "$config"
test "${PROMPTFOO_CONFIG_DIR}" = \
  "${PROMPTFOO_STATE_DIR}/config"
test "${PROMPTFOO_CACHE_PATH}" = \
  "${PROMPTFOO_STATE_DIR}/cache"
test "${PROMPTFOO_SUBJECT_CODEX_HOME}" = \
  "${PROMPTFOO_STATE_DIR}/subject_codex"
test "${PROMPTFOO_JUDGE_CODEX_HOME}" = \
  "${PROMPTFOO_STATE_DIR}/judge_codex"
test "${CODEX_HOME}" = "${PROMPTFOO_SUBJECT_CODEX_HOME}"
if [[ "${EXPECT_REUSED_CODEX_LOGIN:-0}" == 1 ]]; then
  test "$max_concurrency" = 1
  test "${PROMPTFOO_ASSERTIONS_MAX_CONCURRENCY}" = 1
  test "${PROMPTFOO_SUGGESTIONS_MAX_CONCURRENCY}" = 1
  auth_target="$(readlink -f "${CODEX_HOME}/auth.json")"
  if flock -n "$(dirname "$auth_target")" true; then
    echo >&2 "reused Codex auth was not locked"
    exit 1
  fi
  for isolated_codex_home in \
    "${PROMPTFOO_SUBJECT_CODEX_HOME}" \
    "${PROMPTFOO_JUDGE_CODEX_HOME}"; do
    test -f "${isolated_codex_home}/auth.json"
    test -L "${isolated_codex_home}/auth.json"
    test "$(sed -n '1p' "${isolated_codex_home}/auth.json")" = \
      '{"fixture":true}'
  done
  printf '%s\n' '{"fixture":"refreshed"}' > \
    "${PROMPTFOO_SUBJECT_CODEX_HOME}/auth.json"
  test "$(sed -n '1p' "${PROMPTFOO_JUDGE_CODEX_HOME}/auth.json")" = \
    '{"fixture":"refreshed"}'
  test -L "${PROMPTFOO_SUBJECT_CODEX_HOME}/auth.json"
  test -L "${PROMPTFOO_JUDGE_CODEX_HOME}/auth.json"
fi
test "${PROMPTFOO_JUDGE_WORKSPACE}" = \
  "${PROMPTFOO_STATE_DIR}/judge_workspace"
test "${PROMPTFOO_SKILL_WORKSPACE}" = \
  "${PROMPTFOO_STATE_DIR}/workspace"
scratch_root="$(cd "${TEST_TMPDIR:?}" && pwd -P)"
case "${PROMPTFOO_STATE_DIR}" in
  "${scratch_root}"/rules_promptfoo.*) ;;
  *) exit 1 ;;
esac
test "${TMPDIR}" = "${PROMPTFOO_STATE_DIR}/tmp"
test "${TMP}" = "${TMPDIR}"
test "${TEMP}" = "${TMPDIR}"
test -d "${TMPDIR}"
test "$(stat -c '%a' "${TMPDIR}")" = 700
test "${PWD}" = "${PROMPTFOO_SKILL_WORKSPACE}"
test -d "${PROMPTFOO_JUDGE_WORKSPACE}"

skill="${PROMPTFOO_SKILL_WORKSPACE}/.agents/skills/fixture_skill/SKILL.md"
reference="${PROMPTFOO_SKILL_WORKSPACE}/.agents/skills/fixture_skill/reference.txt"
test -f "$skill"
test -f "$reference"
test ! -L "$skill"
test ! -L "$reference"
test ! -e "${PROMPTFOO_JUDGE_WORKSPACE}/.agents"
if [[ -n "${RECORD_STATE_FILE:-}" ]]; then
  printf '%s\n' "${PROMPTFOO_STATE_DIR}" >"${RECORD_STATE_FILE}"
fi

case "$command" in
  validate)
    test -z "$output"
    ;;
  eval)
    test "$no_cache" = 1
    test "$no_write" = 1
    test "$no_share" = 1
    test -n "$output"
    printf '%s\n' '{"ok":true}' >"$output"
    ;;
  *)
    echo >&2 "unexpected command: $command"
    exit 1
    ;;
esac

if [[ -n "${FAKE_PROMPTFOO_EXIT_CODE:-}" ]]; then
  exit "${FAKE_PROMPTFOO_EXIT_CODE}"
fi
