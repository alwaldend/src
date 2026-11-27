---
title: File installer
description: CLI tool to install files
languages:
  - go
  - bzl
tags:
  - bzl_rules
---

## Directory and archive installation

`//tools/file_installer` accepts `install --target DIRECTORY SOURCE...`. A
directory source creates symlinks to its files. Reinstalling an identical
link succeeds, while a conflicting destination reports an error. Traversal
skips a destination nested inside the source and accepts up to 128 directory
levels. A tar source requires an existing destination directory. Each regular
archive entry is staged beside its destination and replaced after the write
completes. Replacing a regular file retains its permissions. New files use the
archive mode subject to the process umask. An interrupted installation can be
retried.

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
