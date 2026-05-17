// @ts-check
/// <reference path="../../dns/types-dnscontrol.d.ts" />

function IncludeDc1Consul1() {
    return [
        A("host1.consul1.dc1", "192.168.10.20"),
        CNAME("consul1.dc1", "host1.consul1.dc1.alwaldend.com."),
    ];
}
