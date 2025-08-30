# 2025-08-30 15:42:48 by RouterOS 7.19.4
# model = L009UiGS-2HaxD
/interface bridge
add admin-mac=78:9A:18:FB:9D:1B auto-mac=no comment=defconf name=bridgeLocal
/interface wifi datapath
add bridge=bridgeLocal comment=defconf disabled=no name=capdp
/interface wifi
# no connection to CAPsMAN
set [ find default-name=wifi1 ] configuration.manager=capsman .mode=ap \
    datapath=capdp
/port
set 0 name=serial0
/interface bridge port
add bridge=bridgeLocal comment=defconf interface=ether1
add bridge=bridgeLocal comment=defconf interface=ether2
add bridge=bridgeLocal comment=defconf interface=ether3
add bridge=bridgeLocal comment=defconf interface=ether4
add bridge=bridgeLocal comment=defconf interface=ether5
add bridge=bridgeLocal comment=defconf interface=ether6
add bridge=bridgeLocal comment=defconf interface=ether7
add bridge=bridgeLocal comment=defconf interface=ether8
add bridge=bridgeLocal comment=defconf interface=sfp1
/interface wifi cap
set discovery-interfaces=bridgeLocal enabled=yes slaves-datapath=capdp
/ip dhcp-client
# DHCP client can not run on slave or passthrough interface!
add comment=defconf interface=ether1
/system clock
set time-zone-name=Europe/Moscow
/system identity
set name=router02.dc01.alwaldend.com
/system routerboard settings
set enter-setup-on=delete-key
