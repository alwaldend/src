## Context

See proposal.md for motivation. Existing low-level rules own execution. Packages
have custom playbooks, renamed assets, scoped Terraform operations and tests
with different data/providers or AL wrapping from operational commands.

## Goals / Non-Goals

Keep executable and test behavior stable while moving repeated target graphs
into small repository macros. Do not change Terraform configuration, Ansible
content, authentication, dependencies or deployed systems.

## Decisions

- `al_ansible_package` belongs in `tools/ansible/rules`; it composes source
  mappings, a package group, an Ansible binary and AL command wrappers. Explicit
  source lists and additional mapped groups retain package-specific inputs.
  Arguments omit the repeated executable; defaults describe the common deploy
  invocation. Separate visibility for the binary retains existing consumers.
- `al_terraform_package` belongs in `tools/al/al_terraform.bzl`; it composes
  existing operation and test maps. Tests default to direct offline execution;
  existing AL-wrapped tests opt in explicitly. Test overrides preserve custom
  data, providers, arguments and names. Specialized targeted maps stay explicit.
- Keep documentation and inventory export targets with their current owners.
  Hiding them in either macro would mix unrelated contracts.

## Risks / Trade-offs

- Macro defaults could change behavior: compare expanded rule attributes for
  every migrated package against the pre-change XML query, ignoring only source
  locations and generator metadata; inspect actual built launcher/config output.
- Package inputs vary: retain explicit globs, renamed assets and test overrides.
- Live commands can mutate hosts: build wrappers and run only offline checks.

## Migration Plan

Capture baseline, add macros, migrate callers, compare graphs, validate and
publish. Reverting the source commit restores the previous definitions without
any infrastructure operation.

## Verification

- Migrated 61 packages (18 Ansible declarations, 45 Terraform declarations).
- Compared all 1,666 expanded rules before and after: names, attributes,
  selectors, visibility and package mappings match. The comparison ignores
  generator metadata, explicit empty rename maps and order of Terraform data
  dependency labels. Buildifier sorted two distinct module labels in the
  Yandex Cloud package; no runtime argument or ordered source mapping changed.
- Built 63 selected packaged binaries and Terraform command targets; inspected
  development, image-generation and local-inference runfiles manifests and the
  development Terraform launcher configuration.
- All 73 selected tests passed, covering 45 offline Terraform checks and the
  root repository quality/Buildifier suites. The AL fixture retained its
  verified invocation receipt. Development Ansible syntax validation passed.
- Semantic lint passed for 1,667 targets. Documentation extraction exposed two
  missing existing AL Starlark library dependencies (Skylib paths and native
  binaries); declaring them at their owner fixed the failure.
- No live Terraform operation or Ansible deployment was executed.

## Continuation

Shared source refactor; publication uses the configured Forgejo origin and
master base through the al-delivery Forgejo procedure. Detailed command logs,
query snapshots and comparison results are in ignored out/package-macros/.
