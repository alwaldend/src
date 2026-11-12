## 1. Implementation

- [x] 1.1 Add template protocol support and preserve existing behavior.
- [x] 1.2 Verify failure cases and lifecycle through offline Terraform acceptance.

## 2. Delivery

- [x] 2.1 Complete scoped checks, local review, and delivery gates.

## 3. Startup transport

- [x] 3.1 Cover compression execution, refresh stability, drift and oversized startup rejection through the real Terraform CLI before implementation.
- [x] 3.2 Add opt-in gzip startup transport and the API byte limit to the instance resource without exposing sensitive contents.
- [x] 3.3 Validate the provider and image-generation consumers and prepare the reviewed fix for delivery; the receipt owns publication state.
- [x] 3.4 Separate synthetic owner-key IDs from generated IDs in the existing team fixture so cleanup does not depend on cross-team execution order.
- [x] 3.5 Consume the SDK-owned startup encoder without changing provider transport behavior; verify existing Terraform CLI fixtures.
- [x] 3.6 Share existing published-port normalization through the SDK while preserving duplicate-binding behavior and invalid/ambiguous diagnostics; verify the controller and provider fixtures.

## 4. Configured model profiles

- [x] 4.1 Extend configured-profile acceptance to every model profile and resolve the shared public asset registry through real Terraform before implementation of the new profile configuration.
- [x] 4.2 Verify model destinations, workflow preservation and startup byte limits using only loopback fixture resources; retain a per-profile artifact.
- [x] 4.3 Verify optional startup workflow selection, preservation of the complete catalog, and the HiDream API-staging profile without live resources.
- [x] 4.4 Extend configured-profile E2E destination checks for the native LoRA and model-patch assets before adding the discussed USO/UMO bundles; retain offline artifacts for every profile.
