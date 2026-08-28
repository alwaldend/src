---
name: full-repo-check
description: >-
  Run every normal Bazel build and test in the root and nested workspaces of
  this repository, preserve restricted command logs, diagnose every failure,
  and report the issues and fixes as a Markdown table. Use for a complete
  repository health check, release-readiness check, or build-and-test audit;
  use repo-bazel for ordinary package-scoped validation.
---

# Check the full repository

## Run the complete check

1. Read the root `AGENTS.md` and follow the `repo-bazel` skill.
2. From the repository root, run:

   ```sh
   bazel run --config=agent \
     //.agents/skills/full-repo-check:run_full_repo_check
   ```

   The runner executes `bazel build --config=agent //...` and
   `bazel test --config=agent //...` in the root workspace and three nested
   workspaces. It continues after failures so one run covers the whole
   repository.
3. Note the printed run directory. Read its `report.md` for the command matrix.
   Raw stdout and stderr are combined in its `logs/` directory with mode 0600.
4. Inspect each failed command's log with targeted searches and the narrowest
   useful context. Do not paste a raw log wholesale. In particular, redact
   credentials, secret-scanner candidates, private infrastructure values, and
   other sensitive data.

Do not add Build Event Protocol output, invocation environment dumps, or
machine-local settings. The Bazel logs are sufficient for diagnostics, and raw
BEP can contain the client environment.

## Diagnose every result

- Identify every failed build or test label printed by Bazel. If one root cause
  appears in multiple commands, group it into one row and enumerate all
  affected labels.
- Trace failures to repository code or configuration when possible. Use
  targeted Bazel queries and focused log excerpts; do not settle for the final
  `Build did NOT complete` line.
- Distinguish repository defects from missing local tools, unavailable external
  services, flaky tests, and other environment limitations.
- Give a concrete fix for each root cause. Name the file, dependency,
  toolchain, configuration, or test behavior that should change.
- Treat skipped or incompatible targets as coverage gaps, not passing results.
  List them after the issue table when Bazel reports them.
- Do not implement fixes unless the user also asks for them.

## Report the outcome

Return a Markdown table with these columns:

| Workspace | Phase | Affected target(s) | Root cause | How to fix | Evidence |
|---|---|---|---|---|---|

Use one row per distinct root cause. If every command passes, include one row
whose affected target and root cause are `None`. After the issue table, include
the eight-command result matrix from `report.md`, with exact commands, exit
statuses, durations, and the run-directory path.

This check covers Bazel's normal `//...` expansion with the repository's
default configuration plus `--config=agent`. Bazel excludes targets tagged
`manual`, skips incompatible targets, and does not expand optional
configuration matrices such as `--config=harbor`; disclose those limits in the
report rather than describing them as checked.
