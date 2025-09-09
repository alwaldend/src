# 2025-09-09 16:28:25 by RouterOS 7.19.4
# model = C53UiG+5HPaxD2HPaxD
/interface bridge
add admin-mac=78:9A:18:38:6C:CA auto-mac=no comment=bridge01 name=bridge01
/interface ethernet
set [ find default-name=ether1 ] comment=ether01 l2mtu=1500 mac-address=F4:28:53:7F:A4:59 name=ether01
set [ find default-name=ether2 ] comment=ether02 name=ether02
set [ find default-name=ether3 ] comment=ether03 name=ether03
set [ find default-name=ether4 ] comment=ether04 name=ether04
set [ find default-name=ether5 ] comment=ether05 name=ether05
/interface wifi
set [ find default-name=wifi1 ] channel.frequency=5000-5400 .skip-dfs-channels=10min-cac comment="wifi01 (5GHz)" configuration.country=Russia .mode=ap .ssid=divinity-5GHz \
    disabled=no name=wifi01 security.authentication-types=wpa2-psk,wpa3-psk .connect-priority=0 .ft=yes .ft-over-ds=yes
set [ find default-name=wifi2 ] channel.skip-dfs-channels=10min-cac comment="wifi02 (2GHz)" configuration.country=Russia .mode=ap .ssid=divinity-2GHz disabled=no name=wifi02 \
    security.authentication-types=wpa2-psk,wpa3-psk .connect-priority=0 .ft=yes .ft-over-ds=yes
/interface list
add comment=defconf name=WAN
add comment=defconf name=LAN
add comment=wireguard name=WG
add comment=wireless name=wireless
/ip pool
add comment=bridge01 name=bridge01 ranges=192.168.1.10-192.168.1.254
/ip dhcp-server
add address-pool=bridge01 comment=bridge01 interface=bridge01 lease-time=10m name=bridge01
/ipv6 pool
add name=dc01 prefix=fd2e:546d:5738::/48 prefix-length=64
/interface bridge port
add bridge=bridge01 comment=bridge01-ether02 interface=ether02
add bridge=bridge01 comment=bridge01-ether03 interface=ether03
add bridge=bridge01 comment=bridge01-ether04 interface=ether04
add bridge=bridge01 comment=bridge01-ether05 interface=ether05
add bridge=bridge01 comment=bridge01-wifi01 interface=wifi01
add bridge=bridge01 comment=bridge01-wifi02 interface=wifi02
add bridge=bridge01 comment=bridge01-LAN interface=LAN
/ip neighbor discovery-settings
set discover-interface-list=LAN
/interface detect-internet
set detect-interface-list=WAN
/interface list member
add comment=LAN-bridge01 interface=bridge01 list=LAN
add comment=WAN-ether01 interface=ether01 list=WAN
add comment=wireless-wifi02 interface=wifi02 list=wireless
add comment=wireless-wifi01 interface=wifi01 list=wireless
/interface ovpn-server server
add mac-address=FE:B3:B4:C4:A4:48 name=ovpn-server1
/ip address
add address=192.168.1.1/24 comment=bridge01 interface=bridge01 network=192.168.1.0
/ip dhcp-client
add comment=defconf interface=ether01 use-peer-dns=no
/ip dhcp-server network
add address=192.168.1.0/24 comment=defconf dns-server=192.168.1.1 gateway=192.168.1.1
add address=192.168.2.0/24 dns-server=192.168.2.1 gateway=192.168.2.1
/ip dns
set allow-remote-requests=yes servers=1.1.1.2,1.0.0.2 use-doh-server=https://odoh.cloudflare-dns.com/dns-query verify-doh-cert=yes
/ip dns static
add address=192.168.88.1 comment=defconf name=router.lan type=A
/ip firewall filter
add action=accept chain=input comment="defconf: accept established,related,untracked" connection-state=established,related,untracked
add action=drop chain=input comment="defconf: drop invalid" connection-state=invalid log-prefix=drop-invalid
add action=accept chain=input comment="defconf: accept ICMP" protocol=icmp
add action=accept chain=input comment="defconf: accept to local loopback (for CAPsMAN)" dst-address=127.0.0.1
add action=drop chain=input comment="defconf: drop all not coming from LAN" in-interface-list=!LAN log-prefix=drop-not-coming-from-lan
add action=accept chain=forward comment="defconf: accept in ipsec policy" ipsec-policy=in,ipsec
add action=accept chain=forward comment="defconf: accept out ipsec policy" ipsec-policy=out,ipsec
add action=fasttrack-connection chain=forward comment="defconf: fasttrack" connection-state=established,related hw-offload=yes
add action=accept chain=forward comment="defconf: accept established,related, untracked" connection-state=established,related,untracked
add action=drop chain=forward comment="defconf: drop invalid" connection-state=invalid log-prefix=drop-invalid
add action=drop chain=forward comment="defconf: drop all from WAN not DSTNATed" connection-nat-state=!dstnat connection-state=new in-interface-list=WAN log-prefix=\
    drop-from-wan-not-dstnated
