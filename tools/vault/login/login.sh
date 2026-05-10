#!/usr/bin/env sh
set -eu

name="${1:-}"
if [ -z "${name}" ]; then
    echo "Missing name" >&2
    exit 1
fi
shift

url=$(p11tool --only-urls --list-tokens | grep URL | tail -n 1 | awk '{ print $2 }')
if [ -z "${url}" ]; then
    echo "Could not get url" >&2
    exit 1
fi

data=$(\
    curl \
        --request POST \
        --cacert "${VAULT_CACERT}" \
        --cert "${url}" \
        --data "{ \"name\": \"${name}\" }" \
        "${VAULT_ADDR}/v1/auth/cert/login"
)
token=$(echo "${data}" | jq -r .auth.client_token)
echo "${token}" >~/.vault-token
