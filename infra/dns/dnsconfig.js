// @ts-check
/// <reference path="types-dnscontrol.d.ts" />

var REG_NONE = NewRegistrar("none");
var DSP_CLOUDFLARE = NewDnsProvider("cloudflare");
var DSP_MIKROTIK = NewDnsProvider("mikrotik");
var DSP_BIND = NewDnsProvider("bind");
var jsons = [
    require("../pve/dnsconfig.json"),
    require("../vault/dnsconfig.json"),
    require("../mikrotik/dnsconfig.json"),
    require("../forgejo/dnsconfig.json"),
    require("../harbor/dnsconfig.json"),
    require("../flux/dnsconfig.json"),
];


DEFAULTS(
    CF_PROXY_DEFAULT_OFF, // turn proxy off when not specified otherwise
);

var modifiers = {
    default: [
        // github pages
        A("@", "185.199.108.153"),
        A("@", "185.199.109.153"),
        A("@", "185.199.110.153"),
        A("@", "185.199.111.153"),
        AAAA("@", "2606:50c0:8000::153"),
        AAAA("@", "2606:50c0:8001::153"),
        AAAA("@", "2606:50c0:8002::153"),
        AAAA("@", "2606:50c0:8003::153"),
        CNAME("www", "@", TTL(300)),

        // Root txt
        TXT(
            "@",
            "protonmail-verification=bdcd133d3f472fa17f66328950d02fbeae1bef75",
            TTL(300),
        ),
        TXT("@", "v=spf1 include:_spf.protonmail.ch ~all", TTL(300)),
        TXT("@", "_globalsign-domain-verification=0QBJgVV_uwcFLTi1Rot3bb1LyJ5uW1WD0ygvIS4OM5", TTL(300)),

        // Mail
        MX("@", 10, "mail.protonmail.ch.", TTL(300)),
        MX("@", 20, "mailsec.protonmail.ch.", TTL(300)),
        CNAME(
            "protonmail._domainkey",
            "protonmail.domainkey.djgwfzcu5fgjtpoijqqomgifmqj6zeiuwdd4mzim4hrxab3zsgwkq.domains.proton.ch.",
            TTL(300),
        ),
        CNAME(
            "protonmail2._domainkey",
            "protonmail2.domainkey.djgwfzcu5fgjtpoijqqomgifmqj6zeiuwdd4mzim4hrxab3zsgwkq.domains.proton.ch.",
            TTL(300),
        ),
        CNAME(
            "protonmail3._domainkey",
            "protonmail3.domainkey.djgwfzcu5fgjtpoijqqomgifmqj6zeiuwdd4mzim4hrxab3zsgwkq.domains.proton.ch.",
            TTL(300),
        ),
        TXT("_dmarc", "v=DMARC1; p=quarantine; adkim=s", TTL(300)),

        // simplelogin
        MX("simplelogin", 10, "mx1.simplelogin.co.", TTL(10800)),
        MX("simplelogin", 20, "mx2.simplelogin.co.", TTL(10800)),
        TXT(
            "simplelogin",
            "sl-verification=bxfzzfjiggzsxyzxhhmkmjqkaskjgy",
            TTL(10800),
        ),
        TXT("simplelogin", "v=spf1 include:simplelogin.co ~all", TTL(10800)),
        TXT(
            "_dmarc.simplelogin",
            "v=DMARC1; p=quarantine; pct=100; adkim=s; aspf=s",
            TTL(10800),
        ),
        CNAME(
            "dkim._domainkey.simplelogin",
            "dkim._domainkey.simplelogin.co.",
            TTL(10800),
        ),
        CNAME(
            "dkim02._domainkey.simplelogin",
            "dkim02._domainkey.simplelogin.co.",
            TTL(10800),
        ),
        CNAME(
            "dkim03._domainkey.simplelogin",
            "dkim03._domainkey.simplelogin.co.",
            TTL(10800),
        ),

        // yandex
        MX("yandex", 10, "mx.yandex.net.", TTL(21600)),
        TXT("yandex", "yandex-verification: b83672f59b3dbe16", TTL(300)),
        TXT("yandex", "v=spf1 redirect=_spf.yandex.net", TTL(300)),
        TXT(
            "mail._domainkey.yandex",
            "v=DKIM1; k=rsa; t=s; p=MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCcYzFVgkeDOhaIIkWM8gNQjxVsv0/aXfU+ax5urB5y6hA6lSjRnjRo6tm0bXbkOJf41GmiwMNgdXpwRtzgzAlX1i2aJbtEr4b9jzibEGLQ7Cvqs44bOYES9f/K3ueQpnvdTOJmFqlRReFL7ZrUyDFCoQ7f4+7h4i8s01cCcRrt5wIDAQAB",
            TTL(300),
        ),
    ]
};

for (var json_i in jsons) {
    var json = jsons[json_i];
    for (var domain_key in json.domains) {
        var domain = json.domains[domain_key];
        var cur = modifiers[domain_key];
        for (var record_key in domain.records) {
            var record = domain.records[record_key]
            if ("A" in record) {
                cur.push(A(record.A.name, record.A.address), TTL(600));
            } else if ("AAAA" in record) {
                cur.push(AAAA(record.AAAA.name, record.AAAA.address, TTL(600)));
            } else if ("CNAME" in record) {
                cur.push(CNAME(record.CNAME.name, record.CNAME.target, TTL(600)));
            }
        }
    }
}


D(
    "alwaldend.com",
    REG_NONE,
    DnsProvider(DSP_CLOUDFLARE),
    // DnsProvider(DSP_MIKROTIK),
    DnsProvider(DSP_BIND),
    modifiers["default"],
);
