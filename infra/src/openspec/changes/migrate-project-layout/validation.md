# Migration validation

The source migration covers 575 tracked inputs in nine product projects and
42 tool components. The inventory confirms every input exists at its new path
and the legacy directories are absent. Active reference searches found no
remaining consumers of those legacy paths. Historical change archives and
upstream source references retain their original spelling.

## Checks and representative outputs

- Eleven consumer builds passed: both Android APKs, the Qt application, the
  Nexus plugin, the CI platform backend, and six Go or Python execution tools.
- All 79 explicitly selected root-workspace tests passed, including command
  behavior, packaging, Gradle catalog freshness, the Android skill's offline
  configuration check, and OpenSpec validation across component workspaces.
- All 42 tests passed across the six affected standalone modules:
  rules_template, rules_binary_toolchain, rules_dnscontrol, rules_skills,
  rules_openspec, and rules_promptfoo.
- The Autoscroll wheel contains its importable package, launcher modules,
  internal implementation, resources, type marker, and console entry point.
  Inspection caught and corrected its old binary-based source selection and
  package prefix. The manifest is in `out/project-layout/wheel-manifest.json`.
- The migrated SRI command produced the SHA-256 integrity value independently
  calculated from the root README. Its output is in
  `out/project-layout/sri-output.txt`.
- All 22 moved PNGs, three SVGs, two text resources, and two active firmware
  patches retain their original bytes. The retained firmware development patch
  changes only its regeneration output path.
- Repository formatting completed. The guarded delivery plan owns final
  exact-candidate repository quality and semantic lint results.
- Semantic lint caught the browser userscript's implicit target retaining its
  old directory name. Renaming it to `download_submissions` fixed the release
  consumer; the remaining affected-package lint then passed.

The APKs were compiled; no device installation or live deployment was required
or performed. Dependency pins remain unchanged.

## Review follow-up

Review found that moving the CI platform's development configuration left
relative paths unresolved. The Compose initialization bind now selects the
canonical SQL file. Air now selects the relocated controller directory,
generated documentation directory, frontend directory, backend command package,
and development configuration. Its generated-output exclusions follow those
same directories.

The pinned Go tool resolves Air's corrected backend package. The canonical
backend build, existing backend tests, and component specification validation
passed together: 73 tests. Source resolution checks confirm that the SQL bind
and Air inputs exist. The full Air/frontend pipeline was not run because Air,
Swag, and Yarn are unavailable; the project's abandoned status and incomplete
frontend guarantee remain unchanged. No services were started.

The user confirmed that `tianyi.yml` and `tianyi.dev.yml` contain development
fixtures and requested scanner ignores. Each JWT secret and database/Redis
password has a same-line `trufflehog:ignore` comment. The pinned YAML parser
confirmed that both documents still parse to exactly their prior values.
The scoped pinned TruffleHog filesystem scan of those two files returned zero
findings with credential verification and update checks disabled. This does
not assert a clean historical repository scan or suppress entire files.

Follow-up evidence is under `out/project-layout/follow-up/`; the guarded
delivery receipt owns final candidate validation and publication results.

## Complete repository verification

The full root and nested Bazel matrix initially found three unrelated failures.
The user requested their repair. The recommender's missing-CSV fixtures now
derive canonical manifest keys through the public directory-based runfiles
resolver. The article pipeline WebP is refreshed through its owning updater
and visually inspected. The complete WASM application suite uses eight
isolated shards, retaining every workflow and the existing 900-second deadline.
The owner-local changes record these repairs.

All three affected E2E targets passed together. The recommender JSON artifact
confirms both missing CSVs, missing/corrupt models, and normal HTTP behavior.
The generated image matches the inspected Git LFS asset. All eight unfiltered
WASM shards passed. Scoped semantic lint and strict OpenSpec validation passed.

The complete normal matrix then passed all 32 commands: `al bazel
build //...` and `al bazel test //...` in the root and 15 nested
workspaces. All 374 root tests passed. Normal expansion excludes manual and
incompatible targets and optional configuration matrices. Full-check reports
under `out/al-checks/` and the guarded delivery receipt own final
exact-candidate results; `out/project-layout/full-repair/` retains repair
artifacts and the selected validation plan.

## Continuation and procedure observations

Use `out/project-layout/full-repair/prepare.json` and the task-local `repo_delivery`
launcher to resume publication. Do not repeat source moves. The initial test
invocation named an obsolete Gradle verification target from the Android
skill; target discovery confirmed only catalog freshness is exposed by the
converter. The skill now documents Gradle's verification ownership correctly.

Git diff and status operations can run Git LFS clean filters, which write to
the managed shared Git directory. Use the permitted managed-cache execution
context for those operations. Large initial reads were truncated; bounded
inventories and per-phase logs kept subsequent inspection reproducible.
