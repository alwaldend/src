## Why

The user requested that the repository's Bazel rule projects become tools and
that their dedicated landing sites be retired. The relocation must preserve
the rule modules and public APIs while removing the selected sites and their
infrastructure through the owning declarative workflows.

## What Changes

- Move all twelve `projects/rules_*` modules to corresponding `tools/`
  directories, preserving module names, public rule interfaces, specifications,
  and supported consumers.
- Update root and nested workspace references, generated registries, skill
  discovery, documentation, and applicable boundary rules for the new owners.
- **BREAKING**: retire the eleven existing rule-project landing repositories,
  Pages sites, and matching DNS records. Remove their site configuration and
  shared repository-catalog entries so GitLab does not import retired sites.
- Execute the explicitly requested retirement using reviewed plans whose
  complete change sets match the selected resources. Preserve unrelated
  repositories, infrastructure, memberships, and AppRoles.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `repository`: define the location and preserved build interfaces of the
  repository's Bazel rule modules, and their exclusion from dedicated project
  landing publication.

## Impact

Affected source owners include the twelve relocated modules, `projects/`,
`tools/`, root Bazel configuration, `infra/repos`, `infra/github`, and
`infra/dns`. GitLab consumes the reduced shared catalog before its initial
repository import. Component infrastructure definitions remain with their
owners; this change records the cross-repository migration and acceptance.

The retirement inventory contains twelve rule modules and eleven matching
landing repositories; `rules_openspec` has no landing repository. The exact
cohort, verified deployment outcomes, and remaining delivery checks are
recorded in [design.md](design.md).
