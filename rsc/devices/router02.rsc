# 2025-09-04 20:00:51 by RouterOS 7.19.4
# model = L009UiGS-2HaxD
/interface bridge
add admin-mac=78:9A:18:FB:9D:1B auto-mac=no comment=defconf name=bridgeLocal
/interface wifi
set [ find default-name=wifi1 ] configuration.mode=ap
/interface wifi datapath
add bridge=bridgeLocal comment=defconf disabled=no name=capdp
/port
set 0 name=serial0
/interface bridge port
add bridge=bridgeLocal comment=defconf interface=ether2
add bridge=bridgeLocal comment=defconf interface=ether3
add bridge=bridgeLocal comment=defconf interface=ether4
add bridge=bridgeLocal comment=defconf interface=ether5
add bridge=bridgeLocal comment=defconf interface=ether6
add bridge=bridgeLocal comment=defconf interface=ether7
add bridge=bridgeLocal comment=defconf interface=ether8
add bridge=bridgeLocal comment=defconf interface=sfp1
/ipv6 settings
set accept-router-advertisements=yes
/interface wifi cap
set discovery-interfaces=""
/ip dhcp-client
add comment=defconf interface=ether1
/ipv6 dhcp-client
add add-default-route=yes comment=dc01 interface=ether1 pool-name=dc01 \
    request=prefix
/system clock
set time-zone-name=Europe/Moscow
/system identity
set name=router02.dc01.alwaldend.com
/system routerboard settings
set enter-setup-on=delete-key
