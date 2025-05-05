#!/usr/bin/env sh
#
# Add commits from a separate repository
# https://plankenau.com/blog/post/copy-commits-separate-git-repos

set -o errexit -o nounset #-o xtrace

echo "Input repo from"
read -r repo_from
if [ ! -d "${repo_from}" ]; then
    echo "invalid directory"
    exit 1
fi
echo "Input repo to"
read -r repo_to
if [ ! -d "${repo_to}" ]; then
    echo "invalid directory"
    exit 1
fi
echo "Input directory to"
read -r repo_to_dir
echo "Is it ok? y/n"
echo "From '${repo_from}' to '${repo_to}' with directory '${repo_to_dir}'"
read -r answer
if [ "${answer}" != "y" ]; then
    echo "Aborted"
    exit 1
fi
patches=$(mktemp -d)
(
    cd "${repo_from}"
    git format-patch -o "${patches}" --root
)
(
    cd "${repo_to}"
    git checkout -B "dev/add-commits-from-$(basename "${repo_from}")"
    git am --directory="${repo_to_dir}" --committer-date-is-author-date --3way "${patches}"/*.patch
)
