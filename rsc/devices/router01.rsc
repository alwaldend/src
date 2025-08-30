# 2025-08-30 15:51:42 by RouterOS 7.19.4
# model = C53UiG+5HPaxD2HPaxD
/interface bridge
add admin-mac=78:9A:18:38:6C:CA auto-mac=no comment=defconf name=bridge
add name=bridge-wireless
/interface ethernet
set [ find default-name=ether1 ] l2mtu=1500 mac-address=F4:28:53:7F:A4:59
/interface wifi
set [ find default-name=wifi2 ] channel.skip-dfs-channels=10min-cac \
    configuration.country=Russia .mode=ap .ssid=divinity-2GHz disabled=no \
    name=divinity-2GHz security.authentication-types=wpa2-psk,wpa3-psk \
    .connect-priority=0 .ft=yes .ft-over-ds=yes
set [ find default-name=wifi1 ] channel.frequency=5000-5400 \
    .skip-dfs-channels=10min-cac configuration.country=Russia .mode=ap .ssid=\
    divinity-5GHz disabled=no name=divinity-5GHz \
    security.authentication-types=wpa2-psk,wpa3-psk .connect-priority=0 .ft=\
    yes .ft-over-ds=yes
/interface list
add comment=defconf name=WAN
add comment=defconf name=LAN
add comment=wireguard name=WG
/ip pool
add name=dhcp-wired ranges=192.168.1.10-192.168.1.254
add name=dhcp-wireless ranges=192.168.2.10-192.168.2.254
/ip dhcp-server
add address-pool=dhcp-wired interface=bridge lease-time=10m name=dhcp-wired
add address-pool=dhcp-wireless interface=bridge-wireless name=dhcp-wireless
/interface bridge port
add bridge=bridge comment=defconf interface=ether2
add bridge=bridge comment=defconf interface=ether3
add bridge=bridge comment=defconf interface=ether4
add bridge=bridge comment=defconf interface=ether5
add bridge=bridge-wireless comment=defconf interface=divinity-5GHz
add bridge=bridge-wireless comment=defconf interface=divinity-2GHz
add bridge=bridge comment=wireless interface=LAN
/ip neighbor discovery-settings
set discover-interface-list=LAN
/interface detect-internet
set detect-interface-list=WAN
/interface list member
add comment=defconf interface=bridge list=LAN
add comment=defconf interface=ether1 list=WAN
add interface=bridge-wireless list=LAN
/interface ovpn-server server
add mac-address=FE:B3:B4:C4:A4:48 name=ovpn-server1
/ip address
add address=192.168.1.1/24 comment=wired interface=bridge network=192.168.1.0
add address=192.168.2.1/24 comment=wireless-dhcp interface=bridge-wireless \
    network=192.168.2.0
/ip dhcp-client
add comment=defconf interface=ether1 use-peer-dns=no
/ip dhcp-server network
add address=192.168.1.0/24 comment=defconf dns-server=192.168.1.1 gateway=\
    192.168.1.1
add address=192.168.2.0/24 dns-server=192.168.2.1 gateway=192.168.2.1
/ip dns
set allow-remote-requests=yes use-doh-server=https://1.1.1.1/dns-query \
    verify-doh-cert=yes
/ip dns static
add address=192.168.88.1 comment=defconf name=router.lan type=A
/ip firewall filter
add action=accept chain=input comment=\
    "defconf: accept established,related,untracked" connection-state=\
    established,related,untracked
add action=drop chain=input comment="defconf: drop invalid" connection-state=\
    invalid
add action=accept chain=input comment="defconf: accept ICMP" protocol=icmp
add action=accept chain=input comment=\
    "defconf: accept to local loopback (for CAPsMAN)" dst-address=127.0.0.1
add action=drop chain=input comment="defconf: drop all not coming from LAN" \
    in-interface-list=!LAN
add action=accept chain=forward comment="defconf: accept in ipsec policy" \
    ipsec-policy=in,ipsec
add action=accept chain=forward comment="defconf: accept out ipsec policy" \
    ipsec-policy=out,ipsec
add action=fasttrack-connection chain=forward comment="defconf: fasttrack" \
    connection-state=established,related hw-offload=yes
add action=accept chain=forward comment=\
    "defconf: accept established,related, untracked" connection-state=\
    established,related,untracked
