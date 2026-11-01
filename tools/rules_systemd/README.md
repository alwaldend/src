---
title: Rules systemd
description: Executable Bazel targets with user systemd resource limits
languages:
  - bzl
  - go
---

`rules_systemd` is a standalone Linux Bazel module. Add a `bazel_dep` on
`rules_systemd` and load its public `systemd/defs.bzl` API. The module uses
its own dependency manifest and lockfile; the parent checkout supplies a local
module override for development.

The `systemd_run` rule exposes a normal executable Bazel target with its
resource policy declared in BUILD.bazel. It runs the application through the
host's `systemd-run --user --scope`, preserving literal arguments, environment,
working directory, standard streams, signals, and exit status. It does not
install services or change persistent host configuration.

```bzl
load("@rules_systemd//systemd:defs.bzl", "systemd_run")

systemd_run(
    name = "server",
    binary = ":server_bin",
    arguments = ["--input", "$(rlocationpath :input.txt)"],
    data = ["input.txt"],
    properties = {
        "MemoryHigh": "96M",
        "MemoryMax": "128M",
        "MemorySwapMax": "0",
    },
)
```

Run it with `bazel_agent bazel run //path/to:server`. Arguments after `--`
are appended to `arguments`. Configure the application's default arguments
on the wrapper; defaults from the raw binary's Bazel `args` attribute are
not inherited. `arguments` supports location/runfiles expressions referring
to `binary` and `data`. The rule merges their runfiles and preserves the
binary's `RunEnvironmentInfo`.

The rule uses the existing hermetic copy action to package the Go launcher
and an adjacent generated JSON configuration. The launcher resolves the
application through Bazel runfiles. Both output files are required when
packaging the executable. No shell launch script or per-target Go build is
needed. All systemd execution occurs at runtime, outside build actions.

The host systemd must support `--expand-environment=no`; the launcher uses
it to preserve dollar expressions literally in child arguments.

Properties use systemd's native names and values; at least one is required.
Unsupported properties, unavailable controllers, or an unavailable user
systemd manager fail without an unbounded fallback. The command runs in
the foreground and Ctrl+C stops it. Ordinary Bazel run-lock behavior remains;
Bazel's persistent server stays outside the new scope.

The CLI remains available for direct commands:

```sh
bazel_agent bazel run @rules_systemd//:systemd_run -- \
  --property=MemoryMax=128M --property=MemorySwapMax=0 -- /usr/bin/true
```

[ComfyUI](https://github.com/alwaldend/src/tree/master/tools/comfyui) declares its existing
model limits through this rule. The small `//test:probe`
fixture writes JSON evidence of real kernel limits, argument forwarding,
environment, and runfiles when supplied `-- --report <absolute-output-path>`.

## Development

Run `bazel_agent bazel build //...` and `bazel_agent bazel test //...` from
this module directory. `//:build_test` checks the public rule against the
runfiles fixture. Runtime verification requires a Linux user systemd manager
and delegated memory controllers; invoke `//test:probe` with an absolute
`--report` path to record actual kernel limits.
