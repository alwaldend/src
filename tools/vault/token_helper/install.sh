#!/usr/bin/env sh
set -eu

echo "token_helper = \"${BUILD_WORKSPACE_DIRECTORY}/tools/vault/token_helper/token_helper.sh\"" >>"${HOME}/.vault"
