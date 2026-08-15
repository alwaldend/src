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

mkdir -p \
    "${config_dir}/config/cinder" \
    "${config_dir}/config/nova"
chmod 0700 "${config_dir}"

install -m 0600 "${globals}" "${config_dir}/globals.yml"
install -m 0600 "${cinder_conf}" "${config_dir}/config/cinder.conf"
install -m 0600 \
    "${cinder_volume_conf}" \
    "${config_dir}/config/cinder/cinder-volume.conf"
install -m 0600 \
    "${nova_compute_conf}" \
    "${config_dir}/config/nova/nova-compute.conf"
install -m 0600 "${KOLLA_PASSWORDS_FILE}" "${config_dir}/passwords.yml"

trap 'rm -f "${config_dir}/passwords.yml"' EXIT

kolla-ansible \
    "${action}" \
    --configdir "${config_dir}" \
    --inventory "${inventory}" \
    --become \
    "$@"
