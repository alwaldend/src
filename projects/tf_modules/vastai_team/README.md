---
title: Vast.ai team module
description: A zero-credit team with separately generated project API keys
---

This module creates one Vast.ai team and a map of project API keys through the
in-tree provider. It declares no compute resources. The caller supplies names
and JSON permission objects and selects provider credentials. Generated key
outputs are sensitive and belong directly in secret resources, with secure state.

Instantiate the module with `for_each` to manage multiple teams. The
[provider integration PR](https://github.com/alwaldend/src/pull/249) demonstrates a checked-in
catalog and direct delivery to Vault. Vault identity and destination paths belong
to the caller; this module grants no Vault permissions.

The provider's Terraform CLI fixture verifies two independent team instances,
three project keys, unchanged plans, and complete cleanup without a live account.
