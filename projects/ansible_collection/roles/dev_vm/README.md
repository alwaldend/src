---
title: Dev VM
description: Set up repositories and packages for a development VM
tags:
  - ansible_role
---

The default task entry point configures package repositories and development
packages. It uses the same `dev_vm_packages` list with either the ordinary
package manager or RPM-OSTree layering. `dev_vm_package_backend` accepts
`auto`, `package`, or `rpm_ostree`; `auto` selects RPM-OSTree when
`/run/ostree-booted` exists. Set `dev_vm_install_packages` to `false` to skip
repository and package management entirely.

RPM-OSTree changes are staged for the next boot by default. The role reports a
pending deployment but does not reboot the host. Set
`dev_vm_rpm_ostree_apply_live` to `true` only when live application is known to
be safe for the selected packages.

The packaged Fedora mirror and Bazelisk tool currently make this role specific
to Linux x86-64 development hosts.

The `packages` task entry point exposes package management separately from
other role tasks. The `bazel` task entry point installs a caller-provided Bazel
rc, Bazelisk, and the repository's `bazel_agent` runner for one developer
account.

Import the Bazel entry point with `tasks_from: bazel` and set:

- `dev_vm_bazel_user`: owner of the installed files.
- `dev_vm_bazel_group`: group of the installed files; defaults to the user.
- `dev_vm_bazel_home`: absolute home directory for the developer account.
- `dev_vm_bazelrc_src`: controller-side Bazel rc source; defaults to `bazelrc`.
- `dev_vm_bazelisk_src`: controller-side Bazelisk source; defaults to
  `bazelisk`.
- `dev_vm_bazel_agent_src`: controller-side `bazel_agent` source; defaults to
  `bazel_agent`.

The role packages its default Linux x86-64 Bazelisk and `bazel_agent` sources.
The calling playbook must package the Bazel rc, plus any overridden binary
sources, in an Ansible file search path. Cache storage and machine-specific
Bazel settings remain the caller's responsibility.
