## Ownership

tools/vastai/cmd/latest_image owns the bounded public Forge tag lookup. tools/vastai/test/latest_image owns its executable-protocol E2E fixture. The personal Terraform root declares the tool as runtime data; its relative runfiles path follows the new owner.

## Verification

Run the relocated protocol fixture, personal Terraform checks, affected semantic lint, and repository quality checks. Inspect the packaged path through the owning Terraform runner without renting capacity.
