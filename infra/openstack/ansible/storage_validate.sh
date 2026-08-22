#!/usr/bin/env bash
set -euo pipefail

array="$1"
expected_level="$2"
volume_group="$3"
shift 3
expected_devices=("$@")

fail() {
    echo "storage validation failed: $*" >&2
    exit 1
}

array_real="$(readlink -f "${array}")"
[[ -b "${array_real}" ]] || fail "${array} is not a block device"

detail="$(mdadm --detail --export "${array}")"
actual_level="$(
    awk -F= '$1 == "MD_LEVEL" { print $2 }' <<<"${detail}"
)"
actual_count="$(
    awk -F= '$1 == "MD_DEVICES" { print $2 }' <<<"${detail}"
)"

[[ "${actual_level}" == "${expected_level}" ]] ||
    fail "${array} is ${actual_level}, expected ${expected_level}"
[[ "${actual_count}" == "${#expected_devices[@]}" ]] ||
    fail "${array} has ${actual_count} devices, expected " \
        "${#expected_devices[@]}"

array_name="$(basename "${array_real}")"
shopt -s nullglob
actual_slaves=(/sys/class/block/"${array_name}"/slaves/*)
(( ${#actual_slaves[@]} > 0 )) || fail "${array} has no active members"

mapfile -t actual_devices < <(
    for device in "${actual_slaves[@]}"; do
        readlink -f "/dev/$(basename "${device}")"
    done | sort
)
mapfile -t resolved_expected_devices < <(
    for device in "${expected_devices[@]}"; do
        readlink -f "${device}"
    done | sort
)

[[ "${actual_devices[*]}" == "${resolved_expected_devices[*]}" ]] ||
    fail "${array} members do not match the declared devices"

if [[ "${volume_group}" != "-" ]]; then
    mapfile -t volume_group_devices < <(
        pvs --noheadings --options pv_name \
            --select "vg_name=${volume_group}" |
            awk '{$1=$1; print}'
    )
    (( ${#volume_group_devices[@]} == 1 )) ||
        fail "${volume_group} must contain exactly one physical volume"
    volume_group_device="$(readlink -f "${volume_group_devices[0]}")"
    [[ "${volume_group_device}" == "${array_real}" ]] ||
        fail "${volume_group} is not backed by ${array}"
fi
