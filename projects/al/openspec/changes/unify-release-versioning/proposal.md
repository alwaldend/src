## Why

Release publication and version calculation currently have separate owners and incompatible guarantees. Projects need independent releases using one calendar version scheme and an operator-run process that makes publication immutable until explicit deletion.

## What Changes

- Own AL as an independent project under projects/al while retaining the rules_al Bazel module identity and consumer-supplied repository configuration.

- Extract the standalone `rules_al` module and consolidate AL, release/versioning, guarded delivery, optional Git changelog parsing, and environment diagnostics there.
- Expose `al version`, `al release`, `al delivery`, `al bazel test`, and `al doctor` beside the existing `config` and `run` commands, without a `repo` namespace.
- Remove generic Bazel forwarding and stamped-build wrappers; ordinary Bazel reads user defaults and release configuration obtains identity from the installed AL status command.
- Keep source preparation Git-only and publish the prepared commit through its receipt. Remove delivery validation plans, the validation ledger, and redundant command and flag surfaces; operators run builds, tests, and lint separately.
- Add `al bazel test` to select tests affected by Git changes through Bazel's dependency graph, independently of source publication. Bazel owns test execution and caching.
- Move lifecycle skills to `al-versioning`, `al-delivery`, `al-workspace`, `al-rebase`, `al-bazel`, `al-checks`, and `al-openspec`; add `al-maintenance`, regenerate discovery, and migrate callers.

- Consolidate version calculation, stamping, source-ref publication, bundle generation, and SSH publication under `projects/al`.
- Bind first-party bundles to version, channel, clean source commit, filenames, sizes, and SHA-256 hashes.
- **BREAKING**: publish complete SSH releases atomically; identical retries succeed and conflicting existing releases fail.
- **BREAKING**: refuse legacy mutable-tag OCI publication; retain OCI manifest generation.
- Add explicit project-release deletion, refusing deletion of an active website; ordinary publication works again after deletion.
- Preserve calendar versions, independent project selection, third-party namespaces, and separate upload, extraction, and activation.
- Make artifact publication upload-only and expose optional website extraction and activation under `al release website`, each with an offline preview.
- Remove the unused full-history Git JSON command, rule, target, and website data input; retain release identity and optional changelog bundles.
- Replace the simulated source-ref lease with a real local operation lock and guarded Git ref updates; support patch releases on diverged release branches.

- **BREAKING**: expose `@rules_al//:defs.bzl` with `al_`-prefixed rule names and `Al`-prefixed providers, and name the lifecycle skills `al-<name>`.
- Restore AL’s visitor-facing landing page through the existing project registry and shared site.
- Remove retired versioning stubs, root compatibility labels, old CLI spellings, and unused OCI executor options; migrate callers to the canonical AL API.
- Fix destination-filtered activation and support verified payloads supplied through runfiles symlinks.
- Correct standalone command execution, redirected output, conflicting plugin exports, and direct rebase publication boundaries.
- Add an offline artifact publication preview and read-only release inspection for recovery without repeating a mutating operation.

## Capabilities

### New Capabilities

- `project-release-process`: shared operator workflow, source identity, independent project bundles, and release-ref publication.

### Modified Capabilities

- `ssh-release-deployment`: immutable atomic publication and explicit deletion.
- `project-bazel-agent`: ordinary Bazel invocation with user defaults, affected-test selection, and bounded AL diagnostics.

## Impact

Describe the consolidated command-line surface under `docs/product-description`,
including source-linked behavior, interruption and recovery, verification
checklists, and defect triage. Fix confirmed behavior defects and describe the
new inspection and preview commands without deploying releases. Assess typed
Bazel configuration providers as a future replacement for repository Lua authoring.

Release/versioning Go packages, Bazel targets and bootstrap, project bundle consumers, versioning skill and discovery. The former Git project is merged into this module, preserving used consumers and removing unused wrappers. No new external software or live deployments. Checked-in host provisioning installs AL instead of the retired runner; no host is changed by implementation.
