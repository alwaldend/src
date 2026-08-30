#!/usr/bin/env bash
set -euo pipefail

mmdc="$1"
input="$2"
output="${TEST_TMPDIR}/smoke.svg"

"${mmdc}" -i "${input}" -o "${output}"
test -s "${output}"
grep -q '<svg' "${output}"
