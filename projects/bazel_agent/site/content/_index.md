---
title: Bazel agent
linkTitle: Bazel agent
description: Bazel runner for repository agents
layout: landing
statuses:
  - active
languages:
  - go
  - sh
tags:
  - bazel
  - agent
---

`bazel_agent` is the Bazel entry point for agents working in this
repository. It keeps Bazel behind a validated subcommand while consistently
applying the agent configuration.

## Features

- Validated `bazel` subcommand instead of arbitrary arguments
- Applies the agent configuration in the required position
- Preserves signals and exit status by replacing itself with the Bazel process
- Read-only `doctor` diagnostics for runner and workspace state
