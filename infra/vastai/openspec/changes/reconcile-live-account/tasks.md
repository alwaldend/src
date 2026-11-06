## 1. Inventory and import preparation

- [x] 1.1 Compare live account inventory with both owning Terraform states and verify both roots produce no-change plans.
- [x] 1.2 Expose the authenticated import command and verify its packaged launcher and help output.
- [x] 1.3 Document observed ownership and reconciliation blockers without publishing secret-bearing inventory.
- [x] 1.4 Validate the changes with Terraform formatting, repository quality, affected semantic lint, and local review.

## 2. Remaining reconciliation

- [x] 2.1 Obtain the intended ENABLE_HTTPS value, declare it as true, and verify the existing variable is imported at its declared resource address.
- [x] 2.2 Apply the reviewed HTTPS value update after explicit live-operation approval and verify a fresh no-change plan.
- [x] 2.3 Import the three existing team identities, remove import blocks, and verify their declared IDs and a no-change plan; record team-scoped inventory limits.
- [x] 2.4 Revoke only unmanaged team API key 36032627 outside IaC as requested and verify that account, owner-context, and Terraform-owned keys are preserved.
- [x] 2.5 Remove hardcoded team IDs and their postcondition, retain the same resource addresses and destruction guards, and verify a no-change plan.
- [x] 2.6 Use a name-keyed map of team configuration objects and verify a no-change plan with the same resource addresses.
- [x] 2.7 Rename the map and resource key to teams, apply the reviewed state migration, and verify a no-change plan without recreating teams.
- [x] 2.8 Consolidate all four teams in one configuration catalog and verify a no-change plan preserving team, module, and Vault secret addresses.
- [x] 2.9 Remove the completed resource-address migration block and verify the owning no-change plan without temporary blocks.
