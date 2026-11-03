## Why

The current repo-infra guidance sends DNS to provisioning roots, coupling
credentials, variables and state. The user requested a separate-root preference.

## What Changes

Prefer component-owned dns/ Terraform roots with independent state. Retain
canonical DNS ownership and explain explicit migration for existing records.
Update conflicting references and add an offline-validated evaluation case.

## Impact

Only canonical repo-infra instructions and their evaluation fixtures change.
This preference grants no live-operation authority and does not require
reorganizing unrelated existing components.
