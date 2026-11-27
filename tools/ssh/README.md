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

`//tools/ssh:forward` provides the AL `ssh_forward` resource. Add the
provider to the command's AL configuration. Declare a resource with the SSH
server `host` and port, SSH `user`, remote TCP `target`, and trusted
`host_ca` public key. Give it one dependency that returns the
`private_key` and `public_key_signed` paths, such as `vault_ssh`.
Use the resource's `.Last.Data.address` in a dependent environment
variable. The provider verifies the SSH host and remote service before
returning an ephemeral loopback address. The forward closes when the
command exits.
