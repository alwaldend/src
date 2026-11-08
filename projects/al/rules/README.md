---
title: AL Bazel rules
description: Configuration, command, Git, and release rules
---

AL rules and macros use the `al_` prefix; AL providers use the `Al` prefix.
Their source filenames follow the same prefix in snake case. The public API is
exported through `@rules_al//:defs.bzl`, with usage documented in
the [project README](../README.md#bazel-module-api).

## Configuration targets

Declare small `al_vault_auth`, `al_vault_conn`, `al_plugin` and
`al_plugin_call` targets, then combine their labels with `al_config(configs = ...)`.
Keep short, single-use declarations in BUILD, including small dictionaries and
direct calls to existing macros. Extract namespaced local `defs.bzl` macros only
when they provide significant reuse or make complex graphs, construction logic,
long templates or bulky settings substantially easier to read. A single KV
source and a few environment exports can stay inline. Put genuinely common
patterns in reusable target macros with their owner; avoid wrappers that merely
bind a few fixed attributes.

```starlark
load("@rules_al//:defs.bzl", "al_config", "al_vault_auth")

al_vault_auth(name = "shared_auth")

al_vault_auth(
    name = "component_auth",
    extends = ":shared_auth",
    approle = "example",
)

al_config(
    name = "config",
    configs = [
        ":component_auth",
        ":shared_auth",
        "@rules_al//plugins/injector:config",
    ],
)
```

`configs` combines definitions independently of list order. An atomic target's
`extends` selects another atomic target of the same kind. For the same runtime
name, the descendant replaces the complete definition; shared ancestors are
included once. Unrelated definitions with the same kind and name fail analysis.
In the example, both auth targets name `default`, so the output contains one
auth method using the component AppRole. Command wrappers combine multiple
config targets with the same rules before passing a generated file to AL.

Extensions between differently named plugin calls retain both calls and place
the descendant after its ancestor. Use this only for deliberate resource or
backend overrides. Independent definitions are sorted deterministically by
extension depth, kind and runtime name. Aggregate reference lists can follow
normal formatter ordering; builds never authenticate or start plugins.

`al_config(config = {...})` remains available for small inputs. Its `vault_auth`,
`vault_conn`, `plugins` and `plugin_calls` sections are dictionaries keyed by
runtime name, with no repeated names inside their values. Entries may use an
explicit `extends` label to override an inherited definition. Large dictionaries
should be split into targets rather than moved wholesale into another file.

Repository components compose `//infra:al` and have their auth target extend
`//tools/al:auth`. A standalone consumer can compose
`@rules_al//plugins/injector:config` to register the injector. Reuse shared
configuration targets instead of copying registrations.

Configuration targets carry plugin executables, their transitive runfiles and
TLS files. An `al_binary_run` consumer only needs `configs = [":config"]`;
do not repeat those plugins in its `data`. Other Bazel consumers can place the
configuration target in `data` to receive the same runtime closure.

`args` contains ordinary JSON-compatible plugin arguments. `argv` contains an
executable's command-line argument list. Plugin arguments support up to 64 nesting levels, including distinct
empty maps and lists and nested null values.

The injector's `args.res` is a dictionary of named resource objects. `env` and
other resource settings remain objects. Dependencies are ordered lists of
resource names: `Last` selects the final listed dependency. Selected call maps
merge in the generated call order, with descendants replacing ancestor resource
definitions of the same name. Missing and cyclic dependencies fail before the child starts.

Vault connection attributes `ca_cert`, `client_cert` and `client_key` select
files by label. Runtime TLS initialization resolves directory or manifest
runfiles. Absolute paths remain filesystem paths; installed commands without
runfiles resolve relative TLS paths against their working directory. Secret
values belong in runtime injection, not build inputs.

A call can declare `substitutions = {"//data:public_file": "@PUBLIC_FILE@"}`.
The action reads this explicit public input, trims surrounding whitespace and
JSON-escapes its contents before replacing the marker in argument string
values. This is suitable for public certificate authorities.

`al_config` compiles to AL JSON. AL loads JSON and YAML; Lua configuration and
the former `srcs` interface have been removed. All four atomic target types provide `AlConfigInfo` and runtime files.

## Command selection

Pass `plugin_labels = {"tf": "setup", "pve_login": ""}` to AL command wrappers.
Keys select labels; nonempty values match an exact value and an empty value
matches the label's presence. The wrapper encodes these into AL CLI arguments.
Use `run_args` only for other AL options and `args` for the child command.

`al_injector_call`, exported through the public API, accepts named `res_env`, `res_file`,
`res_kv`, `res_op`, `res_oidc`, `res_vault_env`, `res_vault_ssh` and `res_process` maps. Entries contain
structured type settings and optional `deps`, `vault_auth` and `vault_conn`
metadata. It builds the existing `args.res` graph and selects the injector.
See the [injector API](../plugins/injector/README.md#configuration).

Plugin calls default their runtime identity and selected plugin name to the
Bazel target name. Use `config_name` or `plugin` when those names must differ.
Injector calls default their identity the same way and always select `injector`.
In the small dictionary API, the entry key supplies the default plugin name;
generated private fragment names do not become runtime identities.
