## 1. Source and acceptance

- [x] 1.1 Add the independent root; verify
      packaged provider, backend, data boundaries and offline checks.
- [x] 1.2 Remove the legacy DNS configuration/commands from PR #233.

## 2. Authorized ownership migration

- [x] 2.1 Inspect private state backups and record an exact migration inventory.
- [x] 2.2 Review and apply only forget/import saved plans; verify unchanged
      record IDs, public DNS and no-change destination/retirement plans.
- [x] 2.3 Remove private artifacts and prepare the validated DNS PR and
      provisioning cleanup. Git and delivery receipts own publication state.

## 3. Remove unnecessary scaffolding

- [x] 3.1 Remove the one-time migration root, legacy backend call, import map
      and instructions; retain only the active independent DNS project.
- [x] 3.2 Verify the simplified project and prepare the corrected PR;
      delivery receipts own publication state.
