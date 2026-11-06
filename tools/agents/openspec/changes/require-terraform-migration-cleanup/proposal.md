## Why

Completed Terraform reconciliations have left temporary migration blocks in
delivered configuration. The infrastructure procedure must include cleanup.

## What Changes

Require removal of temporary import and moved blocks after the authorized
migration, followed by an owning no-change plan. Preserve necessary pending
migrations and report their incomplete status without inventing deployment
authority. This is a procedure correction with no infrastructure deployment.

## Impact

The canonical repo-infra Terraform reference and its offline evaluation cases.
