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

Python dependencies are declared in `pyproject.toml` and locked in `uv.lock`.
Update the lock with `//tools/ansible:uv_lock.update` and review its
diff. Build `//tools/ansible:uv_lock` to generate it without source writes.
These use the shared [uv locking workflow](../py/README.md), whose sandboxed resolution action
requires network access to PyPI.

The current lock uses ansible-core 2.22 for Python 3.15 targets. Its matching
Ansible community bundle is 15.0.0rc1, selected because the stable community
bundle still requires an older core. The [compatibility change](openspec/changes/support-python-315-targets/design.md)
records that release-candidate constraint and deployment validation.

Molecule is not integrated yet. The likely path is to add Molecule to
these Python requirements and invoke it through a Bazel-native Go runner,
but practical scenarios require choosing and provisioning a Molecule
driver before any scenarios can run.
