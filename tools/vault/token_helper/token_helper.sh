#!/usr/bin/env sh
set -eu

key="vault_token_helper"
value="vault_token_helper_token"
label="${key}:${value}"

main() {
    action="${1}"
    case "${action}" in
        get)
            echo "Looking up token" >&2
            secret-tool lookup "${key}" "${value}"
            ;;
        store)
            echo "Storing token" >&2
            cat - | secret-tool store --label="${label}" "${key}" "${value}"
            ;;
        erase)
            secret-tool clear "${key}" "${value}"
            ;;
        *)
            echo "Invalid args: ${@}" >&2
            exit 1
            ;;
    esac

}

main "${@}"
