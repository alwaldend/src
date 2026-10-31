## Why

Release publication and version calculation currently have separate owners and incompatible guarantees. Projects need independent releases using one calendar version scheme and an operator-run process that makes publication immutable until explicit deletion.

## What Changes

- Extract the standalone `rules_al` module and consolidate AL, release/versioning, guarded delivery, Git metadata, and the Bazel agent runner there.
- Expose `al version`, `al release`, `al delivery`, `al bazel`, and `al doctor` beside the existing `config` and `run` commands, without a `repo` namespace.
- Move lifecycle skills to `al-versioning`, `al-delivery`, `al-workspace`, `al-rebase`, `al-bazel`, `al-checks`, and `al-openspec`; add `al-maintenance`, regenerate discovery, and migrate callers.

- Consolidate version calculation, stamping, source-ref publication, bundle generation, and SSH publication under `tools/rules_al`.
- Bind first-party bundles to version, channel, clean source commit, filenames, sizes, and SHA-256 hashes.
- **BREAKING**: publish complete SSH releases atomically; identical retries succeed and conflicting existing releases fail.
- **BREAKING**: refuse legacy mutable-tag OCI publication; retain OCI manifest generation.
- Add explicit project-release deletion, refusing deletion of an active website; ordinary publication works again after deletion.
- Preserve calendar versions, independent project selection, third-party namespaces, and separate upload, extraction, and activation.
- Replace the simulated source-ref lease with a real local operation lock and guarded Git ref updates; support patch releases on diverged release branches.

- **BREAKING**: expose `@rules_al//:defs.bzl`, remove `al_` public rule names, and name the lifecycle skills `al-<name>`.
- Fix destination-filtered activation and support verified payloads supplied through runfiles symlinks.

## Capabilities

### New Capabilities

- `project-release-process`: shared operator workflow, source identity, independent project bundles, and release-ref publication.

### Modified Capabilities

- `ssh-release-deployment`: immutable atomic publication and explicit deletion.

## Impact

Release/versioning Go packages, Bazel targets and bootstrap, project bundle consumers, versioning skill and discovery. The former Git project is merged into this module, preserving used consumers and removing unused wrappers. No new external software or live deployments. Checked-in host provisioning installs AL instead of the retired runner; no host is changed by implementation.
