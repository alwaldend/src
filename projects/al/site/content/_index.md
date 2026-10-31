---
title: AL
linkTitle: AL
description: One command-line toolkit for repository work and independent project releases
layout: landing
statuses:
  - active
languages:
  - bzl
  - go
tags:
  - bzl_rules
  - releases
---

AL brings local commands, validation, source delivery, and project releases into
one command-line toolkit. Its reusable Bazel rules connect those operations to a
repository's own configuration, tools, and release destinations.

## From a change to a release

- Run commands with configured plugins, environment preparation, and cleanup.
- Build and test through Bazel, then prepare and validate a source candidate for
  a pull request.
- Choose a shared repository version and publish each project independently.
- Preview artifact publication offline and inspect an existing release without
  changing it.
- Publish immutable artifacts with `al release`, then extract and activate a
  website with `al release website`. Inactive releases can be deleted and published again.

The process is operator-driven and can run without CI. Source delivery, artifact
publication, and website activation are separate operations, each with its own
validation and recovery boundaries.

## Use AL

AL is developed and built from source with Bazel. A consuming repository supplies
its configuration and the integration targets required by the workflows it uses.

[Read the documentation](/docs/projects/al/) for setup, command behavior, and
release procedures. [Browse the source](https://github.com/alwaldend/src/tree/master/projects/al)
for the executable and reusable rules.
