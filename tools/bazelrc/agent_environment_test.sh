#!/usr/bin/env bash
set -euo pipefail

project_bazelrc="$1"

if grep -Eq -- '--(repo_env|action_env|host_action_env|test_env)=(TEMP|TMP|TMPDIR)([ =]|$)' "${project_bazelrc}"; then
    echo >&2 "agent profile must not propagate ambient temporary directories"
    exit 1
fi
