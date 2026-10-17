## 1. Implement snapshot state export

- [x] 1.1 Add native CLI E2E onboarding and failure cases before implementation; emit config, state and plan artifacts.
- [x] 1.2 Generate bindings during config projection, expose optional state export through App and CLI, and verify native E2E coverage including source preservation, aliases and failed output.

## 2. Document and deliver

- [x] 2.1 Document full onboarding and state output guarantees; verify examples against the native CLI workflow.
- [x] 2.2 Prepare scoped validation and publication inputs requiring affected tests, semantic lint and repository quality; record exact-candidate outcomes in delivery receipts.

## Evidence

The seven-target native/CLI/plan/reconcile test run passed on 2026-10-18 in
48.758 seconds. Onboarding artifacts include TOML bindings and note configs,
plus an empty JSON plan after state installation, for archives and databases.
Ignored `out/export-state/` contains command logs and delivery receipts; those
receipts own final candidate validation and publication state.

Review finding `PRRT_kwDOOgezIM6tkg2G` identified unresolved symlink aliases to
absent state destinations. Added native E2E failure cases before fixing the
shared output resolver and canonical destination comparisons.

The final ownership boundary keeps binding construction in `state.Export`.
Config projection provides paths and native IDs; its note-type binding helper
remains private to state. App exposes only the regular `ColExport` operation.

Review finding `PRRT_kwDOOgezIM6tknLc` added the inverse destination check:
state files cannot be nested beneath generated config files. Native CLI E2E
requires rejection before writing the config tree.
