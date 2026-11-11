## 1. Reusable module

- [x] 1.1 Add template and offer inputs, SSH-only startup, forwarding outputs, packaging, and documentation.
- [x] 1.2 Verify the real module through offline Terraform CLI failure cases and lifecycle evidence.

## 2. Delivery

- [x] 2.1 Run formatting, scoped checks, semantic lint, and repository quality checks.
- [x] 2.2 Review the complete candidate and prepare delivery without renting capacity.

## 3. Shared ComfyUI runtime declaration

- [x] 3.1 Extract fixed runtime defaults into one JSON declaration consumed by the existing module and ephemeral evaluator; preserve personal provisioning.
- [x] 3.2 Export the existing startup template for pure offline compilation and verify unchanged module startup behavior.

## 4. Native model patches

- [x] 4.1 Allow the requested USO projector's model_patch download kind and map it to ComfyUI's models/model_patches directory, retaining URL/filename validation.
- [x] 4.2 Verify destinations and startup limits through the existing configured-profile fixture without creating rentals.
