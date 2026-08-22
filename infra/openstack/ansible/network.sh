#!/usr/bin/env bash
set -euo pipefail

uplink="$1"
bridge="$2"
provider="$3"
provider_peer="$4"
management_cidr="$5"
gateway="$6"
shift 6
dns_servers="$(IFS=,; echo "$*")"

management_connection="openstack-management"
uplink_connection="openstack-uplink"
provider_connection="openstack-provider"
provider_peer_connection="openstack-provider-port"

network_is_ready() {
    ip -4 -o address show dev "${bridge}" 2>/dev/null |
        awk '{ print $4 }' | grep -Fxq "${management_cidr}" || return 1
    ip -o link show "${uplink}" 2>/dev/null |
        grep -Fq "master ${bridge}" || return 1
    ip -o link show "${provider_peer}" 2>/dev/null |
        grep -Fq "master ${bridge}" || return 1
    ip link show "${provider}" >/dev/null 2>&1 || return 1
    ip -4 route show default |
        grep -Fq "default via ${gateway} dev ${bridge}" || return 1
    ping -c 1 -W 2 -I "${bridge}" "${gateway}" >/dev/null 2>&1
}

if network_is_ready; then
    echo unchanged
    exit 0
fi

old_connection="$(
    nmcli --get-values GENERAL.CONNECTION device show "${uplink}" |
        head -n 1
)"

for connection in \
    "${provider_peer_connection}" \
    "${provider_connection}" \
    "${uplink_connection}" \
    "${management_connection}"
do
    nmcli connection delete "${connection}" >/dev/null 2>&1 || true
done

nmcli connection add \
    type bridge \
    ifname "${bridge}" \
    con-name "${management_connection}" \
    connection.autoconnect no \
    bridge.stp no \
    ipv4.method manual \
    ipv4.addresses "${management_cidr}" \
    ipv4.gateway "${gateway}" \
    ipv4.dns "${dns_servers}" \
    ipv6.method disabled

nmcli connection add \
    type ethernet \
    ifname "${uplink}" \
    con-name "${uplink_connection}" \
    master "${bridge}" \
    slave-type bridge \
    connection.autoconnect no \
    ipv4.method disabled \
    ipv6.method disabled

nmcli connection add \
    type veth \
    ifname "${provider}" \
    con-name "${provider_connection}" \
    veth.peer "${provider_peer}" \
    connection.autoconnect no \
    ipv4.method disabled \
    ipv6.method disabled

nmcli connection add \
    type ethernet \
    ifname "${provider_peer}" \
    con-name "${provider_peer_connection}" \
    master "${bridge}" \
    slave-type bridge \
    connection.autoconnect no \
    ipv4.method disabled \
    ipv6.method disabled

nmcli connection up "${provider_connection}"
nmcli connection up "${uplink_connection}"
nmcli connection up "${provider_peer_connection}"

if ! network_is_ready; then
    echo "network verification failed" >&2
    exit 1
fi

for connection in \
    "${management_connection}" \
    "${uplink_connection}" \
    "${provider_connection}" \
    "${provider_peer_connection}"
do
    nmcli connection modify \
        "${connection}" \
        connection.autoconnect yes
done

if [[ -n "${old_connection}" && \
      "${old_connection}" != "${uplink_connection}" ]]; then
    nmcli connection modify \
        "${old_connection}" \
        connection.autoconnect no
fi

echo changed
