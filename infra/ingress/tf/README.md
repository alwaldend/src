---
title: Tf
description: Terraform config
tags:
  - terraform
---

Use this package's `dns.plan`, `dns.show`, and `dns.apply` targets for the
[scoped DNS adoption workflow](../../dns/README.md).

DNS ownership is enabled by default after the verified 2026-09-13 adoption.
The [adoption evidence](../openspec/changes/archive/2026-09-13-adopt-dns-records/design.md)
records four imports, a no-change follow-up plan, and preserved provider inventories.
The DNS wrappers select `dns=1` and `module.dns` within this root and backend;
their validation covers DNS and its dependencies. Apply requires a reviewed
saved plan, and ordinary ingress service authentication remains separate.
