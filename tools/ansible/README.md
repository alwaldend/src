---
title: Ansible
description: Bazel rules for ansible
languages:
  - bzl
tags:
  - bzl_rules
---

`ansible_lint` runs `ansible-lint` from Bazel-managed Python
dependencies. A Go test runner exposes the CLI binaries under their
hyphenated names, configures writable Ansible paths, and invokes the
linter against workspace-relative sources.

The same Python lock supplies Molecule and the seed-image library to the
[local QEMU runner](../molecule/README.md). It uses Ansible-native lifecycle
playbooks and the existing packaged collection mappings.
