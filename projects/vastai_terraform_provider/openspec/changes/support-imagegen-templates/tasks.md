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
