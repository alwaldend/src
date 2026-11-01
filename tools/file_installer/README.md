---
title: File installer
description: CLI tool to install files
languages:
  - go
  - bzl
tags:
  - bzl_rules
---

## Standalone executable installation

`//tools/file_installer/cmd/install_binary` copies a regular executable to a
private temporary sibling, sets mode 0755, and atomically replaces the installed
command. It accepts `--source`, `--name`, and optional `--directory` (default
`~/.local/bin`). It does not install runfiles or follow an automatic cache
refresh protocol. Missing or invalid sources leave an existing installation
unchanged. Tool-specific install targets supply their source and command name.

`//tools/al:install` use this
installer. Host installation requires explicit authority; validation uses an
isolated directory. Existing directory/archive installation remains separate.
