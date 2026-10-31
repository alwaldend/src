## Context

AL already merges YAML, JSON, and Lua configuration through `al_proto.Config`;
the root `al.lua` is packaged by `//tools/al:config`. `bazel_agent tool` provides
an existing source-keyed executable-cache pattern, but its catalog is hard-coded
and tied to the Bazel runner.

## Goals / Non-Goals

**Goals:**

- Keep tool discovery in the same AL configuration used by command execution.
- Reuse the safest properties of the existing cache: source-keyed entries,
  exclusive publication, private permissions, and executable validation.
- Use a declarative Bazel-backed install method instead of hard-coded tool
  fields.
- Remove `bazel_agent tool` so there is one tool execution path.

**Non-Goals:**

- Do not resolve or install unconfigured tools from filesystem discovery.
- Do not create a new package-manager workflow.
- Do not change AL plugin lifecycle semantics.

## Decisions

- **Configuration schema**: add repeated `tools` and repeated `tool_methods` to
  `al_proto.Config`. A tool names a method; a method names a build
  implementation plus method-specific fields.
- **Default configuration**: `al tool --config ...` defaults to the repository
  root `al.lua`, found from the current directory upward.
- **Source key**: hash schema, platform, normalized tool and method fields, the
  absolute Bazel workspace path component, and the relevant configuration and
  Bazel input files. Changed configuration or Bazel files select a new key.
- **Cache root**: use `$XDG_CACHE_HOME/al/tools` (or `$HOME/.cache/al/tools`) by
  default; `AL_TOOL_CACHE` and `--cache-root` override it.
- **Method contract**: the initial `bazel` method builds `target`, locates the
  executable by `output`, and moves the selected executable into the entry as
  `bin`. Future methods may add fields without reusing the same schema.
- **Bazel invocation**: use the host `bazel` command with `--config=agent`, build
  the label, then inspect the requested executable in the workspace output tree.
- **Compatibility**: delete `bazel_agent`'s tool registry, execution, launcher
  probing, and cache code; migrate `mcp_cordis` and `repo_delivery` documentation
  to `al tool` declarations or ordinary Bazel run entry points.

## Risks / Trade-offs

- [Repository-wide configuration inputs can make cache keys broad] → Start with
  declared input paths and keep the key explicit rather than silently reusing
  stale executables.
- [A wrong Bazel output path could install the wrong executable] → Require an
  explicit executable path and validate regular-file and execute-permission bits
  before publication.
- [Removal breaks older launchers] → Document migration in the owning project
  READMEs and update the one repository launcher.

## Evidence

- 2026-09-09 20:27:04: `al tool ... repo_delivery provider` reported a cache
  miss for `sha256:15acc6185387be81fbbabbdc52e2c3e47e8bf8596e0e9d1e9c7c688301f1e881`,
  built `//tools/repo_delivery/cmd/repo_delivery:go`, installed a metadata
  entry, and returned provider JSON successfully.
- 2026-09-09 20:27:18: the same isolated cache returned the provider JSON on a
  cache hit with no stderr. The entry metadata identified kind `bazel` and the
  exact digest.
