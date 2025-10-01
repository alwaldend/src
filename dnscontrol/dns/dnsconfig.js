var REG_NONE = NewRegistrar("none");
var DSP_CLOUDFLARE = NewDnsProvider("cloudflare");
D("alwaldend.com", REG_NONE, DnsProvider(DSP_CLOUDFLARE), A("@", "1.2.3.4"));
