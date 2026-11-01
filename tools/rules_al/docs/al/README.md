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

Repository tools are built through their ordinary Bazel targets. Use explicit
install targets, such as `//tools/al:install`, for standalone host
commands. AL does not cache, build or install repository tool executables.
The former `al tool` command and tool declarations have been removed.

Repository agent skills are written by `bazel run //.agents:write_skills`
and checked by `//.agents:write_skills_test`.

## Links

- [Source code](https://github.com/alwaldend/src/tree/master/tools/rules_al)
- [Bazel rules](https://alwaldend.com/docs/tools/rules_al/rules/)
- [Command and plugin lifecycle](https://alwaldend.com/docs/tools/rules_al/docs/al/lifecycle/)
