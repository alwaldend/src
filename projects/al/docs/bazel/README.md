---
title: Bazel setup, affected tests, and diagnostics
description: Direct Bazel commands, Git-aware test selection, and local configuration reports
statuses:
  - active
languages:
  - go
  - sh
tags:
  - bazel
  - agent
---

Run Bazel directly from the owning workspace:

```sh
bazel test //path/to/package:all
bazel build --config=release //projects/autoscroll/releases:bundle
```

The repository's Bazelisk configuration pins Bazel. Per-user output and resource
defaults belong in `~/.bazelrc`; the [repository setup guide](../../../alwaldend.com/content/docs/misc/repo.md)
describes the agent setup and its managed host configuration. No AL command or
named agent profile is required to build or test. Bazel owns argument parsing,
progress, signals, and exit status.

The root release configuration enables stamping through the installed AL command.
See [Build one project](../../README.md#build-one-project) for prerequisites and
source-version selection. It does not start another Bazel client during a build.

Bazel actions use declared outputs and Bazel-managed temporary storage, and
tests use their test temporary-directory contract. Repository-updating host
tools that need scratch accept or derive an explicit task/run path under
`out/<task>/`; do not propagate ambient `TMPDIR`, `TMP`, or `TEMP` to the whole
build.

## Test affected targets

Use `al bazel test` to select tests affected by changes since an explicit Git
base, then run those tests with ordinary Bazel:

```sh
al bazel test --base master --dry-run
al bazel test --base master
al bazel test --base master //projects/example/...
```

Target patterns must be absolute; the scope defaults to `//...`.
`--workspace PATH` selects another Bazel workspace, resolving a relative path
from the current directory. Otherwise discovery starts from
`BUILD_WORKSPACE_DIRECTORY` when present or the current directory and finds the
nearest workspace root. The base is resolved in the enclosing Git repository
and must have one merge base with HEAD. Selection includes changes from that
merge base through the working tree, staged edits, and nonignored untracked paths.

Bazel queries map changed source files to their reverse dependencies and tests
within the requested scope. Deleted, unknown, configuration, or outside-workspace
paths select the full test scope when a narrower set cannot be established. A query failure
fails the command; partial query results are not accepted as a complete test
selection. The scope belongs to one workspace; check other affected nested
workspaces separately.

When selection narrows the requested tests, work needed only by unrelated tests
stays outside the test invocation. Missing or evicted cache entries for that
unrelated work therefore do not trigger rebuilds or reruns. The dependency query
still loads the requested graph, and selected tests use Bazel's normal cache
behavior. No speedup is guaranteed for a particular change.

The command prints one JSON selection on stdout before any tests run. Query
diagnostics and all test output go to stderr. Its `bazel/test_selection/v1`
report contains the resolved workspace/base/merge base, requested `scope`,
Git-root-relative `changed_paths`, `selection`, `reasons`, `tests`, and `dry_run`.
`selection` is `affected`, `full_scope`,
or `no_changes`; fallback reasons identify the paths that prevented narrowing.

No changed paths means no Bazel query or test invocation. An empty test selection
also succeeds without invoking tests. Wildcard scopes exclude root targets with
the literal `manual` tag before expanding test suites. An explicitly named test
or suite remains eligible, as do the members listed by an eligible suite,
including manual-tagged members. Ordinary Bazel test filters can further filter
the selected labels.

`--dry-run` skips test execution but can still run dependency queries. Normal
test failures return a nonzero status; the preceding JSON is a selection report,
not a passing-test receipt. This command selects tests; it does not implement a generic
Bazel passthrough, prepare a delivery candidate, or record a validation ledger.
Use ordinary Bazel commands for builds, other subcommands, and custom invocations.

## Diagnostics

`al doctor --workspace-root PATH --task-scratch out/<task>/<run>` is a read-only
JSON diagnostic. It reports AL and built-source identity, declared Bazelisk
pins, platform, a configured rc-file inventory, and
task scratch classification without dumping the environment or launching Bazel.

The workspace root defaults to the current directory. Select a different layout
with `--bazelisk-config PATH`, `--source-binary PATH`, and repeated `--rc-file PATH`.
Relative inputs resolve against the selected workspace root. For example:

```sh
al doctor --workspace-root /absolute/workspace \
  --bazelisk-config config/bazelisk.conf \
  --source-binary build/al \
  --rc-file .bazelrc --rc-file "$HOME/.bazelrc" \
  --task-scratch out/task/run
```

The pin-file default is `.bazeliskrc`. Without a source-binary override, comparison
uses `bazel-bin/external/rules_al+/cmd/al/al_/al` under the selected workspace;
the reported source target remains `@rules_al//:al`. A missing comparison file
reports `source.reason: unavailable`; it does not prove the installation is
current. An available comparison reports its resolved absolute path and digest.

The rc inventory defaults to `.bazelrc`. Explicit `--rc-file` entries replace
that default and are reported relative to the workspace when possible. Doctor
does not read those files, verify their existence, or discover Bazel's imports;
the list describes the supplied configuration and does not establish effective
Bazel options. Reading the required pins or available source
binary can fail the command; an unavailable comparison or conflicting scratch
classification remains a report finding.

## Installation

From this repository's root, install or update AL with:

```sh
bazel run //tools/al:install
```

The shared binary installer atomically copies the executable to
`~/.local/bin/al` with mode 0755. It accepts
`-- --directory /absolute/path` for an isolated destination. After every
code update to this project, rerun the install target; the installed binary is
not updated automatically.
