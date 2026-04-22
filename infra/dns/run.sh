#!/usr/bin/env bash
set -eu

function main() {
    vault="tools/vault/vault.script.sh"
    dnscontrol="tools/dnscontrol/dnscontrol.script.sh"
    CLOUDFLAREAPI_ACCOUNT_ID=$("${vault}" kv get --format json sre/data/infra/dns | jq -r .data.data.cloudflare_account_id)
    CLOUDFLAREAPI_API_TOKEN=$("${vault}" kv get --format json sre/data/infra/dns | jq -r .data.data.cloudflare_api_token)
    export CLOUDFLAREAPI_ACCOUNT_ID CLOUDFLAREAPI_API_TOKEN
    trap "rm -f infra/dns/creds.json" EXIT
    envsubst <"infra/dns/creds.json.tpl" >"infra/dns/creds.json"
    "${dnscontrol}" "${@}"
}

main "${@}"
