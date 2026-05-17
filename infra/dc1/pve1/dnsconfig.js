// @ts-check
/// <reference path="../../dns/types-dnscontrol.d.ts" />

function IncludeDc1Pve1() {
    return [
        A("bm2.dc1", "192.168.1.216"),
        CNAME("host1.pve1.dc1", "bm2.dc1.alwaldend.com."),
        A("cloudinit-test.vm.pve1.dc1", "192.168.10.10"),
    ];
}
