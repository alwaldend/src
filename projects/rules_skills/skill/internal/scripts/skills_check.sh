#!/usr/bin/env bash
# Verifies that a skill discovery directory matches its declared skills.
#
# The rule invokes this script with --discovery-dir, --config,
# --workspace-marker, and --marker-suffix. The config file uses the same
# tab-separated data format as the writer. Missing, extra, misdirected, and
# stale entries fail. Written entries must be real directories whose file set
# matches the declaration and whose bytes match their declared sources.
set -euo pipefail

discovery_relative="${SKILLS_WRITE_DISCOVERY_DIR:-}"
config_path="${SKILLS_WRITE_CONFIG:-}"
workspace_marker="${SKILLS_WRITE_WORKSPACE_MARKER:-}"
marker_suffix="${SKILLS_WRITE_MARKER_SUFFIX:-}"
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
    --workspace-marker)
        workspace_marker="$2"
        shift 2
        ;;
    --marker-suffix)
        marker_suffix="$2"
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
if [[ -z "${workspace_marker}" || -z "${marker_suffix}" ]]; then
    echo >&2 "--workspace-marker and --marker-suffix are required"
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

resolve_path() {
    local path="$1"
    local directory name target
    local step
    for ((step = 0; step < 64; step++)); do
        directory="$(cd -P "$(dirname "${path}")" && pwd)" || return 1
        name="$(basename "${path}")"
        path="${directory}/${name}"
        if [[ ! -L "${path}" ]]; then
            printf '%s\n' "${path}"
            return 0
        fi
        target="$(readlink "${path}")" || return 1
        if [[ "${target}" == /* ]]; then
            path="${target}"
        else
            path="${directory}/${target}"
        fi
    done
    echo >&2 "too many symlinks while resolving: $1"
    return 1
}

marker="$(rlocation "${workspace_marker}")"
if [[ -z "${marker}" ]]; then
    echo >&2 "workspace marker is absent from runfiles: ${workspace_marker}"
    exit 1
fi
marker="$(resolve_path "${marker}")" || exit 1
if [[ "${marker}" != *"${marker_suffix}" ]]; then
    echo >&2 "could not locate the workspace from ${marker}"
    exit 1
fi
workspace="${marker%"${marker_suffix}"}"
if [[ ! -f "${workspace}/MODULE.bazel" ]]; then
    echo >&2 "resolved marker is not inside a Bazel workspace: ${marker}"
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
for index in "${!discovery_parts[@]}"; do
    current="${current}/${discovery_parts[$index]}"
    if [[ -L "${current}" ]]; then
        echo >&2 "skill discovery path has a symlink ancestor: ${current}"
        exit 1
    fi
done
if [[ -L "${discovery}" || ! -d "${discovery}" ]]; then
    echo >&2 "skill discovery path is not a real directory: ${discovery}"
    exit 1
fi

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

actual_count=0
shopt -s dotglob nullglob
existing_entries=("${discovery}"/*)
shopt -u dotglob nullglob
for entry in "${existing_entries[@]}"; do
    actual_count=$((actual_count + 1))
    if ! expected_index "${entry##*/}" >/dev/null; then
        echo >&2 "unexpected skill discovery entry: ${entry##*/}"
        exit 1
    fi
done
if [[ "${actual_count}" -ne "${#entry_name[@]}" ]]; then
    echo >&2 "skill discovery has ${actual_count} entries, expected ${#entry_name[@]}"
    exit 1
fi

for index in "${!entry_name[@]}"; do
    entry="${discovery}/${entry_name[$index]}"
    if [[ "${entry_kind[$index]}" == "symlink" ]]; then
        if [[ ! -L "${entry}" ]]; then
            echo >&2 "entry is not a symlink: ${entry}"
            exit 1
        fi
        if [[ "$(readlink "${entry}")" != "${entry_target[$index]}" ]]; then
            echo >&2 "${entry} points to $(readlink "${entry}"), expected ${entry_target[$index]}"
            exit 1
        fi
        canonical="$(resolve_path "${entry}")" || exit 1
        if [[ ! -d "${canonical}" || ! -f "${canonical}/SKILL.md" ]]; then
            echo >&2 "canonical skill root lacks SKILL.md: ${canonical}"
            exit 1
        fi
        continue
    fi
    if [[ -L "${entry}" || ! -d "${entry}" ]]; then
        echo >&2 "written entry is not a real directory: ${entry}"
        exit 1
    fi
    if [[ -n "$(find "${entry}" -type l -print -quit)" ]]; then
        echo >&2 "written entry contains a symlink: ${entry}"
        exit 1
    fi
    actual_files=()
    while IFS= read -r actual_file; do
        actual_files+=("${actual_file#"${entry}/"}")
    done < <(find "${entry}" -type f | sort)
    declared=()
    for position in "${!payload_entry[@]}"; do
        if [[ "${payload_entry[$position]}" == "${index}" ]]; then
            declared+=("${payload_destination[$position]}")
        fi
    done
    if [[ "${#actual_files[@]}" -ne "${#declared[@]}" ]]; then
        echo >&2 "${entry} has ${#actual_files[@]} files, expected ${#declared[@]}"
        exit 1
    fi
    for position in "${!declared[@]}"; do
        if [[ "${actual_files[$position]}" != "${declared[$position]}" ]]; then
            echo >&2 "unexpected file ${actual_files[$position]} in ${entry}"
            exit 1
        fi
    done
    for position in "${!payload_entry[@]}"; do
        if [[ "${payload_entry[$position]}" != "${index}" ]]; then
            continue
        fi
        destination="${entry}/${payload_destination[$position]}"
        if [[ ! -f "${destination}" ]]; then
            echo >&2 "written skill file is missing: ${destination}"
            exit 1
        fi
        source="$(rlocation "${payload_path[$position]}")"
        if [[ -z "${source}" || ! -f "${source}" ]]; then
            echo >&2 "skill file is absent from runfiles: ${payload_path[$position]}"
            exit 1
        fi
        if ! cmp -s "${source}" "${destination}"; then
            echo >&2 "${destination} differs from its declared source"
            exit 1
        fi
    done
done

printf 'verified %s skill discovery entries in %s\n' \
    "${#entry_name[@]}" "${discovery_relative}"
