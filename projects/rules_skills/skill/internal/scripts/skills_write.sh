#!/usr/bin/env bash
# Reconciles a skill discovery directory from declared skill libraries.
#
# The rule invokes this script with --discovery-dir and --config; both default
# to the SKILLS_WRITE_DISCOVERY_DIR and SKILLS_WRITE_CONFIG environment
# variables so the same script serves `bazel run` and test targets. The config
# file is tab-separated data, not code:
#
#   entries
#   <kind>\t<name>\t<symlink target>
#   payloads
#   <entry index>\t<destination>\t<runfiles path>
set -euo pipefail

discovery_relative="${SKILLS_WRITE_DISCOVERY_DIR:-}"
config_path="${SKILLS_WRITE_CONFIG:-}"
while [[ $# -gt 0 ]]; do
    case "$1" in
    --discovery-dir)
        discovery_relative="$2"
        shift 2
        ;;
    --config)
        config_path="$2"
        shift 2
        ;;
    *)
        echo >&2 "unknown option: $1"
        exit 2
        ;;
    esac
done
if [[ -z "${discovery_relative}" || -z "${config_path}" ]]; then
    echo >&2 "--discovery-dir and --config are required"
    exit 2
fi

f=bazel_tools/tools/bash/runfiles/runfiles.bash
# shellcheck disable=SC1090
source "${RUNFILES_DIR:-/dev/null}/$f" 2>/dev/null ||
    source "$(grep -sm1 "^$f " \
        "${RUNFILES_MANIFEST_FILE:-/dev/null}" | cut -f2- -d' ')" \
        2>/dev/null ||
    source "$0.runfiles/$f" 2>/dev/null ||
    source "$(grep -sm1 "^$f " "$0.runfiles_manifest" |
        cut -f2- -d' ')" 2>/dev/null || {
    echo >&2 "ERROR: cannot find $f"
    exit 1
}
runfiles_export_envvars

workspace="${BUILD_WORKSPACE_DIRECTORY:-}"
if [[ -z "${workspace}" || "${workspace}" != /* ]]; then
    echo >&2 "BUILD_WORKSPACE_DIRECTORY must be an absolute path"
    exit 1
fi
if [[ ! -f "${workspace}/MODULE.bazel" ]]; then
    echo >&2 "BUILD_WORKSPACE_DIRECTORY is not a Bazel workspace: ${workspace}"
    exit 1
fi

config_file="$(rlocation "${config_path}")"
if [[ -z "${config_file}" || ! -f "${config_file}" ]]; then
    echo >&2 "skill config is absent from runfiles: ${config_path}"
    exit 1
fi

discovery="${workspace}/${discovery_relative}"
current="${workspace}"
IFS='/' read -r -a discovery_parts <<<"${discovery_relative}"
for ((index = 0; index < ${#discovery_parts[@]} - 1; index++)); do
    current="${current}/${discovery_parts[$index]}"
    if [[ -L "${current}" ]]; then
        echo >&2 "discovery parent must not be a symlink: ${current}"
        exit 1
    fi
    if [[ -e "${current}" && ! -d "${current}" ]]; then
        echo >&2 "discovery parent must be a directory: ${current}"
        exit 1
    fi
    if [[ ! -e "${current}" ]]; then
        mkdir "${current}"
    fi
done
if [[ -L "${discovery}" ]]; then
    echo >&2 "discovery directory must not be a symlink: ${discovery}"
    exit 1
fi
if [[ -e "${discovery}" && ! -d "${discovery}" ]]; then
    echo >&2 "discovery path must be a directory: ${discovery}"
    exit 1
fi
if [[ ! -e "${discovery}" ]]; then
    mkdir "${discovery}"
fi

lock="${discovery}.lock"
if ! mkdir "${lock}" 2>/dev/null; then
    echo >&2 "another skills updater holds ${lock}"
    exit 1
fi
trap 'rmdir "${lock}"' EXIT

entry_kind=()
entry_name=()
entry_target=()
payload_entry=()
payload_destination=()
payload_path=()
section="entries"
while IFS=$'\t' read -r first second third; do
    if [[ "${first}" == "payloads" && -z "${second}" ]]; then
        section="payloads"
        continue
    fi
    [[ -z "${first}" ]] && continue
    if [[ "${section}" == "entries" ]]; then
        entry_kind+=("${first}")
        entry_name+=("${second}")
        entry_target+=("${third}")
    else
        payload_entry+=("${first}")
        payload_destination+=("${second}")
        payload_path+=("${third}")
    fi
done <"${config_file}"
if [[ "${#entry_name[@]}" -eq 0 ]]; then
    echo >&2 "skill config declares no discovery entries"
    exit 1
fi

expected_index() {
    local candidate="$1"
    local index
    for index in "${!entry_name[@]}"; do
        if [[ "${entry_name[$index]}" == "${candidate}" ]]; then
            printf '%s\n' "${index}"
            return 0
        fi
    done
    return 1
}

shopt -s dotglob nullglob
existing_entries=("${discovery}"/*)
shopt -u dotglob nullglob
for entry in "${existing_entries[@]}"; do
    if ! expected_index "${entry##*/}" >/dev/null; then
        rm -rf "${entry}"
    fi
done

for index in "${!entry_name[@]}"; do
    entry="${discovery}/${entry_name[$index]}"
    if [[ "${entry_kind[$index]}" == "symlink" ]]; then
        expected="${entry_target[$index]}"
        if [[ -L "${entry}" &&
            "$(readlink "${entry}")" == "${expected}" ]]; then
            continue
        fi
        rm -rf "${entry}"
        ln -s "${expected}" "${entry}"
        continue
    fi
    rm -rf "${entry}"
    mkdir -p "${entry}"
    for position in "${!payload_entry[@]}"; do
        if [[ "${payload_entry[$position]}" != "${index}" ]]; then
            continue
        fi
        source="$(rlocation "${payload_path[$position]}")"
        if [[ -z "${source}" || ! -f "${source}" ]]; then
            echo >&2 "skill file is absent from runfiles: ${payload_path[$position]}"
            exit 1
        fi
        destination="${entry}/${payload_destination[$position]}"
        mkdir -p "$(dirname "${destination}")"
        cp -f "${source}" "${destination}"
        chmod 0644 "${destination}"
    done
done

printf 'reconciled %s skill discovery entries in %s\n' \
    "${#entry_name[@]}" "${discovery_relative}"
