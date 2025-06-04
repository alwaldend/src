#!/usr/bin/env sh
set -o errexit -o nounset -o xtrace

if [ -n "${TMUX_IS_POPUP:-}" ]; then
    tmux kill-window
    exit
fi

if ! which fzf >/dev/null; then
    exit
fi

choice=$(
    {
        echo ssh
        echo ~
        # shellcheck disable=SC2016
        find ~/repos -maxdepth 6 -type d -name .git -exec dirname "{}" ";" 2>/dev/null
    } | fzf
)

window_name=
window_command=

case "${choice}" in
ssh)
    server=$(awk '/^Host/ { print($2) }' <"${HOME}/.ssh/config" | fzf)
    dirs=$(
        ssh "${server}" -- \
            echo "\${HOME}" "&&" \
            find /etc/ /opt/ /var/ "\${HOME}/repos" -maxdepth 6 -type d 2>/dev/null
    )
    dir=$(echo "${dirs}" | fzf --prompt "ssh dir")
    window_name="${server}//$(basename "${dir}")"
    window_command="ssh -t '${server}' -- cd '${dir}'; exec \"\${SHELL}\" -l"
    ;;
*)
    window_name="//$(basename "${choice}")"
    window_command="cd '${choice}'"
    ;;
esac

tmux rename-window "${window_name}"
tmux send-keys "${window_command}" Enter
