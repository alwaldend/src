// @ts-check
/// <reference path="types-dnscontrol.d.ts" />

var REG_NONE = NewRegistrar("none");
var DSP_CLOUDFLARE = NewDnsProvider("cloudflare");
var BIND = NewDnsProvider("bind");

DEFAULTS(
    CF_PROXY_DEFAULT_OFF, // turn proxy off when not specified otherwise
);

D(
    "alwaldend.com",
    REG_NONE,
    DnsProvider(DSP_CLOUDFLARE),
    DnsProvider(BIND),

    // github pages
    A("@", "185.199.108.153"),
    A("@", "185.199.109.153"),
    A("@", "185.199.110.153"),
    A("@", "185.199.111.153"),
    AAAA("@", "2606:50c0:8000::153"),
    AAAA("@", "2606:50c0:8001::153"),
    AAAA("@", "2606:50c0:8002::153"),
    AAAA("@", "2606:50c0:8003::153"),

    // Protonmail
    TXT("@", "protonmail-verification=bdcd133d3f472fa17f66328950d02fbeae1bef75", TTL(10800)),

    // tutanota
    MX("@", 1, "mail.tutanota.de.", TTL(10800)),
    TXT("@", "t-verify=518bd0ae347aa5e3bcca726766106ad3", TTL(10800)),
    TXT("@", "v=spf1 include:spf.tutanota.de -all", TTL(10800)),
    TXT("_dmarc", "v=DMARC1; p=quarantine; adkim=s", TTL(10800)),
    CNAME("s1._domainkey", "s1.domainkey.tutanota.de.", TTL(10800)),
    CNAME("s2._domainkey", "s2.domainkey.tutanota.de.", TTL(10800)),
    CNAME("_mta-sts", "mta-sts.tutanota.de.", TTL(10800)),
    CNAME("mta-sts", "mta-sts.tutanota.de.", TTL(10800)),

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
    MX(
        "yandex", 10, "mx.yandex.net.", TTL(21600)
    ),
    TXT(
        "yandex", "yandex-verification: b83672f59b3dbe16", TTL(10800)
    ),
    TXT(
        "yandex", "v=spf1 redirect=_spf.yandex.net", TTL(10800)
    ),
    TXT(
        "mail._domainkey.yandex",
        "v=DKIM1; k=rsa; t=s; p=MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCcYzFVgkeDOhaIIkWM8gNQjxVsv0/aXfU+ax5urB5y6hA6lSjRnjRo6tm0bXbkOJf41GmiwMNgdXpwRtzgzAlX1i2aJbtEr4b9jzibEGLQ7Cvqs44bOYES9f/K3ueQpnvdTOJmFqlRReFL7ZrUyDFCoQ7f4+7h4i8s01cCcRrt5wIDAQAB",
        TTL(10800),
    ),

    // dc1
    AAAA("router1.dc1", "fd2e:546d:5738::1"),
    // cl1.dc1
    A("pi1.cl1.dc1", "192.168.1.250"),
    AAAA("pi1.cl1.dc1", "fd2e:546d:5738:0:e65d:726e:80d9:fef1"),
    CNAME("rancher.cl1.dc1", "pi1.cl1.dc1.", TTL(10800)),
);
