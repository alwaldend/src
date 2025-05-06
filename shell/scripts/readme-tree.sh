#!/usr/bin/env sh
set -o errexit -o nounset -o xtrace

if [ -n "${1:-}" ]; then
    cd "${1}"
fi
files=$(git ls-files | grep README.md)
for file in ${files}; do
    if [ "${file}" = "README.md" ]; then
       continue
    fi
    title=$(awk '/title: / { $1=""; print $0 }' "${file}")
    if [ -n "${title}" ]; then
       title=":${title}"
    fi
    echo "$(dirname ${file})" | \
        awk -v title="${title}" -F "/" '
            $0 != "." {
                for (i = 0; i < NF - 1; i++) {
                    printf "  "
                }
                printf "- " "[" $NF "]" "(" $0 ")" title
                printf "\n"
            }
        '
done
