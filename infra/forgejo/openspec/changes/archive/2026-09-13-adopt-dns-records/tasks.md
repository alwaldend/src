## 1. Prepare adoption

- [x] 1.1 Review the shared writer and recovery observations and their coverage
      limits, and verify access through authenticated DNS operations with the
      existing `src_infra_dc1_forgejo1` AppRole; record the observed results.
- [x] 1.2 Match the canonical declaration to the authorized live inventory in
      every managed view and record an unambiguous existing provider ID for each record.

## 2. Transfer ownership

- [x] 2.1 Enable `dns_enabled` by default in `tf_setup/dns.tf` and update the
      setup documentation; verify the receipt's source hash and successful scoped
      wrapper builds, and reuse baseline ownership-linter evidence for the
      unchanged canonical declaration, leaving aggregate validation to delivery.
- [x] 2.2 Prepare the saved import plan through
      `//infra/forgejo/tf_setup:dns.plan` and inspect it through `:dns.show`;
      verify the `dns=1` boundary targets `module.dns`, each planned import ID
      matches its inspected record, and every managed action is `no-op`.

## 3. Verify and record adoption

- [x] 3.1 Apply the reviewed saved plan through
      `//infra/forgejo/tf_setup:dns.apply` and verify a subsequent `:dns.plan`
      has no imports or managed changes; verify owned record
      resolution in each managed view and confirm that owned and unrelated records
      remain unchanged against the before/after inventory.
- [x] 3.2 Record the reviewed source revision, owning root, imported IDs,
      no-change plan/apply result, and DNS comparisons in `design.md`; synchronize
      the owner specification only after this evidence supports adoption completion.
