---
title: Al
description: Repository command runner and Bazel configuration rules
statuses:
  - active
languages:
  - bzl
  - go
tags:
  - bzl_rules
  - fp
  - proto
---

AL is a command runner that prepares credentials and environment variables
through plugins, runs your command, and cleans up afterward. Its Bazel rules
package commands with their configuration and required plugins.

`al tool` runs tools declared in AL configuration from a source-keyed executable
cache. With no `--config`, it uses the repository-root `al.lua`; explicit
`--config` paths replace that default. Cache misses build the tool's configured
Bazel label and publish the selected executable under
`$XDG_CACHE_HOME/al/tools`. `AL_TOOL_CACHE` or `tool --cache-root PATH` selects
an alternate cache root.

Repository agent skills are written by `bazel run //.agents:write_skills`
and checked by `//.agents:write_skills_test`.

## Links

- [Source code](https://github.com/alwaldend/src/tree/master/projects/al)
- [Bazel rules](https://alwaldend.com/docs/projects/al/rules/)
- [Command and plugin lifecycle](https://alwaldend.com/docs/projects/al/docs/)
