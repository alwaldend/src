// @ts-check
/// <reference path="../dns/types-dnscontrol.d.ts" />

require("./vault/dnsconfig.js");
require("./pve1/dnsconfig.js");
require("./consul1/dnsconfig.js");

function IncludeDc1() {
    return [
        A("router1.dc1", "192.168.1.1"),
        AAAA("router1.dc1", "fd2e:546d:5738::1"),
        A("bm1.dc1", "192.168.1.222"),
        IncludeDc1Vault(),
        IncludeDc1Pve1(),
        IncludeDc1Consul1(),
    ];
}
