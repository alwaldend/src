---
title: Dns
description: Dns config
tags:
  - dns
  - dnscontrol
---

## Deployment

- Setup environment variables from [creds.json.tpl](./creds.json.tpl)
- Modify [dnsconfig.js](./dnsconfig.js)
- Preview changes:
  ```sh
  bazel run //infra/dns:preview
  ```
- Push changes:
  ```sh
  bazel run //infra/dns:push
  ```

{{% alwaldend/alert %}}
If BIND files differ after push, test `//infra/dns:preview_test` will fail
{{% /alwaldend/alert %}}

## Links

- Website: https://dnscontrol.org
- Docs: https://docs.dnscontrol.org
- Rules: [../../tools/dnscontrol](../../tools/dnscontrol)

## Zones

### Alwaldend.com

{{< readfile file="zones/alwaldend.com.zone" code="true" lang="zone" >}}
