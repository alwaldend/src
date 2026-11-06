## 1. Rental authentication

- [x] 1.1 Wire the child rental module to an aliased provider using a dedicated alwaldend key; verify packaged configuration without context switching.
- [ ] 1.2 Retrieve the alwaldend key from Vault, replace it through its owning workflow if needed for full team access, and verify the selected team through a read-only request.
- [ ] 1.3 Provision the full-access alwaldend key and its Vault entry through the existing team catalog; verify the actual key context and permissions before renting.
- [x] 1.4 Reuse the shared team module, remove duplicated issuer injection, and verify packaged configuration offline.
- [ ] 1.5 Reconcile alwaldend's existing team address into its team-module address within the same bootstrap state before key apply; verify no team creation or replacement, remove temporary migration blocks, and verify the final plan.

## 2. Rental lifecycle

- [x] 2.1 Add separate ordinary and test rental states, cheapest-offer selection, and interactive apply/release targets; verify formatting and packaged inputs.
- [x] 2.2 Implement the explicit lifecycle test with cleanup on failure and a sanitized receipt; verify its offline failure cases.
- [x] 2.3 Replace affected saved-plan workflows and documentation; verify the generated command arguments.
- [x] 2.4 Correct graceful cancellation through the nested Terraform launcher; verify interrupted creation still reaches cleanup through an end-to-end process fixture.

## 3. Acceptance and delivery

- [ ] 3.1 Execute the authorized cheapest-offer live test and retain a receipt proving the test contract was released.
- [ ] 3.2 Pass repository quality, affected semantic lint, owner OpenSpec checks, local review, and deliver the feature branch and pull request.

## 4. Catalog consolidation

- [x] 4.1 Before implementation, exercise the real catalog root with mocked provider identities and failure cases for duplicate Vault destinations and invalid identity selection.
- [x] 4.2 Remove the separate credentials root and backend, declare the rental key in the existing team catalog, and retain caller-owned Vault writes through their matching AppRoles.
- [x] 4.3 Verify catalog behavior and packaged targets, update maintained documentation, and publish the corrected PR.
