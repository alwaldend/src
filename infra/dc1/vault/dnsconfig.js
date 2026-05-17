// @ts-check
/// <reference path="../../dns/types-dnscontrol.d.ts" />

function IncludeDc1Vault() {
    return [
        A("bm3.dc1", "192.168.1.218"),
        CNAME("vault.dc1", "bm3.dc1.alwaldend.com."),
    ]
}
