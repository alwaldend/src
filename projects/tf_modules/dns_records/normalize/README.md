---
title: DNS normalization
description: Provider-free canonical DNS declaration transformation
---

This child module owns the canonical transformation described by the
[DNS module](../README.md). Its `document` and `zone` inputs produce
`normalized_records` without a provider, backend, or managed resources.
Terraform writers and provider-free inspection consume this output.
