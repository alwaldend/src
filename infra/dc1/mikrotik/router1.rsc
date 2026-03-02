# 2025-09-14 10:53:06 by RouterOS 7.19.4
/interface bridge
add admin-mac=78:9A:18:38:6C:CA auto-mac=no comment="bridge1 (wired)" name=bridge1
add comment="bridge2 (wireless)" name=bridge2
/interface ethernet
set [ find default-name=ether1 ] comment=ether1 l2mtu=1500 mac-address=F4:28:53:7F:A4:59
set [ find default-name=ether2 ] comment=ether2
set [ find default-name=ether3 ] comment=ether3
set [ find default-name=ether4 ] comment=ether4
set [ find default-name=ether5 ] comment=ether5
/interface wifi
set [ find default-name=wifi1 ] channel.frequency=5000-5400 .skip-dfs-channels=10min-cac comment="wifi1 (5GHz)" \
    configuration.country=Russia .mode=ap .ssid=divinity-5GHz datapath.client-isolation=yes disabled=no \
    security.authentication-types=wpa2-psk,wpa3-psk .connect-priority=0 .ft=yes .ft-over-ds=yes
set [ find default-name=wifi2 ] channel.skip-dfs-channels=10min-cac comment="wifi2 (2GHz)" configuration.country=Russia .mode=ap \
    .ssid=divinity-2GHz datapath.client-isolation=yes disabled=no security.authentication-types=wpa2-psk,wpa3-psk \
    .connect-priority=0 .ft=yes .ft-over-ds=yes
/interface list
add comment=defconf name=WAN
add comment=defconf name=LAN
add name=accept-forward-WAN
add name=accept-input-DNS
add name=accept-input-DHCP-server
add name=accept-input-ICMP
add name=accept-input-winbox
add name=accept-input-web-ui
add name=accept-input-mikrotik-neighbor-discovery
/ip pool
add comment=bridge1 name=bridge1 ranges=192.168.1.10-192.168.1.254
add comment=bridge2 name=bridge2 ranges=192.168.2.10-192.168.2.254
/ip dhcp-server
add address-pool=bridge1 comment=bridge1 interface=bridge1 lease-time=10m name=bridge1
add address-pool=bridge2 interface=bridge2 name=bridge2
/ipv6 pool
add name=dc01 prefix=fd2e:546d:5738::/48 prefix-length=64
/interface bridge port
add bridge=bridge1 comment=bridge1-ether2 interface=ether2
add bridge=bridge1 comment=bridge1-ether3 interface=ether3
add bridge=bridge1 comment=bridge1-ether4 interface=ether4
add bridge=bridge1 comment=bridge1-ether5 interface=ether5
add bridge=bridge2 comment=bridge2-wifi1 interface=wifi1
add bridge=bridge2 comment=bridge2-wifi2 interface=wifi2
/ip neighbor discovery-settings
set discover-interface-list=LAN
/interface detect-internet
set detect-interface-list=WAN
/interface list member
add interface=bridge1 list=LAN
add interface=ether1 list=WAN
add interface=bridge2 list=LAN
add interface=bridge2 list=accept-forward-WAN
add interface=bridge1 list=accept-forward-WAN
add interface=bridge1 list=accept-input-DNS
add interface=bridge2 list=accept-input-DNS
add interface=bridge1 list=accept-input-DHCP-server
add interface=bridge2 list=accept-input-DHCP-server
add interface=bridge1 list=accept-input-ICMP
add interface=bridge2 list=accept-input-ICMP
add interface=bridge1 list=accept-input-winbox
add interface=bridge1 list=accept-input-web-ui
add interface=bridge1 list=accept-input-mikrotik-neighbor-discovery
/interface ovpn-server server
add mac-address=FE:B3:B4:C4:A4:48 name=ovpn-server1
/ip address
add address=192.168.1.1/24 comment=bridge1 interface=bridge1 network=192.168.1.0
add address=192.168.2.1/24 comment=bridge2 interface=bridge2 network=192.168.2.0
/ip dhcp-client
add comment=defconf interface=ether1 use-peer-dns=no
/ip dhcp-server network
add address=192.168.1.0/24 comment=defconf dns-server=192.168.1.1 gateway=192.168.1.1
add address=192.168.2.0/24 dns-server=192.168.2.1 gateway=192.168.2.1
/ip dns
set allow-remote-requests=yes servers=1.1.1.2,1.0.0.2 use-doh-server=https://odoh.cloudflare-dns.com/dns-query verify-doh-cert=\
    yes
