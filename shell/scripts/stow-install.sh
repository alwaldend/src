#!/usr/bin/env sh
#
# Install dotfiles

set -o errexit -o nounset -o xtrace

which stow rsync tar >/dev/null

archive="${1:?missing archive location}"
action="${2:?missing action}"
target="${3:-${HOME}}"
stow_dir="${XDG_DATA_HOME:-${HOME}/.local/share}/dotfiles/stow"
temp_dir=$(mktemp -d)
trap 'rm -rf "${temp_dir}"' EXIT
dir="$(dirname "${stow_dir}")"
package="$(basename "${stow_dir}")"

mkdir -p "${stow_dir}"
tar -C "${temp_dir}" --strip-components=1 -xzf "${archive}"

case "${action}" in
--delete) ;;
--stow)
    rsync -a --delete "${temp_dir}/" "${stow_dir}/"
    chmod u+x "${stow_dir}/.local/bin"/*
    ;;
--simulate)
    dir=$(dirname "${temp_dir}")
    package=$(basename "${temp_dir}")
    ;;
*)
    echo "invalid action '${action}'"
    exit 1
    ;;
esac

stow --verbose --no-folding --target "${target}" --dir "${dir}" "${action}" "${package}"

mozilla_dir="${HOME}/.mozilla/firefox"
if [ ! -d "${mozilla_dir}/chrome" ]; then
    exit
fi
find "${mozilla_dir}" -name "*.default-release*" |
    while read -r profile; do
        target="${profile}/chrome"
        mkdir -p "${target}"
        stow --verbose --no-folding --target "${target}" --dir "${mozilla_dir}" "${action}" "chrome"
    done
