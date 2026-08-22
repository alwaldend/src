#!/usr/bin/env bash
set -euo pipefail

action="$1"
inventory="$2"
globals="$3"
cinder_conf="$4"
cinder_volume_conf="$5"
nova_compute_conf="$6"
shift 6

if [[ "${action}" == "install-deps" ]]; then
    exec kolla-ansible install-deps "$@"
fi

state_home="${XDG_STATE_HOME:-${HOME}/.local/state}"
config_dir="${KOLLA_CONFIG_DIR:-${state_home}/alwaldend/openstack/kolla}"

if [[ "${action}" == "init-cloud" ]]; then
    clouds_file="${config_dir}/clouds.yaml"
    if [[ ! -s "${clouds_file}" ]]; then
        echo "run kolla.post_deploy before kolla.init_cloud" >&2
        exit 1
    fi
    export OS_CLIENT_CONFIG_FILE="${clouds_file}"
    for volume_type in fast bulk; do
        if ! openstack \
            --os-cloud kolla-admin \
            volume type show "${volume_type}" >/dev/null 2>&1
        then
            openstack \
                --os-cloud kolla-admin \
                volume type create "${volume_type}"
        fi
        openstack \
            --os-cloud kolla-admin \
            volume type set \
            --property "volume_backend_name=${volume_type}" \
            "${volume_type}"
    done
    exit 0
fi

mkdir -p "${config_dir}"
chmod 0700 "${config_dir}"
if [[ -d "${config_dir}/config" ]]; then
    find "${config_dir}/config" -mindepth 1 -delete
fi
mkdir -p \
    "${config_dir}/config/cinder" \
    "${config_dir}/config/nova"

install -m 0600 "${globals}" "${config_dir}/globals.yml"
install -m 0600 "${cinder_conf}" "${config_dir}/config/cinder.conf"
install -m 0600 \
    "${cinder_volume_conf}" \
    "${config_dir}/config/cinder/cinder-volume.conf"
install -m 0600 \
    "${nova_compute_conf}" \
    "${config_dir}/config/nova/nova-compute.conf"

trap 'rm -f "${config_dir}/passwords.yml"' EXIT
install -m 0600 "${KOLLA_PASSWORDS_FILE}" "${config_dir}/passwords.yml"

kolla-ansible \
    "${action}" \
    --configdir "${config_dir}" \
    --inventory "${inventory}" \
    --become \
    "$@"
