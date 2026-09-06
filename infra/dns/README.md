---
title: Dns
description: Dns setup for alwaldend.com
tags:
  - dns
  - dnscontrol
---

## Links

- Docs: https://docs.dnscontrol.org
- Rules: [../../projects/rules_dnscontrol](../../projects/rules_dnscontrol)

## Deployment

`//infra/dns:config_test` validates the generated configuration with the pinned
DNSControl `check` command. It requires no credentials and does not access DNS
providers.

Operational wrappers pass `providers.json` explicitly. Its environment references
are populated by the AL Vault injector; DNSControl reads the values directly
without a rendered credentials file.

Project landing records live next to their projects in
`projects/<project>/dnsconfig.json`. The main website owns the shared
`pages` address, and project files own their own subdomain CNAMEs.

Interactive:

```sh
bazel run //infra/dns
```

Just preview:

```sh
bazel run //infra/dns:dns.preview
```

Just deploy:

```sh
bazel run //infra/dns:dns.deploy
```

{{% alwaldend/alert %}}
Deploy modifies the bind file, which will cause
//infra/alwaldend.com/dns:preview_test to fail
{{% /alwaldend/alert %}}

## dc1 BIND

{{/*< readfile file="zones/dc1/alwaldend.com.zone" code="true" lang="zone" >*/}}

## global BIND

{{/*< readfile file="zones/global/alwaldend.com.zone" code="true" lang="zone" >*/}}
