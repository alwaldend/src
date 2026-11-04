## 1. Bootstrap root

- [x] 1.1 Define team, generated project key, and Vault secret; verify build, formatting, and declared resource scope.
- [x] 1.2 Verify bootstrap behavior through the provider's Terraform CLI fixture and document the unapplied live deployment prerequisite without exposing tokens.
- [x] 1.3 Select the project's existing Vault identity for secret writes and use standard soft deletion, eliminating cross-role policy grants.
- [x] 1.4 Replace the fixed team/key declarations with a catalog using a reusable Terraform module; verify multiple-team lifecycle through fixtures.
- [x] 1.5 Store catalog data in teams.json and remove moved.tf; verify declared runtime inputs and formatting.
- [x] 1.6 Rebase the integration onto the merged SDK and provider; preserve their implementation and validate the remaining Terraform and credential wiring offline.
