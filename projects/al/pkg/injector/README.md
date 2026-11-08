---
title: Built-in resource injection
description: Generic resource graphs and invocation-owned resource lifetimes
---

## Configuration

AL evaluates selected `inject` configuration entries internally. No injector executable or
plugin registration is needed for the built-in `env`, `file` and `process` types.
`al_inject` is exported through `@rules_al//:defs.bzl`; its named
`res_<type>` maps compile into resources with `type`, `args` and ordered `deps`.
Declarations live in typed `InjectorCall.res` maps, separately from external
plugin calls. Resource keys own identity. Built-in env/file/process settings use generated
protobuf messages; Vault settings remain opaque arguments to their provider.

Each resource can supply JSON-compatible `data` alongside its type-specific
`args`. Environment and file templates read local input through `.Data`;
`.Last.Data` and `.Res[name].Data` read dependency results. Data supports nested
objects, lists, nulls and signed integers without rounding. Provider fetches
receive the same input as `ResourceFetchRequest.data`. File-specific `extra`
settings and `.Extra` have been removed.

The engine validates missing/cyclic dependencies before starting providers,
evaluates each selected node once and supplies ordered dependency results to
its fetcher. `Last` retains the final declared dependency. The root `inject` map
keeps named groups and their `res` maps. Explicit target extensions determine
resource override precedence; key order has no effect. Equal unrelated resource
definitions coalesce; conflicting unrelated definitions fail before startup.
Extension lookup traverses ancestor declarations without selecting their resources.
Environment values remain structured settings. Small declarations stay inline;
extract macros only for meaningful reuse or substantial readability gains.

External plugins marked `resource_provider = True` advertise types through
`PluginStartResponse.resource_types` at startup and implement `ResourceFetch`. AL starts
eligible providers for discovery only when selected resources require types beyond
the built-ins. Explicit plugin label/call selection retains its normal startup behavior.
Advertisements must be nonempty and unique across providers and built-ins.
Missing providers, capability collisions and resource failures prevent child
startup. Provider startup must initialize capabilities without fetching secrets;
AL fetches only selected resources through the generic RPC.

Fetch requests carry the resource name, type, opaque arguments and ordered
`ResourceResult` dependencies. Results carry data, file paths and environment
exports. `PluginStart` returns `PluginStartResponse` with an optional
`ResourceResult` in `result`; the envelope can gain startup fields independently
of resource output. `ResourceFetch` returns `ResourceResult` directly.
Go plugins may return `nil, nil` or a response without `result` when startup has
no resource output. An explicitly empty result retains its presence in
`PluginClient.StartResponse().GetResult()`. Present startup env/files/data
remain available through that result, and environment exports reach the child.
Plugins retain ownership of their startup files through shutdown. Startup data
does not implicitly add dependency nodes to the injection graph.

Data travels as JSON-compatible values; signed integers retain int64 precision.
Diagnostics omit argument/result values. Explicit injection exports override startup exports, including conflicting startup
defaults. Environment-only plugins return their
exports in `ResourceResult.Env` and need no resource service. Invalid capability advertisements roll back startup before the child runs.

AL keeps resources alive until the child finishes, including its graceful
cancellation. Shutdown stops built-in consumer processes, removes their files,
then stops provider processes. Providers drain requests before releasing their
own tokens/files. Failed startup and cancellation roll back attempted resources.
Temporary files use `0600`; SSH directories owned by the Vault provider use `0700`.
Deletion is unlinking, not secure erasure; forced host/process termination can
prevent cleanup.

The [Vault provider](../../plugins/vault/README.md) supplies `kv`, `op`, `oidc`,
`vault_env` and `vault_ssh`. Core injection imports no Vault provider schema or
handler implementation. No global provider registration or live deployment is
performed by building configuration targets.
