## Context

The user clarified that component specifications belong to their projects
and repository evolution belongs to `infra/src`. The previous candidate is
`8e267dc1c3d8463535fe5a9a80ee3f5a4c94b98e`; the initial goal migration's
historical source is `550d7e79b1f5fdbc2b6017b75178471d6914082f`.

## Decisions

Proceed with one native `openspec/` workspace per owner. Keep capability names
and historical deltas unchanged. The repository's own workspace is
`infra/src/openspec`; no root specification directory remains. The CLI selects
`infra/src` by default and accepts the owner's repository-relative directory
through `OPENSPEC_PROJECT`.

A centralized directory with links would preserve the wrong ownership.
Invoking native OpenSpec per owner supports ordinary local edits and archives.
The main integration risk is omitting standalone Bazel modules from validation.
Existing module mappings expose their declared source filegroups to the pinned
CLI tests; sandboxed probes verified both a regular and an external module's
single spec. No JavaScript dependency is added to those project modules.

## Acceptance and continuation

The history relocation check verified all 222 original files (834,320 bytes)
against their source commit. Migrated acceptance and execution metadata remain
unchanged. Reimu Fumo is still open and blocked.

The coordinated package check passed all 59 tests, including all 49 owner
workspaces, archive task checks, routing regressions and Reimu evidence checks.
Representative native CLI output selected `infra/src`, `infra/vault` and the
standalone `projects/rules_skill` owner; Reimu still reported 0 of 12 tasks
complete. Independent review found all 636 live Markdown links resolvable and
all 49 workspace mappings complete. Frozen original history was excluded from
link rewriting and independently checked byte for byte.

Acceptance supports archiving this repository layout change. Final delivery
still verifies the formatted candidate and generated catalogs through the
repository procedure. Temporary logs and exact candidate receipts live in
`out/openspec-locality/`; checked-in acceptance summaries do not validate a
future implementation revision.
