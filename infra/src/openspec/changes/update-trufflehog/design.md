## Context

See `proposal.md`. The module-freshness test must inspect workspace files and
carries `no-sandbox`. Action-query evidence showed that this tag also reached
`GoCompilePkg`, `GoCompilePkgExternal`, `GoTestGenTest`, `GoLink`, and `RunNogo`.
Compilation consequently used the host's quota-limited `/tmp`.

## Goals / Non-Goals

Keep Go build actions sandboxed while retaining the test runner's existing
local-file access. This change does not modify host mounts or temporary-directory
configuration and does not propagate ambient temporary paths into actions.

## Decisions

Proceed with the verified `--modify_execution_info=(Go.*|RunNogo)=-no-sandbox`
setting in shared `tools/bazelrc/project.bazelrc`. Both root and nested modules
already import that configuration. The mnemonic scope preserves `TestRunner`
and unrelated execution metadata.

Leaving the setting per invocation would leave ordinary builds exposed to the
observed quota failure. Passing a host `TMPDIR` into all actions would violate
the existing environment contract and couple builds to a caller-owned path.
The user explicitly authorized the shared action-metadata fix.

## Risks / Trade-offs

- A Go build action with undeclared host inputs may now fail in its sandbox;
  validate actual compilation and retain the local test exception.
- Mnemonic-based selection depends on the pinned rules_go action names;
  inspect the default action graph and confirm actual local runner selection.

## Validation and continuation

Check default action metadata, execute a debug compilation with disk and remote
cache reads disabled, and confirm the compiler uses `linux-sandbox` while the test runner retains
`no-sandbox`. Run the environment-policy test, required quality and lint gates,
skill checks, and a fresh secret scan against the prepared candidate.
Update the existing task PR only after these checks pass.