add action=drop chain=forward comment="defconf: drop invalid" \
    connection-state=invalid
add action=drop chain=forward comment=\
    "defconf: drop all from WAN not DSTNATed" connection-nat-state=!dstnat \
    connection-state=new in-interface-list=WAN
add action=accept chain=input in-interface=bridge protocol=gre
/ip firewall nat
add action=masquerade chain=srcnat comment="defconf: masquerade" ipsec-policy=\
    out,none out-interface-list=WAN
/ip ipsec profile
set [ find default=yes ] dpd-interval=2m dpd-maximum-failures=5
/ip service
set www-ssl certificate=ssl-web-management disabled=no
/ipv6 address
add address=1000::5 advertise=no comment=wireguard interface=*A
/ipv6 firewall address-list
add address=::/128 comment="defconf: unspecified address" list=bad_ipv6
add address=::1/128 comment="defconf: lo" list=bad_ipv6
add address=fec0::/10 comment="defconf: site-local" list=bad_ipv6
add address=::ffff:0.0.0.0/96 comment="defconf: ipv4-mapped" list=bad_ipv6
add address=::/96 comment="defconf: ipv4 compat" list=bad_ipv6
add address=100::/64 comment="defconf: discard only " list=bad_ipv6
add address=2001:db8::/32 comment="defconf: documentation" list=bad_ipv6
add address=2001:10::/28 comment="defconf: ORCHID" list=bad_ipv6
add address=3ffe::/16 comment="defconf: 6bone" list=bad_ipv6
/ipv6 firewall filter
add action=accept chain=input comment=\
    "defconf: accept established,related,untracked" connection-state=\
    established,related,untracked
add action=drop chain=input comment="defconf: drop invalid" connection-state=\
    invalid
add action=accept chain=input comment="defconf: accept ICMPv6" protocol=icmpv6
add action=accept chain=input comment="defconf: accept UDP traceroute" \
    dst-port=33434-33534 protocol=udp
add action=accept chain=input comment=\
    "defconf: accept DHCPv6-Client prefix delegation." dst-port=546 protocol=\
    udp src-address=fe80::/10
add action=accept chain=input comment="defconf: accept IKE" dst-port=500,4500 \
    protocol=udp
add action=accept chain=input comment="defconf: accept ipsec AH" protocol=\
    ipsec-ah
add action=accept chain=input comment="defconf: accept ipsec ESP" protocol=\
    ipsec-esp
add action=accept chain=input comment=\
    "defconf: accept all that matches ipsec policy" ipsec-policy=in,ipsec
add action=drop chain=input comment=\
    "defconf: drop everything else not coming from LAN" in-interface-list=!LAN
add action=accept chain=forward comment=\
    "defconf: accept established,related,untracked" connection-state=\
    established,related,untracked
add action=drop chain=forward comment="defconf: drop invalid" \
    connection-state=invalid
add action=drop chain=forward comment=\
    "defconf: drop packets with bad src ipv6" src-address-list=bad_ipv6
add action=drop chain=forward comment=\
    "defconf: drop packets with bad dst ipv6" dst-address-list=bad_ipv6
add action=drop chain=forward comment="defconf: rfc4890 drop hop-limit=1" \
    hop-limit=equal:1 protocol=icmpv6
add action=accept chain=forward comment="defconf: accept ICMPv6" protocol=\
    icmpv6
add action=accept chain=forward comment="defconf: accept HIP" protocol=139
add action=accept chain=forward comment="defconf: accept IKE" dst-port=\
    500,4500 protocol=udp
add action=accept chain=forward comment="defconf: accept ipsec AH" protocol=\
    ipsec-ah
add action=accept chain=forward comment="defconf: accept ipsec ESP" protocol=\
    ipsec-esp
add action=accept chain=forward comment=\
    "defconf: accept all that matches ipsec policy" ipsec-policy=in,ipsec
add action=drop chain=forward comment=\
    "defconf: drop everything else not coming from LAN" in-interface-list=!LAN
/system clock
set time-zone-name=Europe/Moscow
/system identity
set name=router01.dc01.alwaldend.com
/tool mac-server
set allowed-interface-list=LAN
/tool mac-server mac-winbox
set allowed-interface-list=LAN
