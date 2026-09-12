## 1. Prepare adoption

- [x] 1.1 Verify scoped DNS access for the existing `src_infra_dc1_vault` AppRole through authenticated owner commands, retaining the existing identity, grants, and setup backend.
- [x] 1.2 Complete the applicable shared cutover prerequisites, match provider records, and enable the adoption revision; verify exact matches and source hashes, and reuse the baseline source-check and ownership-linter evidence for unchanged declarations.

## 2. Adopt and verify

- [x] 2.1 Prepare exact declarative imports through `//infra/vault/tf_setup:dns.plan`, inspect the saved plan with `dns.show`, and verify every managed action is no-op within `module.dns`.
- [x] 2.2 Apply the reviewed file through `dns.apply`, verify a no-change follow-up plan, and verify the declared DNS views and preservation of unrelated records against complete pre-adoption inventories.
- [x] 2.3 Record sanitized adoption receipts, verify enabled source defaults and all acceptance scenarios, then synchronize the specification through the archive workflow and validate it.
