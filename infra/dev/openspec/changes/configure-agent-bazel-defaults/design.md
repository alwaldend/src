## Context

The development VM role already copies `ansible/files/bazelrc` to the managed
user's `~/.bazelrc`. The AL change removes wrappers that supplied
`--config=agent`, leaving Bazel to read its ordinary configuration files.

## Goals / Non-Goals

Keep agent output preferences in the existing host configuration and make
ordinary Bazel commands usable from other repositories. Preserve all cache,
SDK, and scheduling settings. This change does not deploy host configuration.

## Decisions

Move the existing nine output options into a labelled block of the managed
user bazelrc, using ordinary `common`, `build`, and `test` entries. A global
`common --config=agent` would depend on each workspace defining that profile;
direct options need no repository-specific profile. The setup page links to
this block instead of maintaining a second copy of the options.

Keep the existing Ansible copy and installation tasks. The host continues to
install Bazelisk and AL; installed AL supplies release workspace status while
Bazel owns build execution.

## Risks / Trade-offs

The new defaults take effect on the live host only after an authorized Ansible
deployment. Operators can override them per invocation using Bazel options.
Other installations must adapt cache paths and SDK locations rather than copy
the entire Host Bot configuration.

## Migration Plan

Validate the packaged Ansible input, documentation, and an isolated invocation
that loads the checked-in output settings. Record source validation without
claiming that the live host has changed.
