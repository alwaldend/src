## 1. Prepare scoped access

- [x] 1.1 Verify scoped DNS grants for the existing `src_infra_threexui` AppRole
      through authenticated provider and backend operations with `dns=1`,
      preserving the owning state and ordinary service authentication.
- [x] 1.2 Record the owner's verified prerequisites in the shared
      [cutover procedure](../../../../../dns/openspec/changes/archive/2026-09-13-migrate-project-dns-to-terraform/cutover.md)
      with writer-audit coverage, applicable recovery evidence, and exact
      existing-record matches for every declared view.

## 2. Adopt existing records

- [x] 2.1 Retain the enabled DNS default and document the equivalent IPv6 spelling
      correction, passing declaration test, source hashes, and wrapper builds;
      leave final aggregate source validation to delivery.
- [x] 2.2 Review 12 exact-ID imports through checked-in import blocks and
      `//infra/threexui/tf_setup:dns.plan`, then apply the reviewed saved plan
      through `:dns.apply` with zero record changes and no non-DNS managed actions.
- [x] 2.3 Verify all ten distinct DNS queries and complete provider-inventory
      equality, then verify the scoped follow-up plan has no changes or imports.

## 3. Record adoption

- [x] 3.1 Record sanitized adoption evidence and the DNS-only validation boundary;
      synchronize and strictly validate the specification with every existing
      scenario name preserved, then archive this owner adoption.
