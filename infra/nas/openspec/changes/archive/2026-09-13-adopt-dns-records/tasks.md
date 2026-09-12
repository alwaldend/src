## 1. Prepare adoption

- [x] 1.1 Provision the prepared dedicated `src_infra_nas` AppRole and scoped DNS access through the owning Vault Terraform flow; verify the reviewed plan, applied grants, and credential access.
- [x] 1.2 Review the applicable shared operational evidence and its coverage limits, match provider records, and enable the adoption revision; verify exact matches, offline source checks, and the runtime ownership linter.

## 2. Adopt and verify

- [x] 2.1 Import existing records through `//infra/nas/tf:tf.import` and verify a reviewed plan with zero record additions, changes, replacements, or deletions and unchanged non-DNS resources.
- [x] 2.2 Apply the reviewed adoption plan and verify the declared DNS views and preservation of unrelated records against the pre-adoption inventories.
- [x] 2.3 Record sanitized adoption receipts, verify enabled source defaults and all acceptance scenarios, then synchronize the specification through the archive workflow and validate it.
