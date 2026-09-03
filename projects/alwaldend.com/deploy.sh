#!/usr/bin/env sh

set -eu

archive=${1:?usage: deploy.sh SITE_ARCHIVE}
archive=$(realpath "${archive}")
temp=$(mktemp -d)
trap 'rm -rf "${temp}"' EXIT
cd "${temp}"

git clone \
    --depth 1 \
    --single-branch \
    --branch pages \
    git@github.com:alwaldend/alwaldend.github.io.git site
cd site

# Remove the previous checkout (except .git) so deleted or renamed pages
# are not served stale; tar -xf does not remove absent files.
find . -mindepth 1 -maxdepth 1 ! -name .git -exec rm -rf {} +
tar -xf "${archive}" --strip-components 1

# Disable Jekyll so GitHub Pages serves the Hugo output verbatim.
touch .nojekyll
git add -A

if git diff --cached --quiet; then
    echo "Site unchanged, nothing to deploy"
    exit 0
fi

git commit -m "Update the site"
git push origin HEAD:pages
