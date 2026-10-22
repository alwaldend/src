## Why

Deferred; implementation has not started. Review finding 6 identifies a
shared Bazel startup setting that assumes a Fedora-style Java truststore path,
while the repository entry point mainly routes readers to agent policy and
external setup documentation. The first supported contributor build needs a
short, verified path with explicit host and certificate prerequisites.

## What Changes

- Present the supported host baseline and one minimal successful development
  workflow before detailed repository reference material.
- Make custom trust configuration an explicit host setup concern, preserving
  certificate verification and existing private-CA users.
- Link authoritative setup and component documentation rather than maintaining
  duplicate installation procedures.
- Accept with a fresh supported environment completing the documented smoke
  workflow, including a host without the current hardcoded truststore path.

## Capabilities

### New Capabilities

- `contributor-setup`: supported host onboarding and explicit trust configuration.

### Modified Capabilities

None.

## Impact

Planning only is authorized. Future implementation concerns root entry-point
documentation, shared Bazel startup configuration, and any host-specific IaC
needed to preserve supported custom trust. Host changes must be checked in;
there is no live configuration or new dependency authorization. CI is excluded.

Anki-specific documentation belongs to the linked
[product onboarding change](../../../../../projects/anki_as_code/openspec/changes/simplify-user-onboarding/proposal.md).
Baseline: `e0d24cb3a`, inspected 2026-10-22. Evidence:
[root README](../../../../../README.md) and
[shared truststore option](../../../../../tools/bazelrc/project.bazelrc).
Portability is not yet proven or disproven by an executed fresh-host build.
