# Requirements and constraints

[Back to durable goal](./)

## User requirements

- The project is named `mcp_cordis` and lives under `projects/`.
- Reuse Cordis itself rather than reimplementing its lifecycle architecture.
- Keep reusable runtime code under the project for later reuse.
- Keep disposable runtime code under `out`.
- Seed the project with packages based on recurring past-session needs.
- Security hardening is not the priority on this dedicated LLM machine.
- Codex must load the MCP from project-scoped configuration so different
  clones and linked worktrees use their own source and disposable state.

## Repository constraints

- Use `bazel_agent` for every Bazel command.
- Pin external dependencies reproducibly and retain required notices.
- Keep disposable task scratch under `out/mcp_cordis`; maintain this explicitly
  requested reusable goal under `projects/mcp_cordis/goals`.
- Preserve unrelated working-tree changes.
- Prefer a narrowly scoped project and focused validation.

## Assumptions

- "Based on past sessions" means extracting generic recurring workflows, not
  copying private conversation text, credentials, or secret-bearing data.
- Host-only JavaScript packages are the initial scope; browser UI packages are
  not required for the first working server.
- A stable gateway invocation tool is required because MCP clients differ in
  when they consume tool-list change notifications.

## Requirement changes

- 2026-08-30: project name changed from `agent_extension_host` to
  `codex_cordis`, then finally to the client-neutral `mcp_cordis`.
- 2026-08-30: storage clarified as two-tier: reusable project packages and
  disposable `out` packages.
- 2026-08-30: the user promoted the runtime-extension goal directory from
  disposable `out` coordination to durable project documentation for future
  reuse.
- 2026-08-30: remote `master` advanced and the user explicitly requested a
  rebase before implementation continued. The task commit was rebased from
  `775f44d3b56146005e44980f3cf948785f963ba0` onto
  `7ad2704cd27757355ab36ec8eb1bb27ef9e1d91d`; all prior checks were treated as
  invalid until rerun.
- 2026-08-30: after publication, the user explicitly expanded the delivery to
  include only the `projects/agents` subtree changes from PR 24. PR 24 changes
  four files there: `bazel-agent`, `goal`, and the new `decision-review`
  package. Import those changes three-way against current `master`; do not
  import its unrelated render or infrastructure content.
- 2026-08-30: the user rejected the custom per-package manifests and requested
  standard solutions. Runtime persistence must use official Cordis loader
  entries and `cordis.yaml`; normal reusable files and Git replace committed
  content-addressed source history.
- 2026-08-30: after review-driven fixes exposed additional correctness gaps,
  the user required `repo-delivery` to invalidate prior correctness verdicts
  after code changes and require fresh diff-focused scrutiny beyond green
  tests.
- 2026-08-30: the user rejected the HMR race-handling complexity and explicitly
  required the simplest robust MCP wrapper for loading DeepSeek/Cordis plugins.
  Source-tool success therefore means validated atomic persistence and an
  official Cordis load/reload request, not a custom synchronous activation
  transaction or on-disk rollback protocol.
- 2026-08-30: the user required per-workspace Codex loading. The checked-in
  `.codex/config.toml` must resolve the active linked worktree; the launcher
  must release Bazel's output-base lock before the MCP begins serving stdio.
- 2026-08-30: remote `master` advanced again to
  `d29f9d471ea467e8dfc75db4eedeedbbae43dc2d`. The user requested another
  rebase and incorporation review. Preserve the new `projects/goal` redesign
  and do not resurrect its deleted predecessor under `projects/agents`.
- 2026-08-30: after the delivery adapter refused the exact nine-commit task
  range, the user explicitly authorized extending `repo_delivery` with a
  guarded exact-head consolidation operation and then repeated the request to
  rebase, review, and incorporate the upstream changes.
