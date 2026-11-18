## Why

Mobile pull request diffs overflow the screen, clipping file controls. The
instance also needs more files per diff page and a managed branch cleanup default.

## What Changes

- Preselect deletion of eligible pull request branches after web merges in
  Terraform-managed Forgejo repositories.
- Render up to 100 diff files per page.
- Keep the pull request diff and its controls within narrow viewports.

## Capabilities

### New Capabilities

- `forgejo-pull-request-review`: Instance and repository settings for pull request review.

### Modified Capabilities

None.

## Impact

Forgejo service Terraform and the Ansible deployment for the instance.
