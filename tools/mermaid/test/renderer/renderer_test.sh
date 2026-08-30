#!/usr/bin/env bash
set -euo pipefail

svg="$1"

test -s "${svg}"
awk '
    index($0, "<svg") {
        found = 1
        exit
    }
    END {
        if (!found) {
            exit 1
        }
    }
' "${svg}"
