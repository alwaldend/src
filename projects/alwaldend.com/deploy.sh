#!/usr/bin/env sh

set -eu

archive=${1:?usage: deploy.sh SITE_ARCHIVE REPOSITORY BRANCH}
repository=${2:?usage: deploy.sh SITE_ARCHIVE REPOSITORY BRANCH}
branch=${3:?usage: deploy.sh SITE_ARCHIVE REPOSITORY BRANCH}
archive=$(realpath "${archive}")
temp=$(mktemp -d)
trap 'rm -rf "${temp}"' EXIT
cd "${temp}"

# A repository that has never been published has no target branch yet, so the
# first deployment creates it as an orphan instead of cloning it.
if git ls-remote --exit-code --heads "git@github.com:${repository}.git" "${branch}" >/dev/null 2>&1; then
    git clone \
        --depth 1 \
        --single-branch \
        --branch "${branch}" \
        "git@github.com:${repository}.git" site
else
    echo "Branch ${branch} does not exist in ${repository}; creating it"
    git init --quiet site
    git -C site remote add origin "git@github.com:${repository}.git"
    git -C site checkout --quiet --orphan "${branch}"
fi
cd site

# Remove the previous checkout (except .git) so deleted or renamed pages
# are not served stale; tar -xf does not remove absent files.
find . -mindepth 1 -maxdepth 1 ! -name .git -exec rm -rf {} +
tar -xf "${archive}" --strip-components 1
# CNAME is at the archive root, so the stripped extraction skips it.
tar -xf "${archive}" CNAME

# Disable Jekyll so GitHub Pages serves the Hugo output verbatim.
touch .nojekyll
git add -A

if git diff --cached --quiet; then
    echo "Site unchanged, nothing to deploy"
    exit 0
fi

git -c user.name=alwaldend -c user.email=alwaldend@alwaldend.com commit -m "Update the site"
git push origin "HEAD:${branch}"
