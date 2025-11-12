#!/usr/bin/env sh

session_name="DEFAULT"
if [ "${TERM_PROGRAM}" = "tmux" ]; then
    exit
fi
if ! tmux has-session -t "${session_name}"; then
    tmux new-session -d -s "${session_name}"
    tmux send-keys tmux-new-window Enter
fi
tmux attach-session -t "${session_name}"
