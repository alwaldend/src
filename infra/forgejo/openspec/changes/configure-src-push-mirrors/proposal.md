## Why

Forgejo is the primary `alwaldend/src` repository, but its GitHub and GitLab
copies need ongoing push synchronization.

## What Changes

- Bootstrap Forgejo native HTTPS push mirrors from a dedicated Terraform root.
  The mirrors use Vault-backed destination tokens and transfer Git LFS objects.
- Synchronize on commits with an eight-hour periodic retry, preserving existing
  matching mirrors and refusing conflicting or ambiguous configuration.
- Reuse each service's Vault authentication on the controller. Persistent HTTPS
  mirror tokens are separate credentials, not controller API tokens.

## Capabilities

### New Capabilities

- `repository-push-mirrors`: Scoped, repeatable setup of the two src mirrors.

### Modified Capabilities

None.

## Impact

Forgejo mirror Terraform packaging and AL configuration, plus catalog consumer
visibility. Use a packaged API tool with the existing dependencies. No host
deployment, repository recreation, branch-policy changes, or key deletion.
