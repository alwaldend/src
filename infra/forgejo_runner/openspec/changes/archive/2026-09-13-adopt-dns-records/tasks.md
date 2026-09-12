## 1. Prepare adoption

- [x] 1.1 Review the shared writer audit and recovery evidence within their
      documented coverage, and verify scoped authenticated access through the
      existing `src_infra_forgejo_runner` AppRole and setup backend.
- [x] 1.2 Match the canonical declaration to the authorized live inventory in
      every managed view and record an unambiguous existing provider ID for each record.

## 2. Transfer ownership

- [x] 2.1 Verify the enabled `tf_setup/dns.tf` source shape and receipt hash,
      passing scoped wrapper builds, and unchanged canonical declaration covered
      by the baseline ownership linter; update the setup documentation.
- [x] 2.2 Plan the matched imports through
      `//infra/forgejo_runner/tf_setup:dns.plan` and inspect the saved plan through
      `:dns.show`; verify exact provider IDs and only no-op DNS resource actions.

## 3. Verify and record adoption

- [x] 3.1 Apply the reviewed saved plan through
      `//infra/forgejo_runner/tf_setup:dns.apply`; verify the no-change follow-up
      DNS plan, owned record resolution, and preservation of all RouterOS records.
- [x] 3.2 Record the reviewed source revision, owning root, imported IDs,
      no-change plan/apply result, and DNS comparisons in `design.md`; synchronize
      the owner specification only after this evidence supports adoption completion.
