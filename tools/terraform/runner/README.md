---
title: Runner
description: Golang code for the terraform runner
languages:
  - go
---

The runner initializes the owning Terraform working directory before running
the requested operation. `--direct` skips initialization.

Wrappers can set `--require-saved-plan` before the Terraform command to accept
only `apply <saved-plan-file>`. The runner rejects bare apply, other commands,
flags, extra arguments, missing files, and directories before initializing
Terraform. Relative plan paths resolve against `--chdir`; absolute paths remain
absolute. Terraform validates the plan contents and applies its recorded scope.
The guard also rejects nonempty `TF_CLI_ARGS` and `TF_CLI_ARGS_*` environment
variables before initialization so injected flags cannot change the command's
meaning. Diagnostics identify the variable name without displaying its value.
The guard is disabled by default, preserving ordinary wrappers.
