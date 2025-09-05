# 2025-09-05 19:25:37 by RouterOS 7.19.4
# model = L009UiGS-2HaxD
/interface bridge
add name=bridge01
/port
set 0 name=serial0
/interface bridge port
add bridge=bridge01 interface=ether2
add bridge=bridge01 interface=ether3
add bridge=bridge01 interface=ether4
add bridge=bridge01 interface=ether5
add bridge=bridge01 interface=ether6
add bridge=bridge01 interface=ether7
add bridge=bridge01 interface=ether8
add bridge=bridge01 interface=ether1
/ip neighbor discovery-settings
set discover-interface-list=!dynamic
/ipv6 settings
set accept-router-advertisements=yes
/ip dhcp-client
add interface=bridge01
/ipv6 dhcp-client
add add-default-route=yes interface=bridge01 pool-name=dc01 request=prefix
/ipv6 nd
set [ find default=yes ] managed-address-configuration=yes \
    other-configuration=yes
/system clock
set time-zone-name=Europe/Moscow
/system identity
set name=router02.dc01.alwaldend.com
/system routerboard settings
set enter-setup-on=delete-key
