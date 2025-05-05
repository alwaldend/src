#!/usr/bin/env sh
set -o errexit -o nounset -o xtrace

selected_config=$(find ~/.config/git/ -name "*.gitconfig" | fzf --prompt "Select config")
target_config="${HOME}/.config/git/config"
echo "
[include]
    path = ${selected_config}
" >"${target_config}"
