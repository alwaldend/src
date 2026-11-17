## Why

The supplied Forgejo run at commit `8b91edb37de21007becef5a37d41ffef2371f102`
failed on shallow Git history, missing Android host prerequisites, stale diagram
projections, synthetic secret findings, and blocked external downloads.

## What Changes

- Fetch complete checkout history for release metadata and secret scanning.
- Repair runner prerequisites and proxy routing through their Ansible owner.
- Refresh diagram projections through their existing generator.
- Review and narrowly suppress demonstrated synthetic secret findings.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

None for the repository build contract. Runner networking belongs to the
owner-local `infra/dev` change. This record repairs existing CI behavior.

## Impact

Forgejo checkout, runner Ansible configuration, generated blog images, and
the repository secret-scan policy. Existing dependency pins remain unchanged.
Acceptance requires representative release builds, diagram comparison tests,
the secret scan, runner configuration validation, and verified proxy bypass.