/ip dns static
add address=192.168.88.1 comment=defconf name=router.lan type=A
/ip firewall filter
add action=accept chain=input comment="defconf: accept established,related,untracked" connection-state=\
    established,related,untracked
add action=drop chain=input comment="defconf: drop invalid" connection-state=invalid log-prefix=drop-invalid
add action=accept chain=input comment="defconf: accept ICMP" in-interface-list=accept-input-ICMP protocol=icmp
add action=accept chain=input comment="defconf: accept to local loopback (for CAPsMAN)" dst-address=127.0.0.1
add action=drop chain=input comment="defconf: drop all not coming from LAN" in-interface-list=!LAN log-prefix=\
    drop-not-coming-from-lan
add action=accept chain=forward comment="defconf: accept in ipsec policy" ipsec-policy=in,ipsec
add action=accept chain=forward comment="defconf: accept out ipsec policy" ipsec-policy=out,ipsec
add action=fasttrack-connection chain=forward comment="defconf: fasttrack" connection-state=established,related hw-offload=yes
add action=accept chain=forward comment="defconf: accept established,related, untracked" connection-state=\
    established,related,untracked
add action=drop chain=forward comment="defconf: drop invalid" connection-state=invalid log-prefix=drop-invalid
add action=drop chain=forward comment="defconf: drop all from WAN not DSTNATed" connection-nat-state=!dstnat connection-state=\
    new in-interface-list=WAN log-prefix=drop-from-wan-not-dstnated
add action=accept chain=input in-interface-list=WAN protocol=gre
add action=accept chain=forward comment="accept forward WAN" in-interface-list=accept-forward-WAN out-interface-list=WAN
add action=accept chain=input comment="accept input DNS (udp)" dst-port=53 in-interface-list=accept-input-DNS protocol=udp
add action=accept chain=input comment="accept input DNS (tcp)" dst-port=53 in-interface-list=accept-input-DNS protocol=tcp
add action=accept chain=input comment="accept input DHCP-server" dst-port=67 in-interface-list=accept-input-DHCP-server \
    log-prefix=accept-DHCP protocol=udp
add action=accept chain=input comment="accept input winbox (tcp)" dst-port=8291 in-interface-list=accept-input-winbox protocol=\
    tcp
add action=accept chain=input comment="accept input winbox (udp)" dst-port=20561 in-interface-list=accept-input-winbox protocol=\
    udp
add action=accept chain=input comment="accept input web ui" dst-port=80,443 in-interface-list=accept-input-web-ui protocol=tcp
add action=accept chain=input comment="accept input mikrotik neighbor discovery" dst-port=5678 in-interface-list=\
    accept-input-mikrotik-neighbor-discovery protocol=udp
add action=drop chain=forward comment="drop forward" log=yes log-prefix=drop-forward
add action=drop chain=input comment="drop input" log=yes log-prefix=drop-input
/ip firewall nat
add action=masquerade chain=srcnat comment="defconf: masquerade" ipsec-policy=out,none out-interface-list=WAN
/ip ipsec profile
set [ find default=yes ] dpd-interval=2m dpd-maximum-failures=5
/ip service
set www disabled=yes
set www-ssl certificate=ssl-web-management disabled=no
/ipv6 address
add address=::1 from-pool=dc01 interface=bridge1
add address=::1 from-pool=dc01 interface=bridge2
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
add action=accept chain=input comment="defconf: accept established,related,untracked" connection-state=\
    established,related,untracked
