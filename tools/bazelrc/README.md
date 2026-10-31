---
title: Bazelrc
description: Bazelrc files
---

## Usage

```sh
bazel run //tools/bazelrc:preset.update
```

The `agent` profile must not propagate ambient `TEMP`, `TMP`, or `TMPDIR`
values into repository rules, actions, host actions, or tests. Those execution
contexts use Bazel-managed temporary storage; host tools own any explicit
task/run scratch they require.

Go build actions (`Go*` and `RunNogo`) remove an inherited `no-sandbox`
execution requirement through the shared `--modify_execution_info` setting.
This keeps compilation in Bazel-managed sandboxes even when its owning test
needs local workspace access. `TestRunner` retains its existing requirements,
so the test can still inspect workspace files. The shared root and nested
workspace configuration both apply this setting.

Agent tests use Bazel's normal result cache. Tests that inspect undeclared
workspace files must use the `external` tag to bypass local test-result reuse;
`no-cache` alone does not disable that cache.

Agent commands use Bazel's default lockfile update mode. Review and commit
generated lockfile changes with dependency changes. CI retains strict lockfile
checking; use `--lockfile_mode=error` for an explicit reproducibility check.

The shared CI profile works in root and nested workspaces. The root-only
`//tools/bazel_configs:ci_flag` belongs in `root.bazelrc`; nested workspaces
import the common CI flags without resolving that root package.

In the root workspace, lint actions fail on violations. Their failure setting
is identical across normal and lint configurations to preserve the analysis
cache when switching modes. `--config=lint` enables the linter aspects and
requests only lint reports and skill validation, not ordinary target outputs.

## Repository caching and protobuf

Bazel 9 enables the local repository contents cache by default under
`{--repository_cache}/contents`. It reuses eligible fetched repositories across
workspaces; the shared configuration preserves this default location.

The shared startup configuration also enables
`--experimental_remote_repo_contents_cache`. It uses configured remote or disk
cache backends for reproducible repository rules. It does not provision a cache
service or configure an endpoint. Existing cache upload policy still applies.
Because the option is a startup flag, changing it restarts the Bazel server.

The root workspace enables
`--@protobuf//bazel/toolchains:prefer_prebuilt_protoc` to use the pinned protobuf
module's prebuilt compiler where supported, with source compilation as the
fallback. The setting is root-only because standalone nested modules can use
different protobuf versions. User overrides remain last.

## Links

- Rules: https://github.com/bazel-contrib/bazelrc-preset.bzl/

The generated preset omits the optional Cloudflare module mirror because it
slow-stalls on the CI network. Bazel uses the source URLs and mirrors declared
by each registry entry. Configure this omission in `BUILD.bazel`, then regenerate
`preset.bazelrc`; a later empty flag does not clear an accumulated mirror.
