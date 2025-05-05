#!/usr/bin/env sh
set -o errexit -o nounset -o xtrace

session_name=$(tmux display-message -p '#{session_name}')
window_id=$(tmux display-message -p '#{window_id}')
popup_name="${session_name}-popup-${window_id}"
is_popup="${TMUX_IS_POPUP:-}"

if [ "${is_popup}" ]; then
    tmux detach-client
    exit
fi

if ! tmux has-session -t "${popup_name}"; then
    tmux new-session \
        -d \
        -e "TMUX_IS_POPUP=1" \
        -c '#{pane_current_path}' \
        -s "${popup_name}"
    tmux set-option -t "${popup_name}" status off
fi

tmux-display-popup -E "tmux attach-session -t '${popup_name}'"
