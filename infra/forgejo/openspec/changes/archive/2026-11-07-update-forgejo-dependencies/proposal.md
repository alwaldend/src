## Why

The checked-in Forgejo dependencies lag their stable upstream releases. Update
the existing pins and verify local consumers without deploying infrastructure.

## What Changes

- Align the Forgejo source archive and Linux binary at 17.0.1 (from 15.0.3).
- Update the Actions runner to 13.2.1 (from 13.1.0) and CLI to 0.7.0 (from 0.5.0).
- Update the Forgejo Terraform provider and consuming constraint to 1.7.0
  (from 1.5.0), with verified archive integrity.
- Validate downloads, local version output, packaging, and repository checks.

## Capabilities

No baseline requirements change. This dependency maintenance record uses
`skip_specs: true`; upstream behavior changes remain upstream release contracts.

## Impact

Existing `third_party/org_codeberg_forgejo*` pins, the shared provider archive in
`third_party/terraform`, and `infra/forgejo/tf/provider.tf`. The record is owned
by the Forgejo service and includes its existing runner and CLI dependencies.
No new dependencies, deployment, live plan, authentication, or service operations
are authorized or required.
