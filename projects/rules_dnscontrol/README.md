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
its record configuration files and emits a generated `requires.js` manifest
whose paths match the Bazel runfiles tree.
