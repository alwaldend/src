#!/usr/bin/env sh

set -eu

if [ "$#" -ne 1 ]; then
    echo "usage: install.sh BAZEL_AGENT_BINARY" >&2
    exit 2
fi
if [ -z "${HOME:-}" ]; then
    echo "HOME is unset" >&2
    exit 1
fi

source_path=$1
target_directory="${HOME}/.local/bin"
target_path="${target_directory}/bazel_agent"
temporary_path="${target_path}.tmp.$$"

cleanup() {
    rm -f "${temporary_path}"
}
trap cleanup 0
trap 'exit 1' 1 2 15

mkdir -p "${target_directory}"
cp "${source_path}" "${temporary_path}"
chmod 0755 "${temporary_path}"
mv -f "${temporary_path}" "${target_path}"
trap - 0 1 2 15

echo "installed ${target_path}"
