## Acceptance

The existing E2E assertion was written before the runtime fixture change.
It first failed with `undefined !== 13`, then passed after the runtime import
and dependency wiring. The retained `node-result.json` reports Node 24.13.0,
Three.js revision 186, and vector length 13. The browser artifact remains a
320 by 120 canvas with visible coordinates; the invalid program still produces
only compiler diagnostic TS2322.

The generated pnpm diff adds only the approved `three@0.186.1` importer,
integrity entry, and empty dependency snapshot. All 366 existing package
entries are unchanged. Lifecycle hooks remain disabled.

## Checks

- `bazel run //:format`: passed; no unrelated source changes.
- `bazel test //:repo_quality_test //tools/js/test/typescript:typescript_test
//tools/mermaid/cmd/mmdc:mmdc_test //tools/mermaid/test/renderer:all
//projects/alwaldend.com:site_test //projects/alwaldend.com/test/mermaid:all`:
  all 36 tests passed.
- `bazel build --config=lint //tools:all //tools/js:all
//tools/js/test/typescript:all //tools/js/test/typescript/node:all`:
  all 15 targets passed.
- `OPENSPEC_PROJECT=tools/js bazel run //tools/openspec -- validate
add-threejs-runtime --strict`: passed.
- `git diff --check`: passed.

Raw command logs and delivery receipts remain in ignored
`out/three-runtime/`. The test retains `node-result.json`, `browser.png`,
`type-error.log`, and `verification.json` with hashes in its Bazel test outputs.