add action=drop chain=input comment="defconf: drop invalid" connection-state=invalid
add action=accept chain=input comment="defconf: accept ICMPv6" in-interface-list=accept-input-ICMP protocol=icmpv6
add action=accept chain=input comment="defconf: accept UDP traceroute" dst-port=33434-33534 protocol=udp
add action=accept chain=input comment="defconf: accept DHCPv6-Client prefix delegation." dst-port=546 protocol=udp src-address=\
    fe80::/10
add action=accept chain=input comment="defconf: accept IKE" dst-port=500,4500 protocol=udp
add action=accept chain=input comment="defconf: accept ipsec AH" protocol=ipsec-ah
add action=accept chain=input comment="defconf: accept ipsec ESP" protocol=ipsec-esp
add action=accept chain=input comment="defconf: accept all that matches ipsec policy" ipsec-policy=in,ipsec
add action=drop chain=input comment="defconf: drop everything else not coming from LAN" in-interface-list=!LAN
add action=accept chain=forward comment="defconf: accept established,related,untracked" connection-state=\
    established,related,untracked
add action=drop chain=forward comment="defconf: drop invalid" connection-state=invalid
add action=drop chain=forward comment="defconf: drop packets with bad src ipv6" src-address-list=bad_ipv6
add action=drop chain=forward comment="defconf: drop packets with bad dst ipv6" dst-address-list=bad_ipv6
add action=drop chain=forward comment="defconf: rfc4890 drop hop-limit=1" hop-limit=equal:1 protocol=icmpv6
add action=accept chain=forward comment="defconf: accept ICMPv6" in-interface-list=accept-input-ICMP protocol=icmpv6
add action=accept chain=forward comment="defconf: accept HIP" protocol=139
add action=accept chain=forward comment="defconf: accept IKE" dst-port=500,4500 protocol=udp
add action=accept chain=forward comment="defconf: accept ipsec AH" protocol=ipsec-ah
add action=accept chain=forward comment="defconf: accept ipsec ESP" protocol=ipsec-esp
add action=accept chain=forward comment="defconf: accept all that matches ipsec policy" ipsec-policy=in,ipsec
add action=drop chain=forward comment="defconf: drop everything else not coming from LAN" in-interface-list=!LAN
add action=accept chain=forward comment="accept forward WAN" in-interface-list=accept-forward-WAN out-interface-list=WAN
add action=accept chain=input comment="accept input DNS (udp)" dst-port=53 in-interface-list=accept-input-DNS protocol=udp
add action=accept chain=input comment="accept input DNS (tcp)" dst-port=53 in-interface-list=accept-input-DNS protocol=tcp
add action=accept chain=input comment="accept input winbox (tcp)" dst-port=8291 in-interface-list=accept-input-winbox protocol=\
    tcp
add action=accept chain=input comment="accept input winbox (udp)" dst-port=20561 in-interface-list=accept-input-winbox protocol=\
    udp
add action=accept chain=input comment="accept input web ui" dst-port=80,443 in-interface-list=accept-input-web-ui protocol=tcp
add action=accept chain=input comment="accept input mikrotik neighbor discovery" dst-port=5678 in-interface-list=\
    accept-input-mikrotik-neighbor-discovery protocol=udp
add action=drop chain=forward comment="drop forward" log=yes log-prefix=drop-forward-ipv6
add action=drop chain=input comment="drop input" log=yes log-prefix=drop-input-ipv6
/ipv6 nd
set [ find default=yes ] interface=bridge1
add interface=bridge2
/system clock
set time-zone-name=Europe/Moscow
/system identity
set name=router1.dc1.alwaldend.com
/system routerboard settings
set auto-upgrade=yes
/tool mac-server
set allowed-interface-list=LAN
/tool mac-server mac-winbox
set allowed-interface-list=LAN
