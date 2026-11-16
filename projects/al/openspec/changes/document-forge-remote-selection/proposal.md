## Why

Delivery guidance does not explicitly keep PR creation on the selected Git
remote's forge, and its Forgejo reference describes an obsolete fj version.

## What Changes

- Make the selected origin/upstream remote authoritative for PR destination.
- Document the pinned fj command, native login and current capabilities.
- Add offline evaluation cases for GitHub and Forgejo origins, upstream selection
  and mirrors. Preserve existing adapter and publication guards.

## Capabilities

No software requirement changes; this is a documentation and skill-evaluation
update with `skip_specs: true`.

## Impact

The AL delivery skill, its existing eval suite and documentation. The related
repository-policy change is `infra/src/openspec/changes/prefer-configured-pr-forge`;
host installation is `infra/dev/openspec/changes/install-forgejo-cli`.
