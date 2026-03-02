---
title: Dns
description: Dns setup for alwaldend.com
tags:
  - dns
  - dnscontrol
---

## Links

- Docs: https://docs.dnscontrol.org
- Rules: [../../tools/dnscontrol](../../tools/dnscontrol)

## Deployment

- Modify `dnsconfig.js`
- Preview changes:
  ```sh
  bazel run //infra/dns:preview
  ```
- Deploy changes:
  ```sh
  bazel run //infra/dns:deploy
  ```

{{% alwaldend/alert %}}
Deploy modifies the bind file, which will cause
//infra/alwaldend.com/dns:preview_test to fail
{{% /alwaldend/alert %}}

## BIND file

{{< readfile file="zones/alwaldend.com.zone" code="true" lang="zone" >}}
