---
title: Injector
description: Secret injector
languages:
  - go
tags:
  - vault
---

Shutdown first drains plugin requests, then stops and waits for resource
processes, revokes invocation-owned Vault credentials, and deletes temporary
files and SSH key directories. Failed cleanup is reported. Registration after
shutdown is rejected; a fetcher removes any unregistered temporary material.
Temporary files use mode `0600` and directories use `0700`. Deletion is filesystem
unlinking, not secure erasure, and cannot run after `SIGKILL` or host failure.

`no_auth` explicitly sets the injected `VAULT_TOKEN` to an empty value so an
inherited token is overridden. This does not remove the user's token-helper
file; commands that independently consult that helper may still authenticate.
Template errors and OIDC status errors omit input and response contents.
OIDC requests honor cancellation and do not follow redirects.

## Configuration

`al_injector_call`, exported by `@rules_al//:defs.bzl`, accepts named maps for
`res_env`, `res_file`, `res_kv`, `res_op`, `res_oidc`, `res_vault_env`, `res_vault_ssh` and `res_process`.

```starlark
al_injector_call(
    name = "github",
    labels = {"tf": "main"},
    res_kv = {"github": {"mount": "secrets", "path": "component/github"}},
    res_env = {"GITHUB_TOKEN": {"deps": ["github"], "value": "{{ .Last.Data.github_token }}"}},
)
```

Each map key owns the resource name. Its structured settings may include `deps`,
`vault_auth` and `vault_conn`; the macro places those fields on the resource and
wraps the remaining fields under the selected type in `args.res`. Resource names
must be unique across groups. Strings are not shorthand for settings objects.
Typed fields such as `process.name` and `oidc.name` retain their own meaning.
The macro consumes resource maps from kwargs by prefix and forwards other call
attributes; adding a runtime resource type needs no new macro parameter.
The runtime schema remains a named resource map; callers do not repeat
`plugin = "injector"` or the `args.res` wrapper.

The runtime call name defaults to the target name, with `config_name` as an
explicit override. The selected plugin is always `injector`. Generic
`al_plugin_call` targets default both their call and plugin names to the target;
set `plugin` when the registration has a different name.

Dependencies are explicit ordered lists: `Last` refers to the final listed
dependency. Selected call maps merge in generated call order, with descendants
replacing ancestor resource definitions of the same name. Declare intentional
precedence through `extends`; aggregate input order has no effect. Environment
exports merge once per evaluated resource, so a dependent takes precedence over
its dependencies. Missing or cyclic references fail before the child starts.

Keep short graphs inline; extract local macros only for significant reuse or a
substantial readability benefit. Compose small targets through `al_config`.

The injector runs as the external @rules_al//plugins/injector executable. AL owns its source and tests. Its local end-to-end fixture runs actual AL/plugin/child processes against synthetic loopback-only Vault responses and retains dictionary-results.json in the Bazel undeclared test outputs. The parent repository also retains a real Vault dev-server integration target at //tools/al:injector_dev_server_test, using the repository's pinned Vault binary. The standalone AL module needs no Vault server executable dependency.

Compose `@rules_al//plugins/injector:config` to reuse the injector registration.
The configuration supplies its executable and runfiles transitively; consumers
do not need a separate plugin entry or plugin `data` dependency. Repository
components receive this registration through `//infra:al`.
