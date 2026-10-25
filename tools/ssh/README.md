---
title: SSH
description: Host SSH client wrapped by AL
---

`ssh_binary` runs the host's `ssh` from `PATH` through `al_binary_run`.
AL configurations and plugins can inject SSH keys before the command runs.

```starlark
load("//tools/ssh:defs.bzl", "ssh_binary")

ssh_binary(
    name = "ssh",
    configs = [":config"],
    args = ["example.invalid"],
)
```

The macro accepts `al_binary_run` attributes, including `configs`, `run_args`,
`data`, `env`, and `visibility`. Declare injection plugins in `data` and select
their calls through `run_args`. The consuming AL configuration owns Vault
references and key injection. No authentication is selected by default.

Run `bazel run //tools/ssh -- -V` to check the host client. Arguments supplied
at invocation are appended after the target's fixed SSH arguments. The host
must provide OpenSSH; it is not downloaded or built by Bazel.
