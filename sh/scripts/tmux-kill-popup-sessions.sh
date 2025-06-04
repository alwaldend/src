#!/usr/bin/env sh
set -o errexit -o nounset -o xtrace

if [ "${TMUX_IS_POPUP:-}" ]; then
    exit
fi

session_name="${1:?missing session name}"
window="${2:-"*"}"

tmux list-sessions -F '#{session_name}' |
    awk -v session_name="^${session_name}-popup-${window}" '$0 ~ session_name' |
    while read -r popup_name; do
        tmux kill-session -C -t "${popup_name}"
    done
