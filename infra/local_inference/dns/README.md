---
title: Local inference DNS
description: Independent Cloudflare DNS root for the local inference service
---

This Terraform project owns the public delegation and T3 Code CNAME declared
in `../dnsconfig.json`. It uses the pinned Cloudflare provider, the component
AppRole and DNS credentials, and its own Vault HTTP state backend at
`alwaldend.com/vault1/approles/src_infra_local_inference/tf_backend/dns`.
It has no dependency on GPU provisioning inputs, Yandex login, or VM state.

Use `//infra/local_inference/dns:tf.plan` to save a private plan, inspect it
with `:tf.show`, and apply the reviewed file with `:tf.apply`. Source checks
use `:tf_tests.fmt_test` and affected Bazel lint/tests. Live operations require
explicit authorization.
