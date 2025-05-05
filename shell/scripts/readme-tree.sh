#!/usr/bin/env sh
set -o errexit -o nounset -o xtrace

echo "$(tree -aifndI ".git" --noreport --filesfirst --sort name --gitignore . "${@}")" |
    awk -F " ->" '{ printf $1 "\n" }' | 
    awk -F "/" '
        $0 != "." {
            for (i = 0; i < NF - 2; i++) {
                printf "  "
            }
            printf "- " "[" $NF "]" "(" $0 ")"
            printf "\n"
        }
    '
