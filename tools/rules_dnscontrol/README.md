---
title: Rules Dnscontrol
description: Bazel-aware DNSControl configuration packaging
statuses:
  - active
tags:
  - dns
  - dnscontrol
  - bzl_rules
---

`rules_dnscontrol` packages a DNSControl JavaScript entrypoint together with
its record configuration files and emits a generated `requires.json` manifest
of relative paths matching the Bazel runfiles tree. Record symlinks retain a
`.json` suffix so DNSControl parses them as JSON. Load the manifest with
`require()` and load each listed path with another `require()` call;
DNSControl's JavaScript loader does not implement CommonJS `module.exports`.