add action=accept chain=input in-interface=bridge01 protocol=gre
add action=drop chain=forward comment="drop wireless not to WAN" in-interface-list=wireless log=yes log-prefix=drop-wireless-not-to-wan out-interface-list=!WAN
add action=drop chain=input comment="drop wireless web access" dst-port=80,443 in-interface-list=wireless log=yes log-prefix=drop-wireless-web-access protocol=tcp
/ip firewall nat
add action=masquerade chain=srcnat comment="defconf: masquerade" ipsec-policy=out,none out-interface-list=WAN
/ip ipsec profile
set [ find default=yes ] dpd-interval=2m dpd-maximum-failures=5
/ip service
set www disabled=yes
set www-ssl certificate=ssl-web-management disabled=no
/ipv6 address
add address=::1 from-pool=dc01 interface=bridge01
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
add action=accept chain=input comment="defconf: accept established,related,untracked" connection-state=established,related,untracked
add action=drop chain=input comment="defconf: drop invalid" connection-state=invalid
add action=accept chain=input comment="defconf: accept ICMPv6" protocol=icmpv6
add action=accept chain=input comment="defconf: accept UDP traceroute" dst-port=33434-33534 protocol=udp
add action=accept chain=input comment="defconf: accept DHCPv6-Client prefix delegation." dst-port=546 protocol=udp src-address=fe80::/10
add action=accept chain=input comment="defconf: accept IKE" dst-port=500,4500 protocol=udp
add action=accept chain=input comment="defconf: accept ipsec AH" protocol=ipsec-ah
add action=accept chain=input comment="defconf: accept ipsec ESP" protocol=ipsec-esp
add action=accept chain=input comment="defconf: accept all that matches ipsec policy" ipsec-policy=in,ipsec
add action=drop chain=input comment="defconf: drop everything else not coming from LAN" in-interface-list=!LAN
add action=accept chain=forward comment="defconf: accept established,related,untracked" connection-state=established,related,untracked
add action=drop chain=forward comment="defconf: drop invalid" connection-state=invalid
add action=drop chain=forward comment="defconf: drop packets with bad src ipv6" src-address-list=bad_ipv6
add action=drop chain=forward comment="defconf: drop packets with bad dst ipv6" dst-address-list=bad_ipv6
add action=drop chain=forward comment="defconf: rfc4890 drop hop-limit=1" hop-limit=equal:1 protocol=icmpv6
add action=accept chain=forward comment="defconf: accept ICMPv6" protocol=icmpv6
add action=accept chain=forward comment="defconf: accept HIP" protocol=139
add action=accept chain=forward comment="defconf: accept IKE" dst-port=500,4500 protocol=udp
add action=accept chain=forward comment="defconf: accept ipsec AH" protocol=ipsec-ah
add action=accept chain=forward comment="defconf: accept ipsec ESP" protocol=ipsec-esp
add action=accept chain=forward comment="defconf: accept all that matches ipsec policy" ipsec-policy=in,ipsec
add action=drop chain=forward comment="defconf: drop everything else not coming from LAN" in-interface-list=!LAN
add action=drop chain=forward comment="drop wireless not to WAN" in-interface-list=wireless log=yes log-prefix=wireless-not-to-wan-ipv6 out-interface-list=!WAN
/ipv6 nd
set [ find default=yes ] interface=bridge01
/system clock
set time-zone-name=Europe/Moscow
/system identity
set name=router01.dc01.alwaldend.com
/tool mac-server
set allowed-interface-list=LAN
/tool mac-server mac-winbox
set allowed-interface-list=LAN
