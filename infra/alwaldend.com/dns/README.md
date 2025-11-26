---
title: Dns
description: Dns setup for alwaldend.com
tags:
  - dns
  - dnscontrol
---

## Deployment

- Setup environment variables from [creds.json.tpl](./creds.json.tpl)
- Modify [dnsconfig.js](./dnsconfig.js)
- Preview changes:
  ```sh
  bazel run //infra/alwaldend.com/dns:preview
  ```
- Deploy changes:
  ```sh
  bazel run //infra/alwaldend.com/dns:deploy
  ```

{{% alwaldend/alert %}}
If BIND files differ after push, test `//infra/alwaldend.com/dns:preview_test` will fail
{{% /alwaldend/alert %}}

## Links

- Website: https://dnscontrol.org
- Docs: https://docs.dnscontrol.org
- Rules: [../../tools/dnscontrol](../../../tools/dnscontrol)

## Zones

### Alwaldend.com

{{< readfile file="zones/alwaldend.com.zone" code="true" lang="zone" >}}
