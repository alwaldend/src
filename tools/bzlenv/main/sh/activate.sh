#!/usr/bin/env sh

bzlenv_deactivate() {
    if [ -n "${_BZLENV_OLD_BAZEL_BINDIR:-}" ]; then
        export BAZEL_BINDIR="${_BZLENV_OLD_BAZEL_BINDIR}"
        unset _BZLENV_OLD_BAZEL_BINDIR
    fi
    if [ -n "${_BZLENV_OLD_PATH:-}" ]; then
        export PATH="${_BZLENV_OLD_PATH}"
        unset _BZLENV_OLD_PATH
    fi
    if [ -n "${_BZLENV_OLD_PS1:-}" ]; then
        export PS1="${_BZLENV_OLD_PS1}"
        unset _BZLENV_OLD_PS1
    fi
    if [ -n "${_BZLENV_OLD_RUNFILES_DIR:-}" ]; then
        export RUNFILES_DIR="${_BZLENV_OLD_RUNFILES_DIR}"
    else
        unset RUNFILES_DIR
    fi
    unset _BZLENV_OLD_RUNFILES_DIR

    # Call hash to forget past locations. Without forgetting
    # past locations the $PATH changes we made may not be respected.
    # See "man bash" for more details. hash is usually a builtin of your shell
    hash -r 2>/dev/null

    if [ ! "${1:-}" = "nondestructive" ]; then
        # Self destruct
        unset -f bzlenv_deactivate
    fi
}

# Unset irrelevant variables
bzlenv_deactivate nondestructive

_BZLENV_WORKSPACE_DIR="${BUILD_WORKSPACE_DIRECTORY}"
_BZLENV_OLD_BAZEL_BINDIR="${BAZEL_BINDIR:-}"
export BAZEL_BINDIR="." # for js rules
_BZLENV_OLD_RUNFILES_DIR="${RUNFILES_DIR:-}"
export RUNFILES_DIR="${BZLENV_RUNFILES_DIR}" # for sh runfiles lib
_BZLENV_OLD_PATH="${PATH:-}"
export PATH="${PATH:-}:${BZLENV_PATH}"
_BZLENV_OLD_PS1="${PS1:-}"
export PS1="(.bzlenv) ${PS1:-}"

if [ -f "${_BZLENV_WORKSPACE_DIR}/.env" ]; then
    set -a
    # shellcheck disable=SC1091
    . "${_BZLENV_WORKSPACE_DIR}/.env"
    set +a
fi
