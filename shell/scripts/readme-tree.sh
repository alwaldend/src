#!/usr/bin/env sh
set -o errexit -o nounset -o xtrace

if [ -n "${1:-}" ]; then
    cd "${1}"
fi

files=$(git ls-files | grep README.md)
echo "<table>"
for file in ${files}; do
    if [ "${file}" = "README.md" ]; then
       continue
    fi
    title=$(awk '/title: / { $1=""; print $0 }' "${file}")
    if [ -n "${title}" ]; then
       title="${title}"
    fi
    echo "$(dirname ${file})" | \
        awk -v title="${title}" -F "/" '
            $0 != "." {
                if (NF == 1) {
                    printf "<tr>"
                    printf "<td><a href=./" $0 ">" $NF "</a></td>"
                    printf "<td>" title "</td>"
                    printf "</tr>\n"
                } else if (NF == 2) {
                    printf "<tr>"
                    printf "<td></td>"
                    printf "<td><a href=./" $0 ">" $NF "</a></td>"
                    printf "<td>" title "</td>"
                    printf "</tr>\n"
                } else {
                    exit 1
                }
            }
        '
done
echo "</table>"
