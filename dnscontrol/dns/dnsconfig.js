// @ts-check
/// <reference path="types-dnscontrol.d.ts" />

var REG_NONE = NewRegistrar("none");
var DSP_CLOUDFLARE = NewDnsProvider("cloudflare");

DEFAULTS(
  CF_PROXY_DEFAULT_OFF, // turn proxy off when not specified otherwise
);

D(
  "alwaldend.com",
  REG_NONE,
  DnsProvider(DSP_CLOUDFLARE),

  // tutanota
  MX("@", 0, "mail.tutanota.de.", TTL(10800)),
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

  // dc1
  A("pi1.dc1", "192.168.1.250"),
  AAAA("pi1.dc1", "fd2e:546d:5738:0:e65d:726e:80d9:fef1"),
  AAAA("router1.dc1", "fd2e:546d:5738::1"),
);
