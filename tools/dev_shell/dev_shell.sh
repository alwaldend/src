#!/usr/bin/env sh

set -eu

export RUNFILES_DIR="${0}.runfiles"
bin_path="${RUNFILES_DIR}/_main/bin"
if [ ! -d "${bin_path}" ]; then
    echo "Missing bin path for some reason: ${bin_path}"
    exit 1
fi
echo "Adding tools '$(ls "${bin_path}")' from '${bin_path}' to \${PATH}"
export PATH="${PATH}:${bin_path}"

echo "Entering workspace dir: ${BUILD_WORKSPACE_DIRECTORY}"
cd "${BUILD_WORKSPACE_DIRECTORY}"

if [ -f .env ]; then
    echo "Sourcing .env"
    set -a
    # shellcheck disable=SC1091
    . "${PWD}/.env"
    set +a
fi

echo "Entering dev shell"
exec "${SHELL}" "${@}"
