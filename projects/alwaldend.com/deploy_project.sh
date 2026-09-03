#!/usr/bin/env sh

set -eu

archive=${1:?usage: deploy_project.sh SITE_ARCHIVE [PAGES_REPO]}
archive=$(realpath "${archive}")
repo=${2:-$(basename "${archive%.tar}")-pages}

temp=$(mktemp -d)
trap 'rm -rf "${temp}"' EXIT
cd "${temp}"

git clone \
    --depth 1 \
    --single-branch \
    --branch pages \
    "git@github.com:alwaldend/${repo}.git" site
cd site

find . -mindepth 1 -maxdepth 1 ! -name .git -exec rm -rf {} +
tar -xf "${archive}" --strip-components 1
touch .nojekyll
git add -A

if git diff --cached --quiet; then
    echo "${repo}: unchanged"
    exit 0
fi

git commit -m "Update ${repo}"
git push origin HEAD:pages
