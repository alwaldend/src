#!/usr/bin/env bash

# shellcheck disable=SC1090
. "${DOTFILES_FUNCTIONS_PATH:-"${HOME}/.local/bin/functions"}"

export FZF_DEFAULT_COMMAND="fd --type f"
export FZF_COMPLETION_TRIGGER="**"

GPG_TTY=$(tty || true)
export GPG_TTY
export EDITOR="nvim"
# color fixes
export TERM="screen-256color"
export COLORTERM="truecolor"
# XDG variables
export XDG_CONFIG_HOME="${HOME}/.config"
export XDG_CACHE_HOME="${HOME}/.cache"
export XDG_DATA_HOME="${HOME}/.local/share"
export XDG_STATE_HOME="${HOME}/.local/state"
# nvim config dir
export NVM_DIR="${XDG_CONFIG_HOME}/nvm"
## the main prompt variable
# https://stackoverflow.com/a/28938235
export PS1='\[\e[35;1m\]\u\[\e[m\]@\[\e[35;1m\]\h\[\e[m\] \[\e[36;1m\]\w\[\e[m\] \[\e[1m\]$(__git_ps1 "%s")\[\e[m\]\n'
export PS2=""
export PROMPT_COMMAND="d_prompt_command"
# git prompt variables
export GIT_PS1_SHOWDIRTYSTATE="true"
export GIT_PS1_SHOWUPSTREAM="auto"
export GIT_PS1_STATESEPARATOR=" "
export GIT_PS1_SHOWCOLORHINTS="true"
export GIT_PS1_SHOWCONFLICTSTATE="yes"
export GIT_PS1_SHOWUNTRACKEDFILES="true"

## all locales are overwritten
export LC_ALL="en_US.UTF-8"
export LANG="en_US.UTF-8"
## silences npm funding messages
export OPEN_SOURCE_CONTRIBUTOR="true"
# git editor
export GIT_EDITOR="nvim"

# history file location
export HISTFILE="${HOME}/.bash_history"
# ignore commands starting with whitespace
export HISTCONTROL="ignorespace"
# bash history size (file size)
export HISTFILESIZE=
# bash history size (number of lines)
export HISTSIZE=10000

### proxy
export all_proxy="http://127.0.0.1:10808"
export ALL_PROXY="${all_proxy}"
export http_proxy="${all_proxy}"
export HTTP_PROXY="${all_proxy}"
export https_proxy="${all_proxy}"
export HTTPS_PROXY="${all_proxy}"

# android sdk
export ANDROID_HOME="${HOME}/Android/Sdk"
### path edits
export GOPATH="${HOME}/.go"
if [ -L "/usr/java/latest" ]; then
    export JAVA_HOME="/usr/java/latest"
elif [ -L "/usr/lib/jvm/java" ]; then
    export JAVA_HOME="/usr/lib/jvm/java"
fi

# https://www.gnu.org/software/bash/manual/html_node/The-Shopt-Builtin.html
# vim mode for the terminal
set -o vi
# check the window size after each command and, if necessary,
# update the values of LINES and COLUMNS.
shopt -s checkwinsize
# minor errors in the spelling of a directory in a cd command will be corrected
#shopt -s cdspell
# attempt spelling correction on directory names during word completion
shopt -s dirspell
# include filenames beginning with a ‘.’ in the results of filename expansion
shopt -s dotglob
# send SIGHUP to all jobs when an interactive login shell exits
shopt -s huponexit
# match filenames in a case-insensitive fashion when performing filename expansion.
shopt -s nocaseglob
# flushes the command to the history file immediately
# otherwise, this would happen only when the shell exits
shopt -s histappend
# attempt to save all lines of a multiple-line
# command in the same history entry.
shopt -s cmdhist

d_add_to_path_front \
    ~/.local/share/gem/ruby/*/bin \
    "${HOME}/.local/bin" \
    /usr/bin \
    "${GOPATH}/bin" \
    "${HOME}/.cargo/bin" \
    "${HOME}/yandex-cloud/bin" \
    "${ANDROID_HOME}/tools" \
    "${ANDROID_HOME}/tools/bin" \
    "${ANDROID_HOME}/cmdline-tools/latest/bin" \
    "${ANDROID_HOME}/platform-tools" \
    /usr/java/latest/bin \
    /usr/mvn/latest/bin \
    /opt/nvim-linux-x86_64/bin

d_source_scripts \
    /etc/profile.d/bash_completion.sh \
    /usr/share/git-core/contrib/completion/git-prompt.sh \
    /usr/share/bash-completion/completions/git \
    /usr/share/fzf/shell/key-bindings.bash \
    "${NVM_DIR}/nvm.sh" \
    "${NVM_DIR}/bash_completion"

if [ "${#}" -gt 0 ]; then
    set -o errexit -o nounset -o xtrace
    "${@}"
fi
