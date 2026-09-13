---
title: Rules Terraform
description: Bazel rules and pinned provider installation for Terraform
---

This standalone module owns reusable Terraform rules. Terraform and provider
archives are declared Bazel inputs. Provider downloads belong to the Bzlmod
extension; Terraform uses a filesystem mirror containing only the providers
selected by its target.

The extension resolves one version per provider source. Consumers select a
provider by Bazel label. Explicit version override policy can be added at this
resolution boundary later; conflicting versions currently fail resolution.

Repository authentication and backend injection are supplied by the caller.
The module does not depend on the enclosing monorepo's AL or Vault packages.

## Declare providers

Use `terraform_providers` from `@rules_terraform//:extensions.bzl`. Each
`archive` tag requires a unique repository `name`, canonical lowercase
`source` (`hostname/namespace/type`), exact `version`, Terraform `platform`
(`os_arch`), immutable HTTPS `urls`, and SHA256 SRI `integrity`. Verify the
integrity against the publisher's release checksums before declaring it.

Import each generated repository with `use_repo`. Its public `:provider`
target exposes `TerraformProviderInfo`. Different platforms can share one
provider version; duplicate archives or conflicting versions fail before
repositories are registered. Provider downloads use Bazel's verified downloader
and repository cache, without executing host tools or querying a registry.

This repository's concrete declarations live in the
[Terraform dependency package](../../third_party/terraform/README.md).
The module is currently developed with a local Bzlmod override; no registry
release is implied.

## Run Terraform

```starlark
load("@rules_terraform//terraform:defs.bzl", "terraform_binary", "terraform_test")

terraform_binary(
    name = "tf_plan",
    srcs = glob(["*.tf"]),
    data = ["//modules/example:source"],
    providers = ["@my_provider_linux_amd64//:provider"],
    arguments = ["plan"],
)

terraform_test(
    name = "tf_fmt_test",
    srcs = glob(["*.tf"]),
    arguments = ["--direct", "fmt", "-check", "-recursive"],
)
```

`arguments` are fixed runner/Terraform arguments; command-line arguments are
appended. `srcs` and `data` declare configuration, local child modules, and files
read by the configuration. `chdir` defaults to the target's package. The
rules package ZIPs at
`HOST/NAMESPACE/TYPE/terraform-provider-TYPE_VERSION_OS_ARCH.zip` beneath a
target-specific runfiles mirror, following Terraform's
[filesystem mirror contract](https://developer.hashicorp.com/terraform/cli/config/config-file#filesystem_mirror).

## Command maps and wrappers

The same `terraform/defs.bzl` exports `terraform_binary_map`,
`terraform_target_binary_map`, and `terraform_test_map`. They preserve explicit
operation names such as `tf.plan`, `tf.apply`, and `tf_tests.fmt_test`.
Targeted maps produce plan/show/apply commands; their apply command requires
one saved plan. No unnamed apply alias is generated.

Maps run the Terraform rules directly by default. Callers can supply a
`wrapper` macro and its `wrapper_kwargs` to compose authentication or other
command setup without adding that framework to rules_terraform. The wrapper
receives the final target `name`, the inner executable in `args`, its runfiles
in `data`, and common target attributes. `wrapper_kwargs` cannot replace those
arguments or duplicate common attributes. Wrapped tests require a test wrapper.

Repository consumers load AL's generic `al_binary_run` or
`al_binary_run_test` and pass their `configs` and `run_args` in
`wrapper_kwargs`. Their plugin binaries remain in `data`. Provider-free tests
can use the Terraform test maps directly. Repository root command bindings
live with [the shared AL configuration](../al/README.md).

The default Terraform executable is the pinned toolchain in
`terraform/binary_toolchain.json`, acquired through `rules_binary_toolchain`.
The initial supported executable platform is Linux amd64. The optional
`terraform` executable attribute permits a caller-supplied Bazel tool.
`@rules_terraform//terraform:cli` exposes the pinned CLI for formatting and
other consumers that do not initialize providers.

Initialization and generated locks remain in runtime workspaces. No source
`.terraform.lock.hcl` is required. The [launcher contract](cmd/terraform/README.md)
describes working directories, direct commands, saved plans, and manifest-only
execution. Provider and backend operations can still use the network when
the operator invokes them; provider acquisition itself has no runtime registry
fallback. Configuration modules must be declared local inputs.
