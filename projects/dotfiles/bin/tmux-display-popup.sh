#!/usr/bin/env sh
#
# Wrapper for display-popup command

set -o errexit -o nounset -o xtrace

tmux display-popup -x C -y C -w 90% -h 75% "${@}"
